package ffMySQL

import (
	"database/sql"
	"ffCommon/log/log"
	"fmt"
)

// mysqlQueryReuslt 实现接口 ffDef.IDBQueryResult
type mysqlQueryReuslt struct {
	sql        string        // 原始语句
	sqlArgs    []interface{} // sql参数
	sqlResult  error         // sql语句的执行情况
	execResult sql.Result    // insert,update等非select操作的实际执行结果
	queryRows  *sql.Rows     // select的实际执行结果
}

func (qr *mysqlQueryReuslt) SQL() string {
	return fmt.Sprintf("sql[%s] with args[%v]", qr.sql, qr.sqlArgs)
}

func (qr *mysqlQueryReuslt) SQLResult() error {
	return qr.sqlResult
}

func (qr *mysqlQueryReuslt) LastInsertId() (int64, error) {
	if qr.sqlResult != nil {
		return 0, qr.sqlResult
	}
	return qr.execResult.LastInsertId()
}
func (qr *mysqlQueryReuslt) RowsAffected() (int64, error) {
	if qr.sqlResult != nil {
		return 0, qr.sqlResult
	}
	return qr.execResult.RowsAffected()
}

func (qr *mysqlQueryReuslt) Next() bool {
	if qr.sqlResult != nil {
		return false
	}
	return qr.queryRows.Next()
}

func (qr *mysqlQueryReuslt) Err() error {
	if qr.sqlResult != nil {
		return qr.sqlResult
	}
	return qr.queryRows.Err()
}

func (qr *mysqlQueryReuslt) Columns() ([]string, error) {
	if qr.sqlResult != nil {
		return nil, qr.sqlResult
	}
	return qr.queryRows.Columns()
}

func (qr *mysqlQueryReuslt) Scan(dest ...interface{}) error {
	if qr.sqlResult != nil {
		return qr.sqlResult
	}
	return qr.queryRows.Scan(dest...)
}

func (qr *mysqlQueryReuslt) initWithExecResult(sql string, args []interface{}, result sql.Result, sqlResult error) {
	qr.sql, qr.sqlArgs, qr.execResult, qr.sqlResult = sql, args, result, sqlResult
}
func (qr *mysqlQueryReuslt) initWithQueryRows(sql string, args []interface{}, rows *sql.Rows, sqlResult error) {
	qr.sql, qr.sqlArgs, qr.queryRows, qr.sqlResult = sql, args, rows, sqlResult
}
func (qr *mysqlQueryReuslt) back(req *mysqlQueryRequest) {
	defer poolQueryResult.back(qr)

	if req.stmt != nil && req.stmt.isQuerySQL {
		err := qr.queryRows.Close()
		if err != nil {
			log.RunLogger.Printf("mysqlQueryReuslt.back: sql[%s] with args[%v] Close rows get error[%v]", req.stmt.sql, req.args, err)
		}
	}
	qr.queryRows, qr.execResult = nil, nil
}

func newDBQueryResult() interface{} {
	return &mysqlQueryReuslt{}
}
