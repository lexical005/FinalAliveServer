package tcpsession

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/pool"
	"ffCommon/uuid"
	"fmt"
	"net"
)

const (
	// DefaultReadDeadTime 没有网络包进入的默认最大间隔(秒), 超过此时限, 则认为对端关闭了
	DefaultReadDeadTime = 60

	// DefaultOnlineCount 默认多少Session同时连接
	DefaultOnlineCount = 60

	// DefaultInitSessionNetEventDataCount 初始默认创建多少sessionNetEventData, 供本进程所有tcpSession使用
	DefaultInitSessionNetEventDataCount = 2
)

// 没有网络包进入的最大间隔(秒)
var readDeadtime int

// 缓存
var sessPool *sessionPool
var eventDataPool *sessionNetEventDataPool

// Init 初始session模块
// 	_readDeadtime: 没有网络包进入的最大间隔(秒), 超过此时限, 则认为对端关闭了. >= DefaultReadDeadTime
// 	_onlineCount: 预计多少Session同时连接. >= 1
// 	_initSessionNetEventDataCount: 初始创建多少sessionNetEventData, 供本进程所有tcpSession使用. >=DefaultInitSessionNetEventDataCount
func Init(
	_readDeadtime int,
	_onlineCount int,
	_initSessionNetEventDataCount int) (err error) {

	log.RunLogger.Printf("tcpsession.Init _readDeadtime[%v] _onlineCount[%v] _initSessionNetEventDataCount[%v]",
		_readDeadtime, _onlineCount, _initSessionNetEventDataCount)

	if _readDeadtime < DefaultReadDeadTime {
		return fmt.Errorf("tcpsession.Init invalid _readDeadtime[%v], must not less than %v",
			_readDeadtime, DefaultReadDeadTime)
	}

	if _onlineCount < 1 {
		return fmt.Errorf("tcpsession.Init invalid _onlineCount[%v], must not less than 1",
			_onlineCount)
	}

	if _initSessionNetEventDataCount < DefaultInitSessionNetEventDataCount {
		return fmt.Errorf("tcpsession.Init invalid _initSessionNetEventDataCount[%v], must not less than %v",
			_initSessionNetEventDataCount, DefaultInitSessionNetEventDataCount)
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
		return newSessionNetEvent()
	}

	eventDataPool = &sessionNetEventDataPool{
		pool: pool.New("tcpsession.eventDataPool", false, funcCreateSessionNetEventData, _initSessionNetEventDataCount, 50),
	}

	return nil
}

// Apply 申请一个空闲session, session将在连接断开后, 自动缓存到sp. 该方法不是多goroutine安全的.
func Apply(conn net.Conn) (s base.Session) {
	sess := sessPool.apply()
	sess.setConn(conn)
	return sess
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
