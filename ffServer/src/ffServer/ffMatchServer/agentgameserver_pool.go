package main

import "ffCommon/pool"

type agentGameServerPool struct {
	pool *pool.Pool
}

func (p *agentGameServerPool) apply() *agentGameServer {
	agent, _ := p.pool.Apply().(*agentGameServer)
	return agent
}

func (p *agentGameServerPool) back(agent *agentGameServer) {
	p.pool.Back(agent)
}

func (p *agentGameServerPool) String() string {
	return p.pool.String()
}

func newAgentGameServerPool(nameOwner string, initOnlineCount int) *agentGameServerPool {
	funcCreator := func() interface{} {
		return newAgentGameServer()
	}
	return &agentGameServerPool{
		pool: pool.New(nameOwner+".agentGameServer.pool", false, funcCreator, initOnlineCount, 50),
	}
}
