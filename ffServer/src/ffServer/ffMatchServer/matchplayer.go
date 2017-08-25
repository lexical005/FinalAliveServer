package main

type matchPlayer struct {
	sourceServer  *agentGameServer // 来源服务器
	uuidPlayerKey uint64           // 与来源服务器确认同一用户(来源服务器发来的每条用户相关的协议的附加数据字段)
	uuidAccount   uint64           // 用户的真实唯一id
	uuidTeam      uint64           // 队伍id

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
	return 1
}

// MatchSuccess 进入了准备组, 匹配完成
func (player *matchPlayer) MatchSuccess() {

}

func (player *matchPlayer) Init(sourceServer *agentGameServer, uuidPlayerKey, uuidAccount, uuidTeam uint64) {
	player.sourceServer, player.uuidPlayerKey, player.uuidAccount, player.uuidTeam = sourceServer, uuidPlayerKey, uuidAccount, uuidTeam
}

func (player *matchPlayer) back() {
	player.sourceServer = nil
}

func newMatchPlayer() *matchPlayer {
	return &matchPlayer{}
}
