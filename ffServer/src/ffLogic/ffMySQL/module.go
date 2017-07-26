package ffMySQL

var appMysqlConfig = &mysqlConfig{}

var poolQueryRequest = &mysqlQueryRequestPool{}

var poolQueryResult = &mysqlQueryResultPool{}

var dbMgr = &MysqlDBManager{}

// NewMysqlManager 返回数据库管理器
func NewMysqlManager() *MysqlDBManager {
	return dbMgr
}
