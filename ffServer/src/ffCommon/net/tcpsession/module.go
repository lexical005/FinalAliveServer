package tcpsession

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/uuid"
	"net"
)

// 配置
var sessionConfig base.SessionConfig

// session uuid 生成器
var uuidSessionGenerator uuid.Generator

// Init 初始session模块
func Init(config *base.SessionConfig, _uuidSessionGenerator uuid.Generator) (err error) {
	log.RunLogger.Printf("tcpsession.Init config[%v]", config)

	sessionConfig = *config

	uuidSessionGenerator = _uuidSessionGenerator

	return nil
}

// Apply 申请一个空闲session, session将在连接断开后, 自动缓存到sp. 该方法不是多goroutine安全的.
func Apply(conn net.Conn) (s base.Session) {
	return newSession(conn, uuidSessionGenerator.Gen())
}
