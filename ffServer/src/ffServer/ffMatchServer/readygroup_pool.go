package main

import "ffCommon/pool"

type readyGroupPool struct {
	pool *pool.Pool
}

func (p *readyGroupPool) apply() *readyGroup {
	agent, _ := p.pool.Apply().(*readyGroup)
	return agent
}

func (p *readyGroupPool) back(agent *readyGroup) {
	p.pool.Back(agent)
}

func (p *readyGroupPool) String() string {
	return p.pool.String()
}

func newReadyGroupPool(initCount int) *readyGroupPool {
	funcCreator := func() interface{} {
		return newReadyGroup()
	}
	return &readyGroupPool{
		pool: pool.New("readyGroupPool", false, funcCreator, initCount, 50),
	}
}
