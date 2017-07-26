package ffMySQL

import (
	"database/sql"
	"ffCommon/log/log"
	"strings"
)

type mysqlStmt struct {
	stmt       *sql.Stmt
	dbConn     *mysqlConn // 使用的数据库连接
	sql        string     // 原始SQL语句
	isQuerySQL bool       // 是不是查询SQL语句
}

func (stmt *mysqlStmt) query(req *mysqlQueryRequest, careResult bool) (queryResult *mysqlQueryReuslt) {
	if careResult {
		queryResult = poolQueryResult.apply()
	}

	if stmt.isQuerySQL {
		rows, err := stmt.stmt.Query(req.args...)
		if err == nil {
			if careResult {
				queryResult.initWithQueryRows(stmt.sql, req.args, rows, nil)
			} else {
				err = rows.Close()
				if err != nil {
					log.RunLogger.Printf("mysqlStmt.query: sql[%s] with args[%v] Close rows get error[%v]", stmt.sql, req.args, err)
				}
			}
		} else if careResult {
			queryResult.initWithQueryRows(stmt.sql, req.args, nil, err)
		}
		return
	}

	result, err := stmt.stmt.Exec(req.args...)
	if err == nil {
		if careResult {
			queryResult.initWithExecResult(stmt.sql, req.args, result, nil)
		}
	} else if careResult {
		queryResult.initWithQueryRows(stmt.sql, req.args, nil, err)
	}

	return
}

func newMysqlStmt(stmt *sql.Stmt, dbConn *mysqlConn, sql string) *mysqlStmt {
	return &mysqlStmt{
		dbConn:     dbConn,
		stmt:       stmt,
		sql:        sql,
		isQuerySQL: strings.HasPrefix(strings.ToUpper(sql), "SELECT "),
	}
}
