package tcpserver

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/uuid"
	"sync"
)

const (

	// DefaultTotalServerNetEventDataCount 初始默认创建多少serverNetEventData(供本进程所有tcpServer使用)
	DefaultTotalServerNetEventDataCount = 256

	// DefaultServerNetEventDataChanCount 一个tcpServer.chNetEventDataInner的缓存有多大
	DefaultServerNetEventDataChanCount = 64
)

// 客户端
var mutexServer sync.Mutex
var mapServers = make(map[uuid.UUID]*tcpServer, 1)
var uuidGenerator uuid.Generator

// init 初始tcpserver模块
func init() {
	log.RunLogger.Printf("tcpserver.Init")

	var err error
	uuidGenerator, err = uuid.NewGeneratorSafe(0)
	if err != nil {
		log.FatalLogger.Printf("tcpserver init failed: %v", err)
	}
}

// PrintModule 输出tcpsession模块信息
func PrintModule() {
	runLogger := log.RunLogger
	runLogger.Println("PrintModule Start[tcpserver]:")

	runLogger.Println("tcpserver.mapServers:")
	for _, server := range mapServers {
		runLogger.Printf("tcpserver: %v", server)
	}

	runLogger.Printf("PrintModule End [tcpserver]\n\n")
}

// NewServer 新建一个 base.Server
// 	addr: 监听地址
func NewServer(addr string) (server base.Server, err error) {
	mutexServer.Lock()
	defer mutexServer.Unlock()

	return newServer(addr, uuidGenerator.Gen())
}
