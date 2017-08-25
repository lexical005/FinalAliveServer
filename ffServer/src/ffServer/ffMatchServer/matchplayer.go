package main

import (
	"ffCommon/uuid"
	"ffProto"
)

type matchPlayer struct {
	sourceServer  *agentGameServer // 来源服务器
	uuidPlayerKey uuid.UUID        // 与来源服务器确认同一用户(来源服务器发来的每条用户相关的协议的附加数据字段)
	uuidAccount   uuid.UUID        // 用户的真实唯一id
	uuidTeam      uuid.UUID        // 队伍id

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
func (player *matchPlayer) MatchSuccess() {
	proto := ffProto.ApplyProtoForSend(ffProto.MessageType_MatchResult)
	message := proto.Message().(*ffProto.MsgMatchResult)
	message.Addr = "127.0.0.1:15201"
	message.Token = player.uuidAccount.Value()
	ffProto.SendProtoExtraDataUUID(player.sourceServer, player.uuidPlayerKey, proto, false)
}

func (player *matchPlayer) Init(sourceServer *agentGameServer, uuidPlayerKey, uuidAccount, uuidTeam uuid.UUID) {
	player.sourceServer, player.uuidPlayerKey, player.uuidAccount, player.uuidTeam = sourceServer, uuidPlayerKey, uuidAccount, uuidTeam
}

func (player *matchPlayer) back() {
	player.sourceServer = nil
}

func newMatchPlayer() *matchPlayer {
	return &matchPlayer{}
}
