package ffMySQL

import (
	"database/sql"
	"ffCommon/log/log"
	"ffCommon/util"
	"fmt"
)

type mysqlDB struct {
	dbConfig *dbConfig
	sqlDB    *sql.DB

	dbConns []*mysqlConn
	stmts   map[int]*mysqlStmt
}

// init 初始化操作，仅且必须执行一次
func (db *mysqlDB) init() error {
	// root:root@tcp(127.0.0.1:13429)/ff_game
	args := fmt.Sprintf("%s:%s@tcp(%s)/%s", db.dbConfig.Account, db.dbConfig.Password, db.dbConfig.Address, db.dbConfig.DataBase)

	mysqlDB, err := sql.Open("mysql", args)
	if err != nil {
		return err
	}

	err = mysqlDB.Ping()
	if err != nil {
		return err
	}

	db.initConn()

	log.RunLogger.Printf("mysqlDB init %s\n", args)

	if !db.initPrepare(mysqlDB) {
		return fmt.Errorf("mysqlDB.initPrepare: Prepare sql failed")
	}
	log.RunLogger.Printf("")

	// 初始化成功
	db.sqlDB = mysqlDB

	mysqlDB.SetMaxIdleConns(len(db.dbConns))
	mysqlDB.SetMaxOpenConns(len(db.dbConns))

	return nil
}

// initConn 与数据库的连接
func (db *mysqlDB) initConn() {
	connCount := len(db.dbConfig.DBConns)
	db.dbConns = make([]*mysqlConn, connCount, connCount)
	db.stmts = make(map[int]*mysqlStmt, len(db.dbConfig.SQL))

	for index := 0; index < connCount; index++ {
		db.dbConns[index] = &mysqlConn{
			valid:          true,
			chClose:        make(chan struct{}),
			chQueryRequest: make(chan *mysqlQueryRequest, appMysqlConfig.MaxQueryCount/connCount),
		}

		go util.SafeGo(db.dbConns[index].queryLoop, db.dbConns[index].queryLoopEnd)
	}
}

// initPrepare 预处理SQL语句
func (db *mysqlDB) initPrepare(mysqlDB *sql.DB) bool {
	result := true
	for groupName, groupConfig := range db.dbConfig.SQL {
		for _, sql := range groupConfig.SQL {
			stmt, err := mysqlDB.Prepare(sql.SQL)
			if err != nil {
				log.FatalLogger.Printf("mysqlDB.Prepare: Prepare group[%s] number[%d] sql[%s] get error[%v]",
					groupName, sql.Number, sql.SQL, err)
				result = false
				continue
			}

			log.RunLogger.Printf("mysqlDB.Prepare: Prepare group[%s] number[%d] sql[%s]",
				groupName, sql.Number, sql.SQL)
			db.stmts[sql.Number] = newMysqlStmt(stmt, db.dbConns[groupConfig.DBConn], sql.SQL)
		}
	}
	return result
}

// close 关闭操作，仅且必须执行一次
func (db *mysqlDB) close() {
	if db.sqlDB != nil {
		for _, info := range db.stmts {
			info.stmt.Close()
		}

		for _, conn := range db.dbConns {
			conn.close()
		}

		db.sqlDB.Close()
		db.sqlDB, db.stmts = nil, nil
	}
}

func newMysqlDB(dbConfig *dbConfig) *mysqlDB {
	return &mysqlDB{
		dbConfig: dbConfig,
	}
}
