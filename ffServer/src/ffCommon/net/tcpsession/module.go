package tcpsession

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/pool"
	"ffCommon/uuid"
	"fmt"
)

const (

	// DefaultTotalSessionNetEventDataCount 初始默认创建多少sessionNetEventData, 供本进程所有tcpSession使用
	DefaultTotalSessionNetEventDataCount = 256
)

// 没有网络包进入的最大间隔(秒)
var readDeadtime int

// 缓存
var sessPool *sessionPool
var eventDataPool *sessionNetEventDataPool

// Init 初始session模块
// 	_readDeadtime: 没有网络包进入的最大间隔(秒), 超过此时限, 则认为对端关闭了. 必须大于等于30
// 	_onlineCount: 预计多少Session同时连接
// 	_totalSessionNetEventDataCount: 初始创建多少sessionNetEventData, 供本进程所有tcpSession使用
func Init(
	_readDeadtime int,
	_onlineCount int,
	_totalSessionNetEventDataCount int) (err error) {

	log.RunLogger.Printf("tcpsession.Init _readDeadtime[%v] _onlineCount[%v] _totalSessionNetEventDataCount[%v]",
		_readDeadtime, _onlineCount, _totalSessionNetEventDataCount)

	if _readDeadtime < 30 {
		return fmt.Errorf("tcpsession.Init invalid _readDeadtime[%v], must not less than 30", _readDeadtime)
	}

	if _onlineCount < 1 {
		return fmt.Errorf("tcpsession.Init invalid _onlineCount[%v], must not less than 1", _onlineCount)
	}

	if _totalSessionNetEventDataCount < DefaultTotalSessionNetEventDataCount {
		return fmt.Errorf("tcpsession.Init invalid _totalSessionNetEventDataCount[%v], must not less than %v",
			_totalSessionNetEventDataCount, DefaultTotalSessionNetEventDataCount)
	}

	uuidGenerator, err := uuid.NewGeneratorSafe(0)
	if err != nil {
		return err
	}

	readDeadtime = _readDeadtime

	funcCreateSession := func() interface{} {
		return newSession()
	}

	sessPool = &sessionPool{
		pool: pool.New("tcpsession.sessPool", false, funcCreateSession, _onlineCount, 50),

		uuidGenerator: uuidGenerator,
	}

	funcCreateSessionNetEventData := func() interface{} {
		return newSessionNetEventData()
	}

	eventDataPool = &sessionNetEventDataPool{
		pool: pool.New("tcpsession.eventDataPool", false, funcCreateSessionNetEventData, _totalSessionNetEventDataCount, 50),
	}

	return nil
}

// Apply 申请一个空闲session, session将在连接断开后, 自动缓存到sp. 该方法不是多goroutine安全的.
func Apply() (s base.Session) {
	return sessPool.apply()
}

// PrintModule 输出tcpsession模块信息
func PrintModule() {
	runLogger := log.RunLogger
	runLogger.Println("PrintModule Start[tcpsession]:")

	runLogger.Println("tcpsession.sessPool:")
	runLogger.Println(sessPool)

	runLogger.Println("tcpsession.eventDataPool:")
	runLogger.Println(eventDataPool)

	runLogger.Printf("PrintModule End [tcpsession]\n\n")
}
