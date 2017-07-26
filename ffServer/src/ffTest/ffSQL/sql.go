package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
)

// nowFunc returns the current time; it's overridden in tests.
var nowFunc = time.Now

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver driver.Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("sql: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("sql: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	// For tests.
	drivers = make(map[string]driver.Driver)
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

// RawBytes is a byte slice that holds a reference to memory owned by
// the database itself. After a Scan into a RawBytes, the slice is only
// valid until the next call to Next, Scan, or Close.
type RawBytes []byte

// NullString represents a string that may be null.
// NullString implements the Scanner interface so
// it can be used as a scan destination:
//
//  var s NullString
//  err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
//  ...
//  if s.Valid {
//     // use s.String
//  } else {
//     // NULL value
//  }
//
type NullString struct {
	String string
	Valid  bool // Valid is true if String is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return convertAssign(&ns.String, value)
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// NullInt64 represents an int64 that may be null.
// NullInt64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullInt64) Scan(value interface{}) error {
	if value == nil {
		n.Int64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Int64, value)
}

// Value implements the driver Valuer interface.
func (n NullInt64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Int64, nil
}

// NullFloat64 represents a float64 that may be null.
// NullFloat64 implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullFloat64 struct {
	Float64 float64
	Valid   bool // Valid is true if Float64 is not NULL
}

// Scan implements the Scanner interface.
func (n *NullFloat64) Scan(value interface{}) error {
	if value == nil {
		n.Float64, n.Valid = 0, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Float64, value)
}

// Value implements the driver Valuer interface.
func (n NullFloat64) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Float64, nil
}

// NullBool represents a bool that may be null.
// NullBool implements the Scanner interface so
// it can be used as a scan destination, similar to NullString.
type NullBool struct {
	Bool  bool
	Valid bool // Valid is true if Bool is not NULL
}

// Scan implements the Scanner interface.
func (n *NullBool) Scan(value interface{}) error {
	if value == nil {
		n.Bool, n.Valid = false, false
		return nil
	}
	n.Valid = true
	return convertAssign(&n.Bool, value)
}

// Value implements the driver Valuer interface.
func (n NullBool) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Bool, nil
}

// Scanner is an interface used by Scan.
type Scanner interface {
	// Scan assigns a value from a database driver.
	//
	// The src value will be of one of the following types:
	//
	//    int64
	//    float64
	//    bool
	//    []byte
	//    string
	//    time.Time
	//    nil - for NULL values
	//
	// An error should be returned if the value cannot be stored
	// without loss of information.
	Scan(src interface{}) error
}

// ErrNoRows is returned by Scan when QueryRow doesn't return a
// row. In such a case, QueryRow returns a placeholder *Row value that
// defers this error until a Scan.
var ErrNoRows = errors.New("sql: no rows in result set")

// connReuseStrategy determines how (*DB).conn returns database connections.
type connReuseStrategy uint8

const (
	// alwaysNewConn forces a new connection to the database.
	alwaysNewConn connReuseStrategy = iota
	// cachedOrNewConn returns a cached connection, if available, else waits
	// for one to become available (if MaxOpenConns has been reached) or
	// creates a new database connection.
	cachedOrNewConn
)

// driverStmt associates a driver.Stmt with the
// *driverConn from which it came, so the driverConn's lock can be
// held during calls.
type driverStmt struct {
	sync.Locker // the *driverConn
	si          driver.Stmt
}

func (ds *driverStmt) Close() error {
	ds.Lock()
	defer ds.Unlock()
	return ds.si.Close()
}

// depSet is a finalCloser's outstanding dependencies
type depSet map[interface{}]bool // set of true bools

// The finalCloser interface is used by (*DB).addDep and related
// dependency reference counting.
type finalCloser interface {
	// finalClose is called when the reference count of an object
	// goes to zero. (*DB).mu is not held while calling it.
	finalClose() error
}

// This is the size of the connectionOpener request chan (DB.openerCh).
// This value should be larger than the maximum typical value
// used for db.maxOpen. If maxOpen is significantly larger than
// connectionRequestQueueSize then it is possible for ALL calls into the *DB
// to block until the connectionOpener can satisfy the backlog of requests.
var connectionRequestQueueSize = 1000000

// Open opens a database specified by its database driver name and a
// driver-specific data source name, usually consisting of at least a
// database name and connection information.
//
// Most users will open a database via a driver-specific connection
// helper function that returns a *DB. No database drivers are included
// in the Go standard library. See https://golang.org/s/sqldrivers for
// a list of third-party drivers.
//
// Open may just validate its arguments without creating a connection
// to the database. To verify that the data source name is valid, call
// Ping.
//
// The returned DB is safe for concurrent use by multiple goroutines
// and maintains its own pool of idle connections. Thus, the Open
// function should be called just once. It is rarely necessary to
// close a DB.
func Open(driverName, dataSourceName string) (*DB, error) {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("sql: unknown driver %q (forgotten import?)", driverName)
	}
	db := &DB{
		driver:   driveri,
		dsn:      dataSourceName,
		openerCh: make(chan struct{}, connectionRequestQueueSize),
		lastPut:  make(map[*driverConn]string),
	}
	go db.connectionOpener()
	return db, nil
}

const defaultMaxIdleConns = 2

// DBStats contains database statistics.
type DBStats struct {
	// OpenConnections is the number of open connections to the database.
	OpenConnections int
}

// connRequest represents one request for a new connection
// When there are no idle connections available, DB.conn will create
// a new connRequest and put it on the db.connRequests list.
type connRequest struct {
	conn *driverConn
	err  error
}

var errDBClosed = errors.New("sql: database is closed")

// putConnHook is a hook for testing.
var putConnHook func(*DB, *driverConn)

// debugGetPut determines whether getConn & putConn calls' stack traces
// are returned for more verbose crashes.
const debugGetPut = false

// maxBadConnRetries is the number of maximum retries if the driver returns
// driver.ErrBadConn to signal a broken connection before forcing a new
// connection to be opened.
const maxBadConnRetries = 2
