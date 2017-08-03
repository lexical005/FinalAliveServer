package main

import "ffCommon/pool"

type userAgentPool struct {
	pool *pool.Pool
}

func (agentPool *userAgentPool) apply() *userAgent {
	agent, _ := agentPool.pool.Apply().(*userAgent)
	return agent
}

func (agentPool *userAgentPool) back(agent *userAgent) {
	agentPool.pool.Back(agent)
}

func (agentPool *userAgentPool) String() string {
	return agentPool.pool.String()
}

func newUserAgentPool(initOnlineCount int) *userAgentPool {
	funcCreator := func() interface{} {
		return newUserAgent()
	}
	return &userAgentPool{
		pool: pool.New("userAgent.pool", false, funcCreator, initOnlineCount, 50),
	}
}
