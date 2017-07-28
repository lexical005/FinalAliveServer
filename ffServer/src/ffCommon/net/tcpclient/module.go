package tcpclient

import (
	"ffCommon/net/base"
	"ffCommon/pool"
)

const (
	// initTotalClientNetEventDataCount 初始默认创建多少 clientNetEventData(供本进程所有tcpClient使用)
	initTotalClientNetEventDataCount = 32

	// initClientNetEventDataChanCount 一个 tcpClient.chNetEventDataInner 的缓存有多大
	initClientNetEventDataChanCount = 8
)

// 客户端数量
var clientCount int

// eventDataPool clientNetEventData Pool
var eventDataPool *clientNetEventDataPool

// Init 初始tcpclient模块
func Init() (err error) {
	funcCreateClientNetEventData := func() interface{} {
		return newClientNetEventData()
	}

	eventDataPool = &clientNetEventDataPool{
		pool: pool.New("tcpclient.eventDataPool", true, funcCreateClientNetEventData, initTotalClientNetEventDataCount, 50),
	}

	return
}

// NewClient create new base.Client
//	addr: 监听地址
func NewClient(addr string) (c base.Client, err error) {
	clientCount++
	return newClient(addr, clientCount)
}
