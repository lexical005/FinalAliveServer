package main

import (
	"ffCommon/log/log"
	"ffCommon/uuid"
	"sync"
)

type matchPlayerManager struct {
	mutexPlayer sync.RWMutex               // 用户操作锁
	players     map[uuid.UUID]*matchPlayer // 当前有效用户. key: uuidPlayerKey
}

// 记录player
func (mgr *matchPlayerManager) addPlayer(player *matchPlayer) {
	mgr.mutexPlayer.Lock()
	defer mgr.mutexPlayer.Unlock()
	mgr.players[player.uuidPlayerKey] = player
}

// 移除player
func (mgr *matchPlayerManager) delPlayer(player *matchPlayer) {
	mgr.mutexPlayer.Lock()
	defer mgr.mutexPlayer.Unlock()

	delete(mgr.players, player.uuidPlayerKey)

	player.back()
}

// 获取player
func (mgr *matchPlayerManager) GetPlayer(uuidPlayerKey uuid.UUID) *matchPlayer {
	mgr.mutexPlayer.RLock()
	defer mgr.mutexPlayer.RUnlock()

	if player, ok := mgr.players[uuidPlayerKey]; ok {
		return player
	}
	return nil
}

// AddPlayer 添加匹配用户
func (mgr *matchPlayerManager) AddPlayer(agent *agentGameServer, uuidPlayerKey, uuidAccount, uuidTeam uuid.UUID) {
	log.RunLogger.Printf("matchPlayerManager.AddPlayer agentGameServer[%v] uuidPlayerKey[%v]", agent.UUID(), uuidPlayerKey)

	player := mgr.GetPlayer(uuidPlayerKey)
	if player == nil {
		player = newMatchPlayer(agent, uuidPlayerKey, uuidAccount, uuidTeam)
		mgr.addPlayer(player)
	}
}

// RemovePlayer 移除匹配用户
func (mgr *matchPlayerManager) RemovePlayer(uuidPlayerKey uuid.UUID) {
	log.RunLogger.Printf("matchPlayerManager.RemovePlayer uuidPlayerKey[%v]", uuidPlayerKey)

	if player, ok := mgr.players[uuidPlayerKey]; ok {
		mgr.delPlayer(player)
	}
}

func (mgr *matchPlayerManager) Start() error {
	mgr.players = make(map[uuid.UUID]*matchPlayer, appConfig.Match.InitMatchCount/2)

	return nil
}
