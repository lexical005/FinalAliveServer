package main

import (
	"ffCommon/pool"

	"fmt"
)

type tunnelClientAgentPool struct {
	pool *pool.Pool
}

func (tcap *tunnelClientAgentPool) apply() *tunnelClientAgent {
	agent, _ := tcap.pool.Apply().(*tunnelClientAgent)
	return agent
}

func (tcap *tunnelClientAgentPool) back(agent *tunnelClientAgent) {
	tcap.pool.Back(agent)
}

func (tcap *tunnelClientAgentPool) String() string {
	return tcap.pool.String()
}

func (tcap *tunnelClientAgentPool) init() error {
	if appConfig.Session.OnlineCount < 1 {
		return fmt.Errorf("tunnelClientAgentPool.Init: invalid appConfig.Session.OnlineCount[%v]", appConfig.Session.OnlineCount)
	}

	tcap.pool = pool.New("tunnelClientAgent.pool", false, newClientAgent, appConfig.Session.OnlineCount, 50)
	return nil
}
