package main

import "ffCommon/pool"

type agentUserPool struct {
	pool *pool.Pool
}

func (p *agentUserPool) apply() *agentUser {
	agent, _ := p.pool.Apply().(*agentUser)
	return agent
}

func (p *agentUserPool) back(agent *agentUser) {
	p.pool.Back(agent)
}

func (p *agentUserPool) String() string {
	return p.pool.String()
}

func newAgentUserPool(nameOwner string, initOnlineCount int) *agentUserPool {
	funcCreator := func() interface{} {
		return newAgentUser()
	}
	return &agentUserPool{
		pool: pool.New(nameOwner+".agentUser.pool", false, funcCreator, initOnlineCount, 50),
	}
}
