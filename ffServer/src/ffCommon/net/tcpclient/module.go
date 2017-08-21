package tcpclient

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/uuid"
	"sync"
)

// 客户端
var mutexClient sync.Mutex
var mapClients = make(map[uuid.UUID]*tcpClient, 1)
var uuidGenerator uuid.Generator

func init() {
	log.RunLogger.Printf("tcpclient.Init")

	var err error
	uuidGenerator, err = uuid.NewGeneratorSafe(0)
	if err != nil {
		log.FatalLogger.Printf("tcpserver init failed: %v", err)
	}

	return
}

// PrintModule 输出tcpsession模块信息
func PrintModule() {
	runLogger := log.RunLogger
	runLogger.Println("PrintModule Start[tcpclient]:")

	runLogger.Println("tcpclient.mapClients:")
	for _, client := range mapClients {
		runLogger.Printf("tcpclient: %v", client)
	}

	runLogger.Printf("PrintModule End [tcpclient]\n\n")
}

// NewClient 创建一个base.Client
//	addr: 监听地址
func NewClient(addr string) (c base.Client, err error) {
	mutexClient.Lock()
	defer mutexClient.Unlock()

	return newClient(addr, uuidGenerator.Gen())
}
