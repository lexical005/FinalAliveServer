package main

import "ffCommon/net/base"

type agentManager struct {
}

func (am *agentManager) Create(s base.Session) base.Agent {
	return &agent{
		sess: s,
	}
}

var am = &agentManager{}
