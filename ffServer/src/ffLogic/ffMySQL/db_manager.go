package ffMySQL

import (
	// mysql driver
	_ "github.com/lexical005/mysql"

	"ffCommon/log/log"
	"ffLogic/ffDef"

	"sync/atomic"
)

const (
	runningOn  = 1
	runningOff = 0
)

// MysqlDBManager 数据库管理, 不支持多协程安全
type MysqlDBManager struct {
	// 所有的数据库
	dbs []*mysqlDB

	// 已完成的数据库操作请求
	completedQuerys chan *mysqlQueryRequest

	// 锁
	running uint32
}

// Open 初始化操作，仅且必须执行一次
func (dbm *MysqlDBManager) Open(tomlPath string) error {
	err := readToml(tomlPath)
	if err != nil {
		return err
	}
	log.RunLogger.Printf("MysqlDBManager.Open: Open[%v] appMysqlConfig:\n%v\n\n", tomlPath, appMysqlConfig)

	poolQueryRequest.init(appMysqlConfig.MaxQueryCount)
	poolQueryResult.init(appMysqlConfig.MaxQueryCount)

	dbm.dbs = make([]*mysqlDB, len(appMysqlConfig.DB), len(appMysqlConfig.DB))
	dbm.completedQuerys = make(chan *mysqlQueryRequest, appMysqlConfig.MaxQueryCount)
	atomic.StoreUint32(&dbm.running, runningOn)

	for _, dbConfig := range appMysqlConfig.DB {
		db := newMysqlDB(dbConfig)
		if err := db.init(); err != nil {
			return err
		}
		dbm.dbs[dbConfig.UniqueID] = db
	}

	return nil
}

// Close 关闭操作，仅且必须执行一次
func (dbm *MysqlDBManager) Close() {
	log.RunLogger.Println("MysqlDBManager.Close")

	atomic.StoreUint32(&dbm.running, runningOff)

	if dbm.dbs != nil {
		for _, db := range dbm.dbs {
			db.close()
		}
		dbm.dbs = nil

		close(dbm.completedQuerys)
		dbm.completedQuerys = nil
	}
}

// Query 返回一个新的 IDBQueryRequest，用于数据库操作。callback和args，在请求者主动取消前或者接收到查询结果前，必须有效！
//  idMysqlDB   int             // 哪个数据库
//  idMysqlStmt int             // 哪个语句
//  callback    DBQueryCallback // 查询结果回调函数
//  args        []interface{}   // 查询参数
func (dbm *MysqlDBManager) Query(idMysqlDB int, idMysqlStmt int, callback ffDef.DBQueryCallback, args ...interface{}) (ffDef.IDBQueryRequest, bool) {
	// 如果已经关闭, 则什么都不做
	if atomic.CompareAndSwapUint32(&dbm.running, runningOff, runningOff) {
		return nil, false
	}

	if idMysqlDB >= len(dbm.dbs) {
		log.FatalLogger.Printf("MysqlDBManager.Query: invalid idMysqlDB[%d] with idMysqlStmt[%d] args[%v]", idMysqlDB, idMysqlStmt, args)
		return nil, false
	}

	db := dbm.dbs[idMysqlDB]
	stmt, ok := db.stmts[idMysqlStmt]
	if !ok {
		log.FatalLogger.Printf("MysqlDBManager.Query: invalid idMysqlStmt[%d] with idMysqlDB[%d] args[%v]", idMysqlStmt, idMysqlDB, args)
		return nil, false
	}

	req := poolQueryRequest.apply()
	req.init(stmt, idMysqlDB, idMysqlStmt, callback, args...)
	return req, true
}

// DispatchResult 分发查询到的结果
func (dbm *MysqlDBManager) DispatchResult() {
	for i := 0; i < appMysqlConfig.MaxQueryCount; i++ {
		// 如果已经关闭, 则什么都不做
		if atomic.CompareAndSwapUint32(&dbm.running, runningOff, runningOff) {
			for req := range dbm.completedQuerys {
				if req != nil {
					req.back()
				}
			}
			break
		}

		select {
		case req := <-dbm.completedQuerys:
			if req != nil {
				req.onResult()
			} else {
				break
			}
		default:
			break
		}
	}
}

// onQueryRequestComplete 数据库操作完成，等待通知请求者
func (dbm *MysqlDBManager) onQueryRequestComplete(req *mysqlQueryRequest) {
	// 如果已经关闭, 则什么都不做
	if atomic.CompareAndSwapUint32(&dbm.running, runningOff, runningOff) {
		req.back()
		return
	}

	dbm.completedQuerys <- req
}
