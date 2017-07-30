package tcpclient

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/pool"
	"ffCommon/uuid"
	"sync"
)

const (
	// DefaultTotalClientNetEventDataCount 初始默认创建多少clientNetEventData(供本进程所有tcpClient使用)
	DefaultTotalClientNetEventDataCount = 32

	// DefaultClientNetEventDataChanCount 一个tcpClient.chNetEventDataInner的缓存有多大
	DefaultClientNetEventDataChanCount = 8
)

// 客户端
var mutexClient sync.Mutex
var mapClients = make(map[uuid.UUID]*tcpClient, 1)
var uuidGenerator uuid.Generator

// eventDataPool clientNetEventData Pool
var eventDataPool *clientNetEventDataPool

// Init 初始tcpclient模块
func Init() (err error) {
	uuidGenerator, err = uuid.NewGeneratorSafe(0)
	if err != nil {
		return err
	}

	funcCreateClientNetEventData := func() interface{} {
		return newClientNetEventData()
	}

	eventDataPool = &clientNetEventDataPool{
		pool: pool.New("tcpclient.eventDataPool", false, funcCreateClientNetEventData, DefaultTotalClientNetEventDataCount, 50),
	}

	return
}

// PrintModule 输出tcpsession模块信息
func PrintModule() {
	runLogger := log.RunLogger
	runLogger.Println("PrintModule Start[tcpclient]:")

	runLogger.Println("tcpclient.eventDataPool:")
	runLogger.Println(eventDataPool)

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
