package main

import (
	"ffCommon/log/log"
	"ffCommon/uuid"
)

type matchCenter struct {
	chMatchGroupSingle chan uuid.UUID
}

// MatchSingle 单人匹配
func (center *matchCenter) MatchSingle(agent *agentGameServer) {
	center.chMatchGroupSingle <- agent.uuid
}

// MatchCancel 取消匹配
func (center *matchCenter) MatchCancel(agent *agentGameServer) {
}

func (center *matchCenter) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("matchCenter.mainLoop: %v", center)
}
func (center *matchCenter) mainLoopEnd() {
	log.RunLogger.Printf("matchCenter.mainLoopEnd: %v", center)
}
