package main

import (
	"ffCommon/log/log"
	"ffCommon/net/netmanager"
	"ffCommon/uuid"
	"fmt"
	"sync/atomic"
)

type agentUserServer struct {
	netManager *netmanager.Manager

	mapAgent  map[uuid.UUID]*agentUser // mapAgent 所有连接
	agentPool *agentUserPool           // agentPool 所有连接缓存
}

// Create 创建
func (server *agentUserServer) Create(netsession netmanager.INetSession) netmanager.INetSessionHandler {
	log.RunLogger.Printf("agentUserServer.Create netsession[%v]", netsession)

	// 申请
	agent := server.agentPool.apply()

	// 初始化
	agent.Init(netsession)

	// 记录
	server.mapAgent[agent.UUID()] = agent

	return agent
}

// Back 回收
func (server *agentUserServer) Back(handler netmanager.INetSessionHandler) {
	log.RunLogger.Printf("agentUserServer.Back handler[%v]", handler)

	agent, _ := handler.(*agentUser)

	// 清除记录
	delete(server.mapAgent, agent.UUID())

	// 回收清理
	agent.Back()

	// 缓存
	server.agentPool.back(agent)
}

// Start 开始建立服务
func (server *agentUserServer) Start() error {
	log.RunLogger.Printf("agentUserServer.Start")

	manager, err := netmanager.NewServer(server, appConfig.ServeUser, &waitApplicationQuit, chApplicationQuit)
	if err != nil {
		log.FatalLogger.Println(err)
		return err
	}

	server.netManager = manager
	server.mapAgent = make(map[uuid.UUID]*agentUser, appConfig.ServeUser.InitOnlineCount)
	server.agentPool = newAgentUserPool("agentUserServer", appConfig.ServeUser.InitOnlineCount)

	atomic.AddInt32(&waitApplicationQuit, 1)

	return err
}

// End 退出完成
func (server *agentUserServer) End() {
	log.RunLogger.Printf("agentUserServer.End")

	atomic.AddInt32(&waitApplicationQuit, 1)
}

// Status 当前服务状态描述
func (server *agentUserServer) Status() string {
	return fmt.Sprintf("mapAgent[%v] agentPool[%v] netManager[%v]",
		len(server.mapAgent), server.agentPool, server.netManager)
}
