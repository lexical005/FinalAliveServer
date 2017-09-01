package main

import (
	"ffCommon/uuid"
	"ffProto"
)

type matchPlayer struct {
	sourceServerUUID uuid.UUID // 来源服务器在本服务器上的唯一标识
	sourceServerID   int32     // 来源服务器id
	uuidPlayerKey    uuid.UUID // 与来源服务器确认同一用户(来源服务器发来的每条用户相关的协议的附加数据字段)
	uuidAccount      uuid.UUID // 用户的真实唯一id
	uuidTeam         uuid.UUID // 队伍id

	uuidBattle uuid.UUID
	uuidToken  uuid.UUID

	mode    matchMode // 匹配模式
	inMatch bool      // 是否在匹配中
}

// StartMatch 开始匹配
func (player *matchPlayer) StartMatch(mode matchMode) bool {
	player.mode, player.inMatch = mode, true

	if player.mode == matchModeSingle {
		group := instMatchMgr.GetMatchGroup(player.mode)
		return group.JoinMatch(player)
	}

	return false
}

// StopMatch 停止匹配
func (player *matchPlayer) StopMatch() bool {
	player.inMatch = false

	if player.mode == matchModeSingle {
		group := instMatchMgr.GetMatchGroup(player.mode)
		return group.LeaveMatch(player)
	}

	return false
}

// IsMatching 是否匹配中
func (player *matchPlayer) IsMatching() bool {
	return player.inMatch
}

// AllReady 匹配单元是不是已经全部准备
func (player *matchPlayer) AllReady() bool {
	return player.inMatch
}

// Count 匹配单元内有多少matchPlayer
func (player *matchPlayer) Count() int {
	return matchModeSingleUnitCount
}

// MatchSuccess 进入了准备组, 匹配完成
func (player *matchPlayer) MatchSuccess(uuidBattle uuid.UUID, uuidTokens []uuid.UUID) {
	player.uuidBattle, player.uuidToken = uuidBattle, uuidTokens[0]

	proto := ffProto.ApplyProtoForSend(ffProto.MessageType_MatchResult)
	message := proto.Message().(*ffProto.MsgMatchResult)
	message.Addr = "127.0.0.1:15201"
	message.UUIDBattle = uuidBattle.Value()
	message.UUIDToken = player.uuidToken.Value()
	instAgentGameServerMgr.SendProtoExtraDataUUID(player, proto, false)
}

func (player *matchPlayer) Init(sourceServer *agentGameServer, uuidPlayerKey, uuidAccount, uuidTeam uuid.UUID) {
	player.sourceServerID, player.sourceServerUUID = sourceServer.serverID, sourceServer.UUID()
	player.uuidPlayerKey, player.uuidAccount, player.uuidTeam = uuidPlayerKey, uuidAccount, uuidTeam
	player.uuidToken = uuid.InvalidUUID
}

func (player *matchPlayer) back() {
}

func newMatchPlayer() *matchPlayer {
	return &matchPlayer{}
}
