// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sql provides a generic interface around SQL (or SQL-like)
// databases.
//
// The sql package must be used in conjunction with a database driver.
// See https://golang.org/s/sqldrivers for a list of drivers.
//
// For more usage examples, see the wiki page at
// https://golang.org/s/sqlwiki.
package main

import (
	"database/sql/driver"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// DB is a database handle representing a pool of zero or more
// underlying connections. It's safe for concurrent use by multiple
// goroutines.
//
// The sql package creates and frees connections automatically; it
// also maintains a free pool of idle connections. If the database has
// a concept of per-connection state, such state can only be reliably
// observed within a transaction. Once DB.Begin is called, the
// returned Tx is bound to a single connection. Once Commit or
// Rollback is called on the transaction, that transaction's
// connection is returned to DB's idle connection pool. The pool size
// can be controlled with SetMaxIdleConns.
type DB struct {
	driver driver.Driver
	dsn    string
	// numClosed is an atomic counter which represents a total number of
	// closed connections. Stmt.openStmt checks it before cleaning closed
	// connections in Stmt.css.
	numClosed uint64

	mu           sync.Mutex // protects following fields
	freeConn     []*driverConn
	connRequests []chan connRequest
	numOpen      int // number of opened and pending open connections
	// Used to signal the need for new connections
	// a goroutine running connectionOpener() reads on this chan and
	// maybeOpenNewConnections sends on the chan (one send per needed connection)
	// It is closed during db.Close(). The close tells the connectionOpener
	// goroutine to exit.
	openerCh    chan struct{}
	closed      bool
	dep         map[finalCloser]depSet
	lastPut     map[*driverConn]string // stacktrace of last conn's put; debug only
	maxIdle     int                    // zero means defaultMaxIdleConns; negative means 0
	maxOpen     int                    // <= 0 means unlimited
	maxLifetime time.Duration          // maximum amount of time a connection may be reused
	cleanerCh   chan struct{}
}

// addDep notes that x now depends on dep, and x's finalClose won't be
// called until all of x's dependencies are removed with removeDep.
func (db *DB) addDep(x finalCloser, dep interface{}) {
	//println(fmt.Sprintf("addDep(%T %p, %T %p)", x, x, dep, dep))
	db.mu.Lock()
	defer db.mu.Unlock()
	db.addDepLocked(x, dep)
}

func (db *DB) addDepLocked(x finalCloser, dep interface{}) {
	if db.dep == nil {
		db.dep = make(map[finalCloser]depSet)
	}
	xdep := db.dep[x]
	if xdep == nil {
		xdep = make(depSet)
		db.dep[x] = xdep
	}
	xdep[dep] = true
}

// removeDep notes that x no longer depends on dep.
// If x still has dependencies, nil is returned.
// If x no longer has any dependencies, its finalClose method will be
// called and its error value will be returned.
func (db *DB) removeDep(x finalCloser, dep interface{}) error {
	db.mu.Lock()
	fn := db.removeDepLocked(x, dep)
	db.mu.Unlock()
	return fn()
}

func (db *DB) removeDepLocked(x finalCloser, dep interface{}) func() error {
	//println(fmt.Sprintf("removeDep(%T %p, %T %p)", x, x, dep, dep))

	xdep, ok := db.dep[x]
	if !ok {
		panic(fmt.Sprintf("unpaired removeDep: no deps for %T", x))
	}

	l0 := len(xdep)
	delete(xdep, dep)

	switch len(xdep) {
	case l0:
		// Nothing removed. Shouldn't happen.
		panic(fmt.Sprintf("unpaired removeDep: no %T dep on %T", dep, x))
	case 0:
		// No more dependencies.
		delete(db.dep, x)
		return x.finalClose
	default:
		// Dependencies remain.
		return func() error { return nil }
	}
}

// Ping verifies a connection to the database is still alive,
// establishing a connection if necessary.
func (db *DB) Ping() error {
	// TODO(bradfitz): give drivers an optional hook to implement
	// this in a more efficient or more reliable way, if they
	// have one.
	dc, err := db.conn(cachedOrNewConn)
	if err != nil {
		return err
	}
	db.putConn(dc, nil)
	return nil
}

// Close closes the database, releasing any open resources.
//
// It is rare to Close a DB, as the DB handle is meant to be
// long-lived and shared between many goroutines.
func (db *DB) Close() error {
	db.mu.Lock()
	if db.closed { // Make DB.Close idempotent
		db.mu.Unlock()
		return nil
	}
	close(db.openerCh)
	if db.cleanerCh != nil {
		close(db.cleanerCh)
	}
	var err error
	fns := make([]func() error, 0, len(db.freeConn))
	for _, dc := range db.freeConn {
		fns = append(fns, dc.closeDBLocked())
	}
	db.freeConn = nil
	db.closed = true
	for _, req := range db.connRequests {
		close(req)
	}
	db.mu.Unlock()
	for _, fn := range fns {
		err1 := fn()
		if err1 != nil {
			err = err1
		}
	}
	return err
}

func (db *DB) maxIdleConnsLocked() int {
	n := db.maxIdle
	switch {
	case n == 0:
		// TODO(bradfitz): ask driver, if supported, for its default preference
		return defaultMaxIdleConns
	case n < 0:
		return 0
	default:
		return n
	}
}

// SetMaxIdleConns sets the maximum number of connections in the idle
// connection pool.
//
// If MaxOpenConns is greater than 0 but less than the new MaxIdleConns
// then the new MaxIdleConns will be reduced to match the MaxOpenConns limit
//
// If n <= 0, no idle connections are retained.
func (db *DB) SetMaxIdleConns(n int) {
	db.mu.Lock()
	if n > 0 {
		db.maxIdle = n
	} else {
		// No idle connections.
		db.maxIdle = -1
	}
	// Make sure maxIdle doesn't exceed maxOpen
	if db.maxOpen > 0 && db.maxIdleConnsLocked() > db.maxOpen {
		db.maxIdle = db.maxOpen
	}
	var closing []*driverConn
	idleCount := len(db.freeConn)
	maxIdle := db.maxIdleConnsLocked()
	if idleCount > maxIdle {
		closing = db.freeConn[maxIdle:]
		db.freeConn = db.freeConn[:maxIdle]
	}
	db.mu.Unlock()
	for _, c := range closing {
		c.Close()
	}
}

// SetMaxOpenConns sets the maximum number of open connections to the database.
//
// If MaxIdleConns is greater than 0 and the new MaxOpenConns is less than
// MaxIdleConns, then MaxIdleConns will be reduced to match the new
// MaxOpenConns limit
//
// If n <= 0, then there is no limit on the number of open connections.
// The default is 0 (unlimited).
func (db *DB) SetMaxOpenConns(n int) {
	db.mu.Lock()
	db.maxOpen = n
	if n < 0 {
		db.maxOpen = 0
	}
	syncMaxIdle := db.maxOpen > 0 && db.maxIdleConnsLocked() > db.maxOpen
	db.mu.Unlock()
	if syncMaxIdle {
		db.SetMaxIdleConns(n)
	}
}

// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
//
// Expired connections may be closed lazily before reuse.
//
// If d <= 0, connections are reused forever.
func (db *DB) SetConnMaxLifetime(d time.Duration) {
	if d < 0 {
		d = 0
	}
	db.mu.Lock()
	// wake cleaner up when lifetime is shortened.
	if d > 0 && d < db.maxLifetime && db.cleanerCh != nil {
		select {
		case db.cleanerCh <- struct{}{}:
		default:
		}
	}
	db.maxLifetime = d
	db.startCleanerLocked()
	db.mu.Unlock()
}

// startCleanerLocked starts connectionCleaner if needed.
func (db *DB) startCleanerLocked() {
	if db.maxLifetime > 0 && db.numOpen > 0 && db.cleanerCh == nil {
		db.cleanerCh = make(chan struct{}, 1)
		go db.connectionCleaner(db.maxLifetime)
	}
}

func (db *DB) connectionCleaner(d time.Duration) {
	const minInterval = time.Second

	if d < minInterval {
		d = minInterval
	}
	t := time.NewTimer(d)

	for {
		select {
		case <-t.C:
		case <-db.cleanerCh: // maxLifetime was changed or db was closed.
		}

		db.mu.Lock()
		d = db.maxLifetime
		if db.closed || db.numOpen == 0 || d <= 0 {
			db.cleanerCh = nil
			db.mu.Unlock()
			return
		}

		expiredSince := nowFunc().Add(-d)
		var closing []*driverConn
		for i := 0; i < len(db.freeConn); i++ {
			c := db.freeConn[i]
			if c.createdAt.Before(expiredSince) {
				closing = append(closing, c)
				last := len(db.freeConn) - 1
				db.freeConn[i] = db.freeConn[last]
				db.freeConn[last] = nil
				db.freeConn = db.freeConn[:last]
				i--
			}
		}
		db.mu.Unlock()

		for _, c := range closing {
			c.Close()
		}

		if d < minInterval {
			d = minInterval
		}
		t.Reset(d)
	}
}

// Stats returns database statistics.
func (db *DB) Stats() DBStats {
	db.mu.Lock()
	stats := DBStats{
		OpenConnections: db.numOpen,
	}
	db.mu.Unlock()
	return stats
}

// Assumes db.mu is locked.
// If there are connRequests and the connection limit hasn't been reached,
// then tell the connectionOpener to open new connections.
func (db *DB) maybeOpenNewConnections() {
	numRequests := len(db.connRequests)
	if db.maxOpen > 0 {
		numCanOpen := db.maxOpen - db.numOpen
		if numRequests > numCanOpen {
			numRequests = numCanOpen
		}
	}
	for numRequests > 0 {
		db.numOpen++ // optimistically
		numRequests--
		if db.closed {
			return
		}
		db.openerCh <- struct{}{}
	}
}

// Runs in a separate goroutine, opens new connections when requested.
func (db *DB) connectionOpener() {
	for range db.openerCh {
		db.openNewConnection()
	}
}

// Open one new connection
func (db *DB) openNewConnection() {
	// maybeOpenNewConnctions has already executed db.numOpen++ before it sent
	// on db.openerCh. This function must execute db.numOpen-- if the
	// connection fails or is closed before returning.
	ci, err := db.driver.Open(db.dsn)
	db.mu.Lock()
	defer db.mu.Unlock()
	if db.closed {
		if err == nil {
			ci.Close()
		}
		db.numOpen--
		return
	}
	if err != nil {
		db.numOpen--
		db.putConnDBLocked(nil, err)
		db.maybeOpenNewConnections()
		return
	}
	dc := &driverConn{
		db:        db,
		createdAt: nowFunc(),
		ci:        ci,
	}
	if db.putConnDBLocked(dc, err) {
		db.addDepLocked(dc, dc)
	} else {
		db.numOpen--
		ci.Close()
	}
}

// conn returns a newly-opened or cached *driverConn.
func (db *DB) conn(strategy connReuseStrategy) (*driverConn, error) {
	db.mu.Lock()
	if db.closed {
		db.mu.Unlock()
		return nil, errDBClosed
	}
	lifetime := db.maxLifetime

	// Prefer a free connection, if possible.
	numFree := len(db.freeConn)
	if strategy == cachedOrNewConn && numFree > 0 {
		conn := db.freeConn[0]
		copy(db.freeConn, db.freeConn[1:])
		db.freeConn = db.freeConn[:numFree-1]
		conn.inUse = true
		db.mu.Unlock()
		if conn.expired(lifetime) {
			conn.Close()
			return nil, driver.ErrBadConn
		}
		return conn, nil
	}

	// Out of free connections or we were asked not to use one. If we're not
	// allowed to open any more connections, make a request and wait.
	if db.maxOpen > 0 && db.numOpen >= db.maxOpen {
		// Make the connRequest channel. It's buffered so that the
		// connectionOpener doesn't block while waiting for the req to be read.
		req := make(chan connRequest, 1)
		db.connRequests = append(db.connRequests, req)
		db.mu.Unlock()
		ret, ok := <-req
		if !ok {
			return nil, errDBClosed
		}
		if ret.err == nil && ret.conn.expired(lifetime) {
			ret.conn.Close()
			return nil, driver.ErrBadConn
		}
		return ret.conn, ret.err
	}

	db.numOpen++ // optimistically
	db.mu.Unlock()
	ci, err := db.driver.Open(db.dsn)
	if err != nil {
		db.mu.Lock()
		db.numOpen-- // correct for earlier optimism
		db.maybeOpenNewConnections()
		db.mu.Unlock()
		return nil, err
	}
	db.mu.Lock()
	dc := &driverConn{
		db:        db,
		createdAt: nowFunc(),
		ci:        ci,
	}
	db.addDepLocked(dc, dc)
	dc.inUse = true
	db.mu.Unlock()
	return dc, nil
}

// noteUnusedDriverStatement notes that si is no longer used and should
// be closed whenever possible (when c is next not in use), unless c is
// already closed.
func (db *DB) noteUnusedDriverStatement(c *driverConn, si driver.Stmt) {
	db.mu.Lock()
	defer db.mu.Unlock()
	if c.inUse {
		c.onPut = append(c.onPut, func() {
			si.Close()
		})
	} else {
		c.Lock()
		defer c.Unlock()
		if !c.finalClosed {
			si.Close()
		}
	}
}

// putConn adds a connection to the db's free pool.
// err is optionally the last error that occurred on this connection.
func (db *DB) putConn(dc *driverConn, err error) {
	db.mu.Lock()
	if !dc.inUse {
		if debugGetPut {
			fmt.Printf("putConn(%v) DUPLICATE was: %s\n\nPREVIOUS was: %s", dc, stack(), db.lastPut[dc])
		}
		panic("sql: connection returned that was never out")
	}
	if debugGetPut {
		db.lastPut[dc] = stack()
	}
	dc.inUse = false

	for _, fn := range dc.onPut {
		fn()
	}
	dc.onPut = nil

	if err == driver.ErrBadConn {
		// Don't reuse bad connections.
		// Since the conn is considered bad and is being discarded, treat it
		// as closed. Don't decrement the open count here, finalClose will
		// take care of that.
		db.maybeOpenNewConnections()
		db.mu.Unlock()
		dc.Close()
		return
	}
	if putConnHook != nil {
		putConnHook(db, dc)
	}
	added := db.putConnDBLocked(dc, nil)
	db.mu.Unlock()

	if !added {
		dc.Close()
	}
}

// Satisfy a connRequest or put the driverConn in the idle pool and return true
// or return false.
// putConnDBLocked will satisfy a connRequest if there is one, or it will
// return the *driverConn to the freeConn list if err == nil and the idle
// connection limit will not be exceeded.
// If err != nil, the value of dc is ignored.
// If err == nil, then dc must not equal nil.
// If a connRequest was fulfilled or the *driverConn was placed in the
// freeConn list, then true is returned, otherwise false is returned.
func (db *DB) putConnDBLocked(dc *driverConn, err error) bool {
	if db.closed {
		return false
	}
	if db.maxOpen > 0 && db.numOpen > db.maxOpen {
		return false
	}
	if c := len(db.connRequests); c > 0 {
		req := db.connRequests[0]
		// This copy is O(n) but in practice faster than a linked list.
		// TODO: consider compacting it down less often and
		// moving the base instead?
		copy(db.connRequests, db.connRequests[1:])
		db.connRequests = db.connRequests[:c-1]
		if err == nil {
			dc.inUse = true
		}
		req <- connRequest{
			conn: dc,
			err:  err,
		}
		return true
	} else if err == nil && !db.closed && db.maxIdleConnsLocked() > len(db.freeConn) {
		db.freeConn = append(db.freeConn, dc)
		db.startCleanerLocked()
		return true
	}
	return false
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's Close method
// when the statement is no longer needed.
func (db *DB) Prepare(query string) (*Stmt, error) {
	var stmt *Stmt
	var err error
	for i := 0; i < maxBadConnRetries; i++ {
		stmt, err = db.prepare(query, cachedOrNewConn)
		if err != driver.ErrBadConn {
			break
		}
	}
	if err == driver.ErrBadConn {
		return db.prepare(query, alwaysNewConn)
	}
	return stmt, err
}

func (db *DB) prepare(query string, strategy connReuseStrategy) (*Stmt, error) {
	// TODO: check if db.driver supports an optional
	// driver.Preparer interface and call that instead, if so,
	// otherwise we make a prepared statement that's bound
	// to a connection, and to execute this prepared statement
	// we either need to use this connection (if it's free), else
	// get a new connection + re-prepare + execute on that one.
	dc, err := db.conn(strategy)
	if err != nil {
		return nil, err
	}
	dc.Lock()
	si, err := dc.prepareLocked(query)
	dc.Unlock()
	if err != nil {
		db.putConn(dc, err)
		return nil, err
	}
	stmt := &Stmt{
		db:            db,
		query:         query,
		css:           []connStmt{{dc, si}},
		lastNumClosed: atomic.LoadUint64(&db.numClosed),
	}
	db.addDep(stmt, stmt)
	db.putConn(dc, nil)
	return stmt, nil
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
func (db *DB) Exec(query string, args ...interface{}) (Result, error) {
	var res Result
	var err error
	for i := 0; i < maxBadConnRetries; i++ {
		res, err = db.exec(query, args, cachedOrNewConn)
		if err != driver.ErrBadConn {
			break
		}
	}
	if err == driver.ErrBadConn {
		return db.exec(query, args, alwaysNewConn)
	}
	return res, err
}

func (db *DB) exec(query string, args []interface{}, strategy connReuseStrategy) (res Result, err error) {
	dc, err := db.conn(strategy)
	if err != nil {
		return nil, err
	}
	defer func() {
		db.putConn(dc, err)
	}()

	if execer, ok := dc.ci.(driver.Execer); ok {
		dargs, err := driverArgs(nil, args)
		if err != nil {
			return nil, err
		}
		dc.Lock()
		resi, err := execer.Exec(query, dargs)
		dc.Unlock()
		if err != driver.ErrSkip {
			if err != nil {
				return nil, err
			}
			return driverResult{dc, resi}, nil
		}
	}

	dc.Lock()
	si, err := dc.ci.Prepare(query)
	dc.Unlock()
	if err != nil {
		return nil, err
	}
	defer withLock(dc, func() { si.Close() })
	return resultFromStatement(driverStmt{dc, si}, args...)
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *DB) Query(query string, args ...interface{}) (*Rows, error) {
	var rows *Rows
	var err error
	for i := 0; i < maxBadConnRetries; i++ {
		rows, err = db.query(query, args, cachedOrNewConn)
		if err != driver.ErrBadConn {
			break
		}
	}
	if err == driver.ErrBadConn {
		return db.query(query, args, alwaysNewConn)
	}
	return rows, err
}

func (db *DB) query(query string, args []interface{}, strategy connReuseStrategy) (*Rows, error) {
	ci, err := db.conn(strategy)
	if err != nil {
		return nil, err
	}

	return db.queryConn(ci, ci.releaseConn, query, args)
}

// queryConn executes a query on the given connection.
// The connection gets released by the releaseConn function.
func (db *DB) queryConn(dc *driverConn, releaseConn func(error), query string, args []interface{}) (*Rows, error) {
	if queryer, ok := dc.ci.(driver.Queryer); ok {
		dargs, err := driverArgs(nil, args)
		if err != nil {
			releaseConn(err)
			return nil, err
		}
		dc.Lock()
		rowsi, err := queryer.Query(query, dargs)
		dc.Unlock()
		if err != driver.ErrSkip {
			if err != nil {
				releaseConn(err)
				return nil, err
			}
			// Note: ownership of dc passes to the *Rows, to be freed
			// with releaseConn.
			rows := &Rows{
				dc:          dc,
				releaseConn: releaseConn,
				rowsi:       rowsi,
			}
			return rows, nil
		}
	}

	dc.Lock()
	si, err := dc.ci.Prepare(query)
	dc.Unlock()
	if err != nil {
		releaseConn(err)
		return nil, err
	}

	ds := driverStmt{dc, si}
	rowsi, err := rowsiFromStatement(ds, args...)
	if err != nil {
		dc.Lock()
		si.Close()
		dc.Unlock()
		releaseConn(err)
		return nil, err
	}

	// Note: ownership of ci passes to the *Rows, to be freed
	// with releaseConn.
	rows := &Rows{
		dc:          dc,
		releaseConn: releaseConn,
		rowsi:       rowsi,
		closeStmt:   si,
	}
	return rows, nil
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (db *DB) QueryRow(query string, args ...interface{}) *Row {
	rows, err := db.Query(query, args...)
	return &Row{rows: rows, err: err}
}

// Begin starts a transaction. The isolation level is dependent on
// the driver.
func (db *DB) Begin() (*Tx, error) {
	var tx *Tx
	var err error
	for i := 0; i < maxBadConnRetries; i++ {
		tx, err = db.begin(cachedOrNewConn)
		if err != driver.ErrBadConn {
			break
		}
	}
	if err == driver.ErrBadConn {
		return db.begin(alwaysNewConn)
	}
	return tx, err
}

func (db *DB) begin(strategy connReuseStrategy) (tx *Tx, err error) {
	dc, err := db.conn(strategy)
	if err != nil {
		return nil, err
	}
	dc.Lock()
	txi, err := dc.ci.Begin()
	dc.Unlock()
	if err != nil {
		db.putConn(dc, err)
		return nil, err
	}
	return &Tx{
		db:  db,
		dc:  dc,
		txi: txi,
	}, nil
}

// Driver returns the database's underlying driver.
func (db *DB) Driver() driver.Driver {
	return db.driver
}
