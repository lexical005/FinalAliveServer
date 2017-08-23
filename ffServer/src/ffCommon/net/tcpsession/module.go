package tcpsession

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/pool"
	"ffCommon/uuid"
	"net"
)

// 配置
var sessionConfig base.SessionConfig

// 缓存
var sessPool *sessionPool
var eventDataPool *sessionNetEventDataPool

// Init 初始session模块
func Init(config *base.SessionConfig) (err error) {
	log.RunLogger.Printf("tcpsession.Init config[%v]", config)

	sessionConfig = *config

	uuidGenerator, err := uuid.NewGeneratorSafe(0)
	if err != nil {
		return err
	}

	funcCreateSession := func() interface{} {
		return newSession()
	}

	sessPool = &sessionPool{
		pool: pool.New("tcpsession.sessPool", false, funcCreateSession, sessionConfig.InitOnlineCount, 50),

		uuidGenerator: uuidGenerator,
	}

	funcCreateSessionNetEventData := func() interface{} {
		return newSessionNetEvent()
	}

	eventDataPool = &sessionNetEventDataPool{
		pool: pool.New("tcpsession.eventDataPool", false, funcCreateSessionNetEventData, sessionConfig.InitNetEventDataCount, 50),
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
