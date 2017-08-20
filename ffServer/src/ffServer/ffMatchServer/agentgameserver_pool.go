package main

import "ffCommon/pool"

type agentGameServerPool struct {
	pool *pool.Pool
}

func (agentPool *agentGameServerPool) apply() *agentGameServer {
	agent, _ := agentPool.pool.Apply().(*agentGameServer)
	return agent
}

func (agentPool *agentGameServerPool) back(agent *agentGameServer) {
	agentPool.pool.Back(agent)
}

func (agentPool *agentGameServerPool) String() string {
	return agentPool.pool.String()
}

func newAgentGameServerPool(initOnlineCount int) *agentGameServerPool {
	funcCreator := func() interface{} {
		return newAgentGameServer()
	}
	return &agentGameServerPool{
		pool: pool.New("agentGameServer.pool", false, funcCreator, initOnlineCount, 50),
	}
}
