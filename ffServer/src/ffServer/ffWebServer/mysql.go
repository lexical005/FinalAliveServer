package main

import (
	"ffCommon/util"
	"ffLogic/ffDef"
	"ffLogic/ffMySQL"
	"sync"
	"time"
)

type mysqlManager struct {
	dbm *ffMySQL.MysqlDBManager

	rwMutex sync.RWMutex
}

func (mm *mysqlManager) start() {
	mm.dbm = ffMySQL.NewMysqlManager()
	mm.dbm.Open("toml/sql.toml")

	go util.SafeGo(mm.dispatchMySQL, nil)
}

func (mm *mysqlManager) close() {
	mm.rwMutex.Lock()
	defer mm.rwMutex.Unlock()

	mm.dbm.Close()
	mm.dbm = nil
}

func (mm *mysqlManager) query(idMysqlDB int, idMysqlStmt int, callback ffDef.DBQueryCallback, args ...interface{}) bool {
	mm.rwMutex.RLock()
	defer mm.rwMutex.RUnlock()

	if mm.dbm != nil {
		query, ok := mm.dbm.Query(idMysqlDB, idMysqlStmt, callback, args...)
		if !ok {
			return false
		}
		query.Query()
		return true
	}
	return false
}

func (mm *mysqlManager) dispatchMySQL(params ...interface{}) {
	for {
		select {
		case <-time.After(time.Millisecond * 10):
			if mm.doDispatch() {
				break
			}
		}
	}
}

func (mm *mysqlManager) doDispatch() bool {
	mm.rwMutex.RLock()
	defer mm.rwMutex.RUnlock()

	if mm.dbm != nil {
		mm.dbm.DispatchResult()
		return false
	}
	return true
}
