package main

import (
	"ffCommon/log/log"
	"ffCommon/net/netmanager"
	"ffCommon/uuid"
	"fmt"
	"sync/atomic"
)

type agentGameServerManager struct {
	netManager *netmanager.Manager

	mapAgent  map[uuid.UUID]*agentGameServer // mapAgent 所有连接
	agentPool *agentGameServerPool           // agentPool 所有连接缓存
}

// Create 创建
func (mgr *agentGameServerManager) Create(netsession netmanager.INetSession) netmanager.INetSessionHandler {
	log.RunLogger.Printf("agentUserServer.Create netsession[%v]", netsession)

	// 申请
	agent := mgr.agentPool.apply()

	// 初始化
	agent.Init(netsession)

	// 记录
	mgr.mapAgent[agent.UUID()] = agent

	return agent
}

// Back 回收
func (mgr *agentGameServerManager) Back(handler netmanager.INetSessionHandler) {
	log.RunLogger.Printf("agentUserServer.Back handler[%v]", handler)

	agent, _ := handler.(*agentGameServer)

	// 清除记录
	delete(mgr.mapAgent, agent.UUID())

	// 回收清理
	agent.Back()

	// 缓存
	mgr.agentPool.back(agent)
}

// Start 开始建立服务
func (mgr *agentGameServerManager) Start() error {
	log.RunLogger.Printf("agentUserServer.Start")

	manager, err := netmanager.NewServer(mgr, appConfig.ServeAgentGameServer, &waitApplicationQuit, chApplicationQuit)
	if err != nil {
		log.FatalLogger.Println(err)
		return err
	}

	mgr.netManager = manager
	mgr.mapAgent = make(map[uuid.UUID]*agentGameServer, appConfig.ServeAgentGameServer.InitOnlineCount)
	mgr.agentPool = newAgentGameServerPool("agentUserServer", appConfig.ServeAgentGameServer.InitOnlineCount)

	atomic.AddInt32(&waitApplicationQuit, 1)

	return err
}

// End 退出完成
func (mgr *agentGameServerManager) End() {
	log.RunLogger.Printf("agentUserServer.End")

	atomic.AddInt32(&waitApplicationQuit, -1)
}

// Status 当前服务状态描述
func (mgr *agentGameServerManager) Status() string {
	return fmt.Sprintf("mapAgent[%v] agentPool[%v] netManager[%v]",
		len(mgr.mapAgent), mgr.agentPool, mgr.netManager.Status())
}
