package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

// connStmt is a prepared statement on a particular connection.
type connStmt struct {
	dc *driverConn
	si driver.Stmt
}

// Stmt is a prepared statement.
// A Stmt is safe for concurrent use by multiple goroutines.
type Stmt struct {
	// Immutable:
	db        *DB    // where we came from
	query     string // that created the Stmt
	stickyErr error  // if non-nil, this error is returned for all operations

	closemu sync.RWMutex // held exclusively during close, for read otherwise.

	// If in a transaction, else both nil:
	tx   *Tx
	txsi *driverStmt

	mu     sync.Mutex // protects the rest of the fields
	closed bool

	// css is a list of underlying driver statement interfaces
	// that are valid on particular connections. This is only
	// used if tx == nil and one is found that has idle
	// connections. If tx != nil, txsi is always used.
	css []connStmt

	// lastNumClosed is copied from db.numClosed when Stmt is created
	// without tx and closed connections in css are removed.
	lastNumClosed uint64
}

// Exec executes a prepared statement with the given arguments and
// returns a Result summarizing the effect of the statement.
func (s *Stmt) Exec(args ...interface{}) (Result, error) {
	s.closemu.RLock()
	defer s.closemu.RUnlock()

	var res Result
	for i := 0; i < maxBadConnRetries; i++ {
		dc, releaseConn, si, err := s.connStmt()
		if err != nil {
			if err == driver.ErrBadConn {
				continue
			}
			return nil, err
		}

		res, err = resultFromStatement(driverStmt{dc, si}, args...)
		releaseConn(err)
		if err != driver.ErrBadConn {
			return res, err
		}
	}
	return nil, driver.ErrBadConn
}

func driverNumInput(ds driverStmt) int {
	ds.Lock()
	defer ds.Unlock() // in case NumInput panics
	return ds.si.NumInput()
}

func resultFromStatement(ds driverStmt, args ...interface{}) (Result, error) {
	want := driverNumInput(ds)

	// -1 means the driver doesn't know how to count the number of
	// placeholders, so we won't sanity check input here and instead let the
	// driver deal with errors.
	if want != -1 && len(args) != want {
		return nil, fmt.Errorf("sql: expected %d arguments, got %d", want, len(args))
	}

	dargs, err := driverArgs(&ds, args)
	if err != nil {
		return nil, err
	}

	ds.Lock()
	defer ds.Unlock()
	resi, err := ds.si.Exec(dargs)
	if err != nil {
		return nil, err
	}
	return driverResult{ds.Locker, resi}, nil
}

// removeClosedStmtLocked removes closed conns in s.css.
//
// To avoid lock contention on DB.mu, we do it only when
// s.db.numClosed - s.lastNum is large enough.
func (s *Stmt) removeClosedStmtLocked() {
	t := len(s.css)/2 + 1
	if t > 10 {
		t = 10
	}
	dbClosed := atomic.LoadUint64(&s.db.numClosed)
	if dbClosed-s.lastNumClosed < uint64(t) {
		return
	}

	s.db.mu.Lock()
	for i := 0; i < len(s.css); i++ {
		if s.css[i].dc.dbmuClosed {
			s.css[i] = s.css[len(s.css)-1]
			s.css = s.css[:len(s.css)-1]
			i--
		}
	}
	s.db.mu.Unlock()
	s.lastNumClosed = dbClosed
}

// connStmt returns a free driver connection on which to execute the
// statement, a function to call to release the connection, and a
// statement bound to that connection.
func (s *Stmt) connStmt() (ci *driverConn, releaseConn func(error), si driver.Stmt, err error) {
	if err = s.stickyErr; err != nil {
		return
	}
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		err = errors.New("sql: statement is closed")
		return
	}

	// In a transaction, we always use the connection that the
	// transaction was created on.
	if s.tx != nil {
		s.mu.Unlock()
		ci, err = s.tx.grabConn() // blocks, waiting for the connection.
		if err != nil {
			return
		}
		releaseConn = func(error) {}
		return ci, releaseConn, s.txsi.si, nil
	}

	s.removeClosedStmtLocked()
	s.mu.Unlock()

	// TODO(bradfitz): or always wait for one? make configurable later?
	dc, err := s.db.conn(cachedOrNewConn)
	if err != nil {
		return nil, nil, nil, err
	}

	s.mu.Lock()
	for _, v := range s.css {
		if v.dc == dc {
			s.mu.Unlock()
			return dc, dc.releaseConn, v.si, nil
		}
	}
	s.mu.Unlock()

	// No luck; we need to prepare the statement on this connection
	dc.Lock()
	si, err = dc.prepareLocked(s.query)
	dc.Unlock()
	if err != nil {
		s.db.putConn(dc, err)
		return nil, nil, nil, err
	}
	s.mu.Lock()
	cs := connStmt{dc, si}
	s.css = append(s.css, cs)
	s.mu.Unlock()

	return dc, dc.releaseConn, si, nil
}

// Query executes a prepared query statement with the given arguments
// and returns the query results as a *Rows.
func (s *Stmt) Query(args ...interface{}) (*Rows, error) {
	s.closemu.RLock()
	defer s.closemu.RUnlock()

	var rowsi driver.Rows
	for i := 0; i < maxBadConnRetries; i++ {
		dc, releaseConn, si, err := s.connStmt()
		if err != nil {
			if err == driver.ErrBadConn {
				continue
			}
			return nil, err
		}

		rowsi, err = rowsiFromStatement(driverStmt{dc, si}, args...)
		if err == nil {
			// Note: ownership of ci passes to the *Rows, to be freed
			// with releaseConn.
			rows := &Rows{
				dc:    dc,
				rowsi: rowsi,
				// releaseConn set below
			}
			s.db.addDep(s, rows)
			rows.releaseConn = func(err error) {
				releaseConn(err)
				s.db.removeDep(s, rows)
			}
			return rows, nil
		}

		releaseConn(err)
		if err != driver.ErrBadConn {
			return nil, err
		}
	}
	return nil, driver.ErrBadConn
}

func rowsiFromStatement(ds driverStmt, args ...interface{}) (driver.Rows, error) {
	ds.Lock()
	want := ds.si.NumInput()
	ds.Unlock()

	// -1 means the driver doesn't know how to count the number of
	// placeholders, so we won't sanity check input here and instead let the
	// driver deal with errors.
	if want != -1 && len(args) != want {
		return nil, fmt.Errorf("sql: statement expects %d inputs; got %d", want, len(args))
	}

	dargs, err := driverArgs(&ds, args)
	if err != nil {
		return nil, err
	}

	ds.Lock()
	rowsi, err := ds.si.Query(dargs)
	ds.Unlock()
	if err != nil {
		return nil, err
	}
	return rowsi, nil
}

// QueryRow executes a prepared query statement with the given arguments.
// If an error occurs during the execution of the statement, that error will
// be returned by a call to Scan on the returned *Row, which is always non-nil.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
//
// Example usage:
//
//  var name string
//  err := nameByUseridStmt.QueryRow(id).Scan(&name)
func (s *Stmt) QueryRow(args ...interface{}) *Row {
	rows, err := s.Query(args...)
	if err != nil {
		return &Row{err: err}
	}
	return &Row{rows: rows}
}

// Close closes the statement.
func (s *Stmt) Close() error {
	s.closemu.Lock()
	defer s.closemu.Unlock()

	if s.stickyErr != nil {
		return s.stickyErr
	}
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true

	if s.tx != nil {
		err := s.txsi.Close()
		s.mu.Unlock()
		return err
	}
	s.mu.Unlock()

	return s.db.removeDep(s, s)
}

func (s *Stmt) finalClose() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.css != nil {
		for _, v := range s.css {
			s.db.noteUnusedDriverStatement(v.dc, v.si)
			v.dc.removeOpenStmt(v.si)
		}
		s.css = nil
	}
	return nil
}
