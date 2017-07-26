package main

import (
	"database/sql/driver"
	"errors"
	"sync"
)

// Tx is an in-progress database transaction.
//
// A transaction must end with a call to Commit or Rollback.
//
// After a call to Commit or Rollback, all operations on the
// transaction fail with ErrTxDone.
//
// The statements prepared for a transaction by calling
// the transaction's Prepare or Stmt methods are closed
// by the call to Commit or Rollback.
type Tx struct {
	db *DB

	// dc is owned exclusively until Commit or Rollback, at which point
	// it's returned with putConn.
	dc  *driverConn
	txi driver.Tx

	// done transitions from false to true exactly once, on Commit
	// or Rollback. once done, all operations fail with
	// ErrTxDone.
	done bool

	// All Stmts prepared for this transaction. These will be closed after the
	// transaction has been committed or rolled back.
	stmts struct {
		sync.Mutex
		v []*Stmt
	}
}

var ErrTxDone = errors.New("sql: Transaction has already been committed or rolled back")

func (tx *Tx) close(err error) {
	if tx.done {
		panic("double close") // internal error
	}
	tx.done = true
	tx.db.putConn(tx.dc, err)
	tx.dc = nil
	tx.txi = nil
}

func (tx *Tx) grabConn() (*driverConn, error) {
	if tx.done {
		return nil, ErrTxDone
	}
	return tx.dc, nil
}

// Closes all Stmts prepared for this transaction.
func (tx *Tx) closePrepared() {
	tx.stmts.Lock()
	for _, stmt := range tx.stmts.v {
		stmt.Close()
	}
	tx.stmts.Unlock()
}

// Commit commits the transaction.
func (tx *Tx) Commit() error {
	if tx.done {
		return ErrTxDone
	}
	tx.dc.Lock()
	err := tx.txi.Commit()
	tx.dc.Unlock()
	if err != driver.ErrBadConn {
		tx.closePrepared()
	}
	tx.close(err)
	return err
}

// Rollback aborts the transaction.
func (tx *Tx) Rollback() error {
	if tx.done {
		return ErrTxDone
	}
	tx.dc.Lock()
	err := tx.txi.Rollback()
	tx.dc.Unlock()
	if err != driver.ErrBadConn {
		tx.closePrepared()
	}
	tx.close(err)
	return err
}

// Prepare creates a prepared statement for use within a transaction.
//
// The returned statement operates within the transaction and can no longer
// be used once the transaction has been committed or rolled back.
//
// To use an existing prepared statement on this transaction, see Tx.Stmt.
func (tx *Tx) Prepare(query string) (*Stmt, error) {
	// TODO(bradfitz): We could be more efficient here and either
	// provide a method to take an existing Stmt (created on
	// perhaps a different Conn), and re-create it on this Conn if
	// necessary. Or, better: keep a map in DB of query string to
	// Stmts, and have Stmt.Execute do the right thing and
	// re-prepare if the Conn in use doesn't have that prepared
	// statement. But we'll want to avoid caching the statement
	// in the case where we only call conn.Prepare implicitly
	// (such as in db.Exec or tx.Exec), but the caller package
	// can't be holding a reference to the returned statement.
	// Perhaps just looking at the reference count (by noting
	// Stmt.Close) would be enough. We might also want a finalizer
	// on Stmt to drop the reference count.
	dc, err := tx.grabConn()
	if err != nil {
		return nil, err
	}

	dc.Lock()
	si, err := dc.ci.Prepare(query)
	dc.Unlock()
	if err != nil {
		return nil, err
	}

	stmt := &Stmt{
		db: tx.db,
		tx: tx,
		txsi: &driverStmt{
			Locker: dc,
			si:     si,
		},
		query: query,
	}
	tx.stmts.Lock()
	tx.stmts.v = append(tx.stmts.v, stmt)
	tx.stmts.Unlock()
	return stmt, nil
}

// Stmt returns a transaction-specific prepared statement from
// an existing statement.
//
// Example:
//  updateMoney, err := db.Prepare("UPDATE balance SET money=money+? WHERE id=?")
//  ...
//  tx, err := db.Begin()
//  ...
//  res, err := tx.Stmt(updateMoney).Exec(123.45, 98293203)
//
// The returned statement operates within the transaction and can no longer
// be used once the transaction has been committed or rolled back.
func (tx *Tx) Stmt(stmt *Stmt) *Stmt {
	// TODO(bradfitz): optimize this. Currently this re-prepares
	// each time. This is fine for now to illustrate the API but
	// we should really cache already-prepared statements
	// per-Conn. See also the big comment in Tx.Prepare.

	if tx.db != stmt.db {
		return &Stmt{stickyErr: errors.New("sql: Tx.Stmt: statement from different database used")}
	}
	dc, err := tx.grabConn()
	if err != nil {
		return &Stmt{stickyErr: err}
	}
	dc.Lock()
	si, err := dc.ci.Prepare(stmt.query)
	dc.Unlock()
	txs := &Stmt{
		db: tx.db,
		tx: tx,
		txsi: &driverStmt{
			Locker: dc,
			si:     si,
		},
		query:     stmt.query,
		stickyErr: err,
	}
	tx.stmts.Lock()
	tx.stmts.v = append(tx.stmts.v, txs)
	tx.stmts.Unlock()
	return txs
}

// Exec executes a query that doesn't return rows.
// For example: an INSERT and UPDATE.
func (tx *Tx) Exec(query string, args ...interface{}) (Result, error) {
	dc, err := tx.grabConn()
	if err != nil {
		return nil, err
	}

	if execer, ok := dc.ci.(driver.Execer); ok {
		dargs, err := driverArgs(nil, args)
		if err != nil {
			return nil, err
		}
		dc.Lock()
		resi, err := execer.Exec(query, dargs)
		dc.Unlock()
		if err == nil {
			return driverResult{dc, resi}, nil
		}
		if err != driver.ErrSkip {
			return nil, err
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
func (tx *Tx) Query(query string, args ...interface{}) (*Rows, error) {
	dc, err := tx.grabConn()
	if err != nil {
		return nil, err
	}
	releaseConn := func(error) {}
	return tx.db.queryConn(dc, releaseConn, query, args)
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
func (tx *Tx) QueryRow(query string, args ...interface{}) *Row {
	rows, err := tx.Query(query, args...)
	return &Row{rows: rows, err: err}
}
