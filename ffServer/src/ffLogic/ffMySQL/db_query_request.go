package ffMySQL

import (
	"ffCommon/log/log"
	"ffLogic/ffDef"
	"sync"
)

// mysqlQueryRequest 实现了接口 IDBQueryRequest
type mysqlQueryRequest struct {
	muValid sync.Mutex // 是否有效的保护锁
	valid   bool       // 本次请求是否还有效

	idMysqlDB   int                   // 哪个数据库
	idMysqlStmt int                   // 哪个语句
	callback    ffDef.DBQueryCallback // 查询结果回调函数
	args        []interface{}         // 查询参数

	stmt   *mysqlStmt        // 预处理语句
	result *mysqlQueryReuslt // 查询结果
}

// IsCancel 是否已经被请求者取消
func (req *mysqlQueryRequest) IsValid() bool {
	return req.valid
}

// Cancel 请求者主动取消，执行此操作后，请求者不能再持有此实例
// 由于数据库操作是异步进行的，一旦在执行Cancel操作时，数据库操作正在锁定输入参数状态，将导致此操作被阻塞，直到数据库操作结束
func (req *mysqlQueryRequest) Cancel() {
	// 日志记录
	log.RunLogger.Printf("mysqlQueryRequest.Cancel: idMysqlDB[%d] idMysqlStmt[%d] args[%v]", req.idMysqlDB, req.idMysqlStmt, req.args)

	defer req.muValid.Unlock()
	req.muValid.Lock()
	req.clear()
}

// Query 请求者请求执行数据库操作
func (req *mysqlQueryRequest) Query() {
	req.stmt.dbConn.addQuery(req)
}

func (req *mysqlQueryRequest) init(stmt *mysqlStmt, idMysqlDB int, idMysqlStmt int, callback ffDef.DBQueryCallback, args ...interface{}) {
	req.stmt, req.result = stmt, nil
	req.idMysqlDB, req.idMysqlDB, req.callback, req.args = idMysqlDB, idMysqlDB, callback, args
	req.valid = true
}

func (req *mysqlQueryRequest) back() {
	req.clear()

	poolQueryRequest.back(req)
}

func (req *mysqlQueryRequest) clear() {
	req.valid = false

	if req.result != nil {
		req.result.back(req)
		req.result = nil
	}

	req.stmt = nil
	req.callback, req.args = nil, nil
}

func (req *mysqlQueryRequest) doQuery() {
	defer req.muValid.Unlock()
	req.muValid.Lock()

	// 请求者已经主动取消，直接回收即可
	if !req.IsValid() {
		req.back()
		return
	}

	// 请求者是否关注数据库操作结果
	careResult := req.callback != nil

	// 执行数据库操作
	result := req.stmt.query(req, careResult)

	// 请求者不关注数据库操作结果时，直接回收即可
	if !careResult {
		req.back()
		return
	}

	// 记录数据库操作结果，添加到待处理列表内
	req.result = result
	dbMgr.onQueryRequestComplete(req)
}

// onResult 通知查询请求者
func (req *mysqlQueryRequest) onResult() {
	defer req.back()

	// 逻辑处理
	if req.IsValid() && req.callback != nil {
		req.callback(req.result)
	}
}

func newDBQueryRequest() interface{} {
	return &mysqlQueryRequest{}
}
