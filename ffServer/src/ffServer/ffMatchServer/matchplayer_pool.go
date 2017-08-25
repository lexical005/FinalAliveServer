package main

import "ffCommon/pool"

type matchPlayerPool struct {
	pool *pool.Pool
}

func (p *matchPlayerPool) apply() *matchPlayer {
	agent, _ := p.pool.Apply().(*matchPlayer)
	return agent
}

func (p *matchPlayerPool) back(agent *matchPlayer) {
	p.pool.Back(agent)
}

func (p *matchPlayerPool) String() string {
	return p.pool.String()
}

func newMatchPlayerPool(nameOwner string, initCount int) *matchPlayerPool {
	funcCreator := func() interface{} {
		return newMatchPlayer()
	}
	return &matchPlayerPool{
		pool: pool.New(nameOwner+".matchPlayer.pool", false, funcCreator, initCount, 50),
	}
}
