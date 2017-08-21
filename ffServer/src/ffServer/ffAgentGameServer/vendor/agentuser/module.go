package agentuser

import (
	"ffCommon/net/base"
	"ffCommon/util"
	"ffCommon/uuid"
	"fmt"
)

// NewServer 根据配置返回一个Server管理器
func NewServer(config *base.ServeConfig,
	countApplicationQuit *int32,
	chApplicationQuit chan struct{}) (mgr *Manager, err error) {

	net, err := newNetServer(config)
	if err != nil {
		return
	}

	name := fmt.Sprintf("ServerManager[%v]", net.UUID())

	mgr = &Manager{
		name: name,

		net: net,

		chAgentClosed: make(chan *agentSession, 2),
		mapAgent:      make(map[uuid.UUID]*agentSession, config.InitOnlineCount),
		agentPool:     newAgentSessionPool(name, config.InitOnlineCount),

		countApplicationQuit: countApplicationQuit,
		chApplicationQuit:    chApplicationQuit,
	}

	go util.SafeGo(mgr.mainLoop, mgr.mainLoopEnd)

	return
}

// NewClient 根据配置返回一个Client管理器
func NewClient(config *base.ConnectConfig,
	countApplicationQuit *int32,
	chApplicationQuit chan struct{}) (mgr *Manager, err error) {

	net, err := newNetClient(config)
	if err != nil {
		return
	}

	name := fmt.Sprintf("ClientManager[%v]", net.UUID())

	mgr = &Manager{
		name: name,

		net: net,

		chAgentClosed: make(chan *agentSession, 2),
		mapAgent:      make(map[uuid.UUID]*agentSession, 1),
		agentPool:     newAgentSessionPool(name, 1),

		countApplicationQuit: countApplicationQuit,
		chApplicationQuit:    chApplicationQuit,
	}

	go util.SafeGo(mgr.mainLoop, mgr.mainLoopEnd)

	return
}
