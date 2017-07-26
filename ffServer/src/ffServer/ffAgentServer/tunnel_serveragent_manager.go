package main

import (
	"ffCommon/net/base"
	"ffCommon/net/tcpserver"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"

	"sync"
	"time"
)

// 服务端侧管理
type tunnelServerAgentManager struct {
	serverType string
	server     base.Server

	muAgents sync.RWMutex
	agents   map[uuid.UUID]*tunnelServerAgent
}

func (tsam *tunnelServerAgentManager) Create(s base.Session) base.Agent {
	agent := &tunnelServerAgent{
		sess: s,

		serverID: 0,
	}

	// 等待底层通知连接建立完成
	agent.muSessionOn.Lock()

	return agent
}

func (tsam *tunnelServerAgentManager) onConnect(agent *tunnelServerAgent) {
	tsam.muAgents.Lock()

	// 解锁, 可发送协议了
	agent.muSessionOn.Unlock()

	tsam.agents[agent.uuid()] = agent

	tsam.muAgents.Unlock()
}

func (tsam *tunnelServerAgentManager) onDisconnect(id uuid.UUID) {
	tsam.muAgents.Lock()

	delete(tsam.agents, id)

	tsam.muAgents.Unlock()
}

func (tsam *tunnelServerAgentManager) getServerAgent(serverID int) *tunnelServerAgent {
	tsam.muAgents.RLock()

	for _, agent := range tsam.agents {
		if agent.serverID == serverID {

			tsam.muAgents.RUnlock()

			return agent
		}
	}

	tsam.muAgents.RUnlock()

	return nil
}

// 保持现有连接不断开
func (tsam *tunnelServerAgentManager) keepAliveLoop(params ...interface{}) {
	interval := time.Duration(appConfig.Session.ReadDeadTime*40/100) * time.Second

	for {
		select {
		case <-time.After(interval):
			tsam.keepAliveOneLoop()
		}
	}
}
func (tsam *tunnelServerAgentManager) keepAliveOneLoop() {
	tsam.muAgents.RLock()
	defer tsam.muAgents.RUnlock()

	// 向每一个现有连接发送一条MsgKeepAlive协议, 以避免协议被物理断开
	for _, agent := range tsam.agents {
		agent.sendProto(ffProto.ApplyProtoForSend(ffProto.MessageType_MT_MsgServerKeepAlive), 0)
	}
}

// 根据配置创建Server
func (tsam *tunnelServerAgentManager) init() error {
	server, err := tcpserver.NewServer(appConfig.ServerListen.ListenAddr)
	if err != nil {
		return err
	}

	tsam.server = server
	tsam.serverType = appConfig.ServerListen.ServerType
	tsam.agents = make(map[uuid.UUID]*tunnelServerAgent, 1)

	return nil
}

// Server启动监听
func (tsam *tunnelServerAgentManager) start() error {
	// 新起一个协程, 用于定时向已建立的连接发送KeepAlive协议, 以避免连接被自然断开
	go util.SafeGo(tsam.keepAliveLoop)

	return tsam.server.Start(tsam, ffProto.ExtraDataTypeUUID)
}
