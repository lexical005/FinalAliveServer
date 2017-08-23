package main

import (
	"ffCommon/log/log"
	"ffCommon/net/netmanager"
	"fmt"
	"sync/atomic"
)

type matchServerClient struct {
	netManager *netmanager.Manager

	matchServer *matchServer // 连接对象
}

// Create 创建
func (client *matchServerClient) Create(netsession netmanager.INetSession) netmanager.INetSessionHandler {
	log.RunLogger.Printf("matchServerClient.Create netsession[%v]", netsession)

	// 初始化
	client.matchServer.Init(netsession)

	return client.matchServer
}

// Back 回收
func (client *matchServerClient) Back(handler netmanager.INetSessionHandler) {
	log.RunLogger.Printf("matchServerClient.Back handler[%v]", handler)

	// 回收清理
	client.matchServer.Back()
}

// Start 开始建立服务
func (client *matchServerClient) Start() error {
	log.RunLogger.Printf("matchServerClient.Start")

	manager, err := netmanager.NewClient(client, appConfig.ConnectMatchServer, appConfig.Session, &waitApplicationQuit, chApplicationQuit)
	if err != nil {
		log.FatalLogger.Println(err)
		return err
	}

	client.netManager = manager
	client.matchServer = newMatchServer()

	atomic.AddInt32(&waitApplicationQuit, 1)

	return err
}

// End 退出完成
func (client *matchServerClient) End() {
	log.RunLogger.Printf("matchServerClient.End")

	atomic.AddInt32(&waitApplicationQuit, -1)
}

// Status 当前状态描述
func (client *matchServerClient) Status() string {
	return fmt.Sprintf("matchServer[%v] netManager[%v]",
		client.matchServer, client.netManager.Status())
}
