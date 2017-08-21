package agentuser

import "ffCommon/pool"

type agentSessionPool struct {
	pool *pool.Pool
}

func (p *agentSessionPool) apply() *agentSession {
	agent, _ := p.pool.Apply().(*agentSession)
	return agent
}

func (p *agentSessionPool) back(agent *agentSession) {
	p.pool.Back(agent)
}

func (p *agentSessionPool) String() string {
	return p.pool.String()
}

func newAgentSessionPool(nameOwner string, initOnlineCount int) *agentSessionPool {
	funcCreator := func() interface{} {
		return newAgentSession()
	}
	return &agentSessionPool{
		pool: pool.New(nameOwner+".agentSession.pool", false, funcCreator, initOnlineCount, 50),
	}
}
