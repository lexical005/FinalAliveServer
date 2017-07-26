package ffMySQL

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"fmt"

	"github.com/lexical005/toml"
)

type sqlConfig struct {
	Number int    // 在SQL语句中全局唯一的编号，目前定义规则为 (所属组的DBConn+1)*100+从0开始递增的编号
	SQL    string // SQL语句
}
type sqlGroupConfig struct {
	DBConn int          // 使用的数据库的连接
	SQL    []*sqlConfig // 所有SQL配置
}
type dbConfig struct {
	DataBase string // 数据库名称
	Address  string // 数据库所在地址
	Account  string // 连接数据库时使用的帐号
	Password string // 连接数据库时使用的帐号的密码

	UniqueID int   // 数据库唯一编号，全局唯一，从0开始递增
	DBConns  []int // 与数据库的连接，同一连接上的操作，会确保时序性，而不同连接之间的时序性不作任何保证，不能重复，从0开始递增

	SQL map[string]sqlGroupConfig // 该数据库上的所有预处理语句. key: 组别, 一般为表名; value: sql语句定义
}

func (db *dbConfig) String() string {
	return fmt.Sprintf("DataBase:%s Address:%s Account:%s",
		db.DataBase, db.Address, db.Account)
}

type mysqlConfig struct {
	MaxQueryCount int // 同时最大数据库操作数量

	DB []*dbConfig
}

func (db *mysqlConfig) String() string {
	s := fmt.Sprintf("MaxQueryCount:%d", db.MaxQueryCount)
	for _, one := range db.DB {
		s += fmt.Sprintf("%s\n%v", s, one)
	}
	return s
}

func (db *mysqlConfig) check() bool {
	result := true
	for index, dbConfig := range db.DB {
		if dbConfig.UniqueID != index {
			result = false
			log.RunLogger.Printf("dbConfig[%s:%d] UniqueID must start with 0 and increase in order by step 1",
				dbConfig.DataBase, dbConfig.UniqueID)
		}

		dbConnCount := len(dbConfig.DBConns)
		if dbConnCount != dbConfig.DBConns[dbConnCount-1]+1 {
			result = false
			log.RunLogger.Printf("dbConfig[%s:%d] DBConns must start with 0 and increase in order by step 1",
				dbConfig.DataBase, dbConfig.UniqueID)
		}

		globalNumbers := make(map[int]struct{}, 128)
		for groupName, groupConfig := range dbConfig.SQL {
			numberGroup := -1
			// DBConn
			if groupConfig.DBConn < 0 || groupConfig.DBConn >= dbConnCount {
				result = false
				log.RunLogger.Printf("dbConfig[%s:%d] group[%s] invalid groupConfig.DBConn:%d",
					dbConfig.DataBase, dbConfig.UniqueID, groupName, groupConfig.DBConn)
			}

			for _, sql := range groupConfig.SQL {
				// Number
				if numberGroup == -1 {
					numberGroup = sql.Number / 100
				}
				if numberGroup != sql.Number/100 {
					result = false
					log.RunLogger.Printf("dbConfig[%s:%d] group[%s] invalid SQL.Number:%d",
						dbConfig.DataBase, dbConfig.UniqueID, groupName, sql.Number)
				} else if _, ok := globalNumbers[sql.Number]; ok {
					result = false
					log.RunLogger.Printf("dbConfig[%s:%d] group[%s] multi SQL.Number:%d",
						dbConfig.DataBase, dbConfig.UniqueID, groupName, sql.Number)
				}
				globalNumbers[sql.Number] = struct{}{}
			}
		}
	}
	return result
}

func readToml(tomlPath string) error {
	// 读取文件内容
	fileContent, err := util.ReadFile(tomlPath)
	if err != nil {
		return err
	}

	// 解析
	err = toml.Unmarshal(fileContent, appMysqlConfig)
	if err != nil {
		return err
	}

	// 数据有效性检查
	if !appMysqlConfig.check() {
		return fmt.Errorf("readToml: invalid toml config")
	}

	return nil
}
