package tcpsession

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/pool"
	"ffCommon/uuid"
	"fmt"
)

// 没有网络包进入的最大间隔(秒)
var readDeadtime int

// Session缓存
var sessPool *sessionPool
var eventDataPool *sessionNetEventDataPool

// Init 初始session模块
// _readDeadtime: 没有网络包进入的最大间隔(秒)，超过此时限，则认为对端关闭了
// _onlineCount: 预计多少Session同时连接
func Init(_readDeadtime int, _onlineCount int) (err error) {
	if _readDeadtime < 30 {
		return fmt.Errorf("tcpsession.Init: invalid _readDeadtime[%v]", _readDeadtime)
	}

	if _onlineCount < 1 {
		return fmt.Errorf("tcpsession.Init: invalid _onlineCount[%v]", _onlineCount)
	}

	uuidGenerator, err := uuid.NewGenerator(0)
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
		pool: pool.New("tcpsession.eventDataPool", true, funcCreateSessionNetEventData, _onlineCount/2, 50),
	}

	return nil
}

// Apply 申请一个空闲session, session将在连接断开后, 自动缓存到sp. 该方法不是多goroutine安全的.
func Apply() (s base.Session) {
	return sessPool.apply()
}

// PrintModule 输出session模块信息
func PrintModule() {
	runLogger := log.RunLogger
	runLogger.Println("PrintModule Start[tcpsession]:")

	runLogger.Println("tcpsession.sessPool:")
	runLogger.Println(sessPool)

	runLogger.Println("tcpsession.eventDataPool:")
	runLogger.Println(eventDataPool)

	runLogger.Printf("PrintModule End [tcpsession]\n\n")
}
