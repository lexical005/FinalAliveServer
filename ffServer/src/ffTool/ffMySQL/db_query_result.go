package main

import "database/sql"

type myDBQueryResult struct {
	sqlStr string

	isQueryRows bool
	execResult  sql.Result
	queryRows   *sql.Rows
}

func (qr *myDBQueryResult) SQL() string {
	return qr.sqlStr
}

func (qr *myDBQueryResult) Close() error {
	if qr.isQueryRows {
		return qr.queryRows.Close()
	}
	return nil
}

func (qr *myDBQueryResult) LastInsertId() (int64, error) {
	return qr.execResult.LastInsertId()
}
func (qr *myDBQueryResult) RowsAffected() (int64, error) {
	return qr.execResult.RowsAffected()
}

func (qr *myDBQueryResult) Next() bool {
	return qr.queryRows.Next()
}

func (qr *myDBQueryResult) Err() error {
	return qr.queryRows.Err()
}

func (qr *myDBQueryResult) Columns() ([]string, error) {
	return qr.queryRows.Columns()
}

func (qr *myDBQueryResult) Scan(dest ...interface{}) error {
	return qr.queryRows.Scan()
}
