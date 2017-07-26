package main

import (
	"ffCommon/log/log"
	"ffLogic/ffDef"
	"ffLogic/ffMySQL"
	"time"
)

var dbm *ffMySQL.MysqlDBManager

func testMYSQL() {

	dbm = ffMySQL.NewMysqlManager()
	err := dbm.Open("sql.toml")
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	callback2 := func(result2 ffDef.IDBQueryResult) {
		err := result2.SQLResult()
		if err != nil {
			log.RunLogger.Printf("%s excute get error[%v]", result2.SQL(), err)
		} else {
			count, err := result2.RowsAffected()
			if err != nil {
				log.RunLogger.Printf("%s RowsAffected get error[%v]", result2.SQL(), err)
			} else if count == 0 {
				log.RunLogger.Printf("%s RowsAffected count zero", result2.SQL())
			} else {
				log.RunLogger.Printf("%s success", result2.SQL())
			}
		}
	}

	callback1 := func(result1 ffDef.IDBQueryResult) {
		for result1.Next() {
			var signature, storeOrder, storeAmount, vivoOrder, vivoFee string
			var status, vivoPayType int
			var target = []interface{}{&signature, &status, &storeOrder, &storeAmount, &vivoOrder, &vivoPayType, &vivoFee}
			err := result1.Scan(target...)
			if err != nil {
				log.RunLogger.Println("callback1 Scan get error", err)
			} else {
				log.RunLogger.Println("callback1 Scan success", signature, status, storeOrder, storeAmount, vivoOrder, vivoPayType, vivoFee)
			}
		}

		if query, ok := dbm.Query(0, 101, callback2, "148247180307597123268", "1bc2f86cc8d24715db927a449d4ea12b", "152fe331-6ba7-460c-bd1a-9e0852f763fc"); ok {
			log.RunLogger.Println("start callback2 success")
			query.Query()
		} else {
			log.RunLogger.Println("start callback2 failed")
		}
	}

	callback0 := func(result0 ffDef.IDBQueryResult) {
		err := result0.SQLResult()
		if err != nil {
			log.RunLogger.Printf("%s excute get error[%v]", result0.SQL(), err)
		} else {
			count, err := result0.RowsAffected()
			if err != nil {
				log.RunLogger.Printf("%s RowsAffected get error[%v]", result0.SQL(), err)
			} else if count == 0 {
				log.RunLogger.Printf("%s RowsAffected count zero", result0.SQL())
			} else {
				log.RunLogger.Printf("%s success", result0.SQL())
			}
		}

		if query, ok := dbm.Query(0, 103, callback1, "152fe331-6ba7-460c-bd1a-9e0852f763fc"); ok {
			query.Query()
		}
	}

	if query, ok := dbm.Query(0, 100, callback0, "152fe331-6ba7-460c-bd1a-9e0852f763fc", "1bc2f86cc8d24715db927a449d4ea12b", "1", "100.00"); ok {
		query.Query()
	}

	c := make(chan struct{}, 1)
	go mainDispatch(c)

	select {
	case <-time.After(time.Second * 5):
		log.RunLogger.Println("time over")
	}

	dbm.Close()
}

func mainDispatch(c chan struct{}) {
	for {
		select {
		case <-time.After(time.Millisecond * 10):
			dbm.DispatchResult()
		case <-c:
			break
		}
	}
}
