package main

import (
	"ffAutoGen/ffError"
	"ffCommon/net/base"
	"ffCommon/net/tcpserver"
	"ffCommon/uuid"
	"ffProto"

	"sync"
)

// 客户端侧管理
type tunnelClientAgentManager struct {
	clientType string
	server     base.Server

	muAgents sync.RWMutex
	agents   map[uuid.UUID]*tunnelClientAgent
}

func (tcam *tunnelClientAgentManager) Create(s base.Session) base.Agent {
	agent := clientAgentPool.apply()
	agent.sess = s

	// 等待底层通知连接建立完成
	agent.muSessionOn.Lock()

	return agent
}

func (tcam *tunnelClientAgentManager) sendProto(id uuid.UUID, p *ffProto.Proto) bool {
	tcam.muAgents.RLock()
	defer tcam.muAgents.RUnlock()

	if agent, ok := tcam.agents[id]; ok {
		return agent.sendProto(p)
	}
	return false
}

func (tcam *tunnelClientAgentManager) kick(id uuid.UUID, reason int32) {
	tcam.muAgents.RLock()
	defer tcam.muAgents.RUnlock()

	if agent, ok := tcam.agents[id]; ok {
		agent.kick(true, ffError.ErrByCode(reason))
	}
}

func (tcam *tunnelClientAgentManager) onConnect(agent *tunnelClientAgent) {
	tcam.muAgents.Lock()

	// 解锁, 可发送协议了
	agent.muSessionOn.Unlock()

	tcam.agents[agent.uuid()] = agent

	tcam.muAgents.Unlock()
}

func (tcam *tunnelClientAgentManager) onDisconnect(id uuid.UUID) {
	tcam.muAgents.Lock()

	delete(tcam.agents, id)

	tcam.muAgents.Unlock()
}

// 根据配置创建Server
func (tcam *tunnelClientAgentManager) init() error {
	server, err := tcpserver.NewServer(appConfig.ClientListen.ListenAddr)
	if err != nil {
		return err
	}

	tcam.server = server
	tcam.clientType = appConfig.ClientListen.ClientType
	tcam.agents = make(map[uuid.UUID]*tunnelClientAgent, appConfig.Session.OnlineCount)

	return nil
}

// Server启动监听
func (tcam *tunnelClientAgentManager) start() error {
	return tcam.server.Start(tcam, ffProto.ExtraDataTypeNormal)
}
