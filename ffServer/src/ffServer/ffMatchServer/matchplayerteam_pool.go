package main

import "ffCommon/pool"

// matchTeamPool
type matchTeamPool struct {
	pool *pool.Pool
}

func (p *matchTeamPool) apply() *matchPlayer {
	agent, _ := p.pool.Apply().(*matchPlayer)
	return agent
}

func (p *matchTeamPool) back(agent *matchPlayer) {
	p.pool.Back(agent)
}

func (p *matchTeamPool) String() string {
	return p.pool.String()
}

func newMatchTeamPool(nameOwner string, initCount int) *matchTeamPool {
	funcCreator := func() interface{} {
		return newMatchTeam()
	}
	return &matchTeamPool{
		pool: pool.New(nameOwner+".matchTeam.pool", false, funcCreator, initCount, 50),
	}
}
