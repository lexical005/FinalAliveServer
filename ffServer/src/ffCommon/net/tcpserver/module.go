package tcpserver

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/uuid"
	"sync"
)

// 缓存
var eventDataPool *serverNetEventDataPool

// 客户端
var mutexServer sync.Mutex
var mapServers = make(map[uuid.UUID]*tcpServer, 1)
var uuidGenerator uuid.Generator

// Init 初始tcpserver模块
func Init() (err error) {
	log.RunLogger.Printf("tcpserver.Init")

	uuidGenerator, err = uuid.NewGeneratorSafe(0)
	if err != nil {
		return err
	}

	return
}

// NewServer 新建一个 base.Server
// 	addr: 监听地址
func NewServer(addr string) (server base.Server, err error) {
	mutexServer.Lock()
	defer mutexServer.Unlock()

	return newServer(addr, uuidGenerator.Gen())
}
