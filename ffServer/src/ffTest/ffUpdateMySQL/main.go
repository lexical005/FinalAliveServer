package main

import (
	"encoding/json"
	"ffCommon/log/log"
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

var mysql = &mysqlManager{}

func doUpdateWeaponComponents(_id int, components string) {
	log.RunLogger.Printf("doUpdateWeaponComponents _id:%d components:%s", _id, components)
	callback := func(result ffDef.IDBQueryResult) {
		err := result.SQLResult()
		if err != nil {
			log.RunLogger.Printf("doUpdateWeaponComponents %s excute get error[%v]", result.SQL(), err)
		} else {
			count, err := result.RowsAffected()
			if err != nil {
				log.RunLogger.Printf("doUpdateWeaponComponents %s RowsAffected get error[%v]", result.SQL(), err)
			} else if count == 0 {
				log.RunLogger.Printf("doUpdateWeaponComponents %s RowsAffected count zero", result.SQL())
			}
		}
		exitLock--
		if exitLock == 0 {
			cWait <- struct{}{}
		}
	}
	mysql.query(0, 100, callback, components, _id)
}

func doSelectAllWeaponComponents() {
	callback := func(result ffDef.IDBQueryResult) {
		err := result.SQLResult()
		if err != nil {
			log.RunLogger.Printf("doSelectAllWeaponComponents %s excute get error[%v]", result.SQL(), err)
		} else {
			for result.Next() {
				var _id int
				var components string
				if err = result.Scan(&_id, &components); err != nil {
					log.RunLogger.Printf("doSelectAllWeaponComponents %s Scan get error[%v]", result.SQL(), err)
				} else {
					var ary []int
					if err = json.Unmarshal([]byte(components), &ary); err != nil {
						log.RunLogger.Printf("doSelectAllWeaponComponents %s Unmarshal %s get error[%v]", result.SQL(), components, err)
					} else {
						needUpdate := false
						for index := 0; index < len(ary); index++ {
							if ary[index] == -1 {
								ary[index] = 0
								needUpdate = true
							}
						}

						if needUpdate {
							if co, err := json.Marshal(ary); err != nil {
								log.RunLogger.Printf("doSelectAllWeaponComponents %s Marshal %v get error[%v]", result.SQL(), ary, err)
							} else {
								doUpdateWeaponComponents(_id, string(co))
								exitLock++
							}
						}
					}
				}
			}

			if exitLock == 0 {
				cWait <- struct{}{}
			}
		}
	}
	mysql.query(0, 101, callback)
}

var cWait = make(chan struct{}, 1)
var exitLock int

func main() {
	// 异常保护
	defer util.PanicProtect()

	// 数据库配置
	mysql.start()

	doSelectAllWeaponComponents()

	<-cWait
}
