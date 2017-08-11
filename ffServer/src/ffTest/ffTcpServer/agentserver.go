package main

import (
	"ffCommon/uuid"
	"ffProto"
	"sync"

	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpserver"
)

// agentManager 创建Server, 并维护连接对应的am实例
type agentManager struct {
	extraDataType ffProto.ExtraDataType // extraDataType 附加数据类型

	server base.Server

	mutexClient sync.RWMutex
	mapClients  map[uuid.UUID]*agentClient
}

// CreateNetAgent 创建NetAgent对象, 不会失败!
func (am *agentManager) CreateNetAgent(uuidSession uuid.UUID) base.NetAgent {
	agent := newAgentClient(uuidSession)

	am.mutexClient.RLock()
	defer am.mutexClient.RUnlock()
	am.mapClients[uuidSession] = agent

	return agent
}

// CloseAllNetAgent 通知所有NetAgent, 让其关闭其对应的Session(close(NetAgent.SendProtoChannel()))
func (am *agentManager) CloseAllNetAgent() {
	am.mutexClient.RLock()
	defer am.mutexClient.RUnlock()
	for _, agent := range am.mapClients {
		agent.CloseSession()
	}
}

// OnSessionNetEventData 被Client或Server调用, 接收到网络事件, 同步
func (am *agentManager) OnSessionNetEventData(data base.SessionNetEventData) {
	log.RunLogger.Printf("agentManager.onNetDataEvent data[%v]: %v", data, am)

	defer data.Back()

	am.mutexClient.RLock()
	defer am.mutexClient.RUnlock()
	if agent, ok := am.mapClients[data.Session().UUID()]; ok {
		switch data.NetEventType() {
		case base.NetEventOn:
			agent.OnNetConnect()
		case base.NetEventOff:
			agent.OnNetDisConnect(data.ManualClose())
		case base.NetEventProto:
			agent.OnNetProto(data.Proto())
		}
	} else {
		log.FatalLogger.Printf("agentManager.OnSessionNetEventData session[%v] not exist agent!", data.Session())
	}
}

func (am *agentManager) Start() error {
	return am.server.Start(am, am.extraDataType)
}

func newServerAgent(addr string, extraDataType ffProto.ExtraDataType) (am *agentManager, err error) {
	server, err := tcpserver.NewServer(addr)
	if err != nil {
		return nil, err
	}

	am = &agentManager{
		extraDataType: extraDataType,

		server: server,
	}
	return
}

var gServer *agentManager
