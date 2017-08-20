package agentuser

import "ffCommon/pool"

type agentUserPool struct {
	pool *pool.Pool
}

func (agentPool *agentUserPool) apply() *agentUser {
	agent, _ := agentPool.pool.Apply().(*agentUser)
	return agent
}

func (agentPool *agentUserPool) back(agent *agentUser) {
	agentPool.pool.Back(agent)
}

func (agentPool *agentUserPool) String() string {
	return agentPool.pool.String()
}

func newAgentUserPool(initOnlineCount int) *agentUserPool {
	funcCreator := func() interface{} {
		return newAgentUser()
	}
	return &agentUserPool{
		pool: pool.New("userAgent.pool", false, funcCreator, initOnlineCount, 50),
	}
}
