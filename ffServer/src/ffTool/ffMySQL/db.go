package main

import (
	"database/sql"

	// 每次查询，优先获取可用的数据库连接，如果没有且未达到最大同时连接数，则创建连接，否则等待可用的数据库连接
	"ffCommon/log/log"

	_ "github.com/go-sql-driver/mysql"
)

type stmtInfo struct {
	stmt *sql.Stmt
	str  string // 原始sql语句
}

type myDB struct {
	sqlDB *sql.DB

	stmts map[int]*stmtInfo
}

// Open 初始化操作，仅且必须执行一次
func (db *myDB) Open(args string) error {
	sqlDB, err := sql.Open("mysql", args)
	if err == nil {
		err = sqlDB.Ping()
		if err == nil {
			db.sqlDB = sqlDB
			db.stmts = make(map[int]*stmtInfo, 128)

			sqlDB.SetMaxIdleConns(10)
			sqlDB.SetMaxOpenConns(10)
		}
	}
	return err
}

// Close 关闭操作，仅且必须执行一次
func (db *myDB) Close() {
	if db.sqlDB != nil {
		for _, info := range db.stmts {
			info.stmt.Close()
		}
		db.stmts = make(map[int]*stmtInfo)

		db.sqlDB.Close()
		db.sqlDB = nil
	}
}

// Prepare 不支持并发执行，应在Open操作成功时，连续将所需的所有sql语句，转换为Stmt
func (db *myDB) Prepare(strSQL string) (int, bool) {
	for idSTMT, info := range db.stmts {
		if info.str == strSQL {
			return idSTMT, true
		}
	}

	stmt, err := db.sqlDB.Prepare(strSQL)
	if err != nil {
		log.RunLogger.Printf("myDB.Prepare: Prepare strSQL[%s] get error[%v]", strSQL, err)
		return 0, false
	}

	id := len(db.stmts) + 1
	db.stmts[id] = &stmtInfo{
		stmt: stmt,
		str:  strSQL,
	}
	return id, true
}

// Exec 执行stmt语句，返回影响的行数以及是否成功
// 主要用于非查询操作
// 是否成功，仅表明了语句的执行结果，是否达成了逻辑上的预期，需要同步判定返回的影响行数
func (db *myDB) Exec(idSTMT int, args ...interface{}) (int64, bool) {
	info, ok := db.stmts[idSTMT]
	if !ok {
		log.RunLogger.Printf("myDB.Exec: invalid idSTMT[%d]", idSTMT)
		return 0, false
	}

	result, err := info.stmt.Exec(args...)
	if err != nil {
		log.RunLogger.Printf("myDB.Exec: stmt[%s] Exec failed[%v]", info.str, err)
		return 0, false
	}

	rowsCount, err := result.RowsAffected()
	if err != nil {
		log.RunLogger.Printf("myDB.Exec: stmt[%s] RowsAffected failed[%v]", info.str, err)
		return rowsCount, false
	}
	return rowsCount, true
}

// Query 执行stmt语句，返回查询结果集以及是否成功
// 用于查询操作
// 外界需要手动关闭查询结果集
func (db *myDB) Query(idSTMT int, args ...interface{}) (*sql.Rows, bool) {
	info, ok := db.stmts[idSTMT]
	if !ok {
		log.RunLogger.Printf("myDB.Query: invalid idSTMT[%d]", idSTMT)
		return nil, false
	}

	rows, err := info.stmt.Query(args...)
	if err != nil {
		log.RunLogger.Printf("myDB.Query: stmt[%s] Query failed[%v]", info.str, err)
		return nil, false
	}

	return rows, true
}
