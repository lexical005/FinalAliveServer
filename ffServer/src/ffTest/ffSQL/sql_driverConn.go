package main

import (
	"database/sql/driver"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// driverConn wraps a driver.Conn with a mutex, to
// be held during all calls into the Conn. (including any calls onto
// interfaces returned via that Conn, such as calls on Tx, Stmt,
// Result, Rows)
type driverConn struct {
	db        *DB
	createdAt time.Time

	sync.Mutex  // guards following
	ci          driver.Conn
	closed      bool
	finalClosed bool // ci.Close has been called
	openStmt    map[driver.Stmt]bool

	// guarded by db.mu
	inUse      bool
	onPut      []func() // code (with db.mu held) run when conn is next returned
	dbmuClosed bool     // same as closed, but guarded by db.mu, for removeClosedStmtLocked
}

func (dc *driverConn) releaseConn(err error) {
	dc.db.putConn(dc, err)
}

func (dc *driverConn) removeOpenStmt(si driver.Stmt) {
	dc.Lock()
	defer dc.Unlock()
	delete(dc.openStmt, si)
}

func (dc *driverConn) expired(timeout time.Duration) bool {
	if timeout <= 0 {
		return false
	}
	return dc.createdAt.Add(timeout).Before(nowFunc())
}

func (dc *driverConn) prepareLocked(query string) (driver.Stmt, error) {
	si, err := dc.ci.Prepare(query)
	if err == nil {
		// Track each driverConn's open statements, so we can close them
		// before closing the conn.
		//
		// TODO(bradfitz): let drivers opt out of caring about
		// stmt closes if the conn is about to close anyway? For now
		// do the safe thing, in case stmts need to be closed.
		//
		// TODO(bradfitz): after Go 1.2, closing driver.Stmts
		// should be moved to driverStmt, using unique
		// *driverStmts everywhere (including from
		// *Stmt.connStmt, instead of returning a
		// driver.Stmt), using driverStmt as a pointer
		// everywhere, and making it a finalCloser.
		if dc.openStmt == nil {
			dc.openStmt = make(map[driver.Stmt]bool)
		}
		dc.openStmt[si] = true
	}
	return si, err
}

// the dc.db's Mutex is held.
func (dc *driverConn) closeDBLocked() func() error {
	dc.Lock()
	defer dc.Unlock()
	if dc.closed {
		return func() error { return errors.New("sql: duplicate driverConn close") }
	}
	dc.closed = true
	return dc.db.removeDepLocked(dc, dc)
}

func (dc *driverConn) Close() error {
	dc.Lock()
	if dc.closed {
		dc.Unlock()
		return errors.New("sql: duplicate driverConn close")
	}
	dc.closed = true
	dc.Unlock() // not defer; removeDep finalClose calls may need to lock

	// And now updates that require holding dc.mu.Lock.
	dc.db.mu.Lock()
	dc.dbmuClosed = true
	fn := dc.db.removeDepLocked(dc, dc)
	dc.db.mu.Unlock()
	return fn()
}

func (dc *driverConn) finalClose() error {
	dc.Lock()

	for si := range dc.openStmt {
		si.Close()
	}
	dc.openStmt = nil

	err := dc.ci.Close()
	dc.ci = nil
	dc.finalClosed = true
	dc.Unlock()

	dc.db.mu.Lock()
	dc.db.numOpen--
	dc.db.maybeOpenNewConnections()
	dc.db.mu.Unlock()

	atomic.AddUint64(&dc.db.numClosed, 1)
	return err
}
