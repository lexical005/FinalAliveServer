package main

import "time"
import "ffProto"
import "ffCommon/uuid"

// readyGroup 准备组, 正在补齐人数, 筹备战场, 等待战场开启
type readyGroup struct {
	uuidBattle    uuid.UUID    // 战场唯一标识
	units         []iMatchUnit // 涉及的所有匹配单元
	lackCount     int          // 距离满员还差多少玩家
	stopEnterTime time.Time    // 战场最晚开启时间

	sourceServerUUID uuid.UUID // 来源服务器连接标识
	sourceServerID   int32     // 来源服务器
}

// Init 初始化准备组
func (group *readyGroup) Init(fullCount int, stopEnterTime time.Time, units []iMatchUnit, count int) {
	// todo: 测试代码
	for k, v := range instAgentGameServerMgr.mapAgent {
		group.sourceServerUUID = k
		group.sourceServerID = v.serverID
	}

	group.lackCount, group.stopEnterTime = fullCount, stopEnterTime
	group.uuidBattle = instMatchMgr.uuidBattleGenerator.Gen()

	proto := ffProto.ApplyProtoForSend(ffProto.MessageType_ServerNewBattle)
	message := proto.Message().(*ffProto.MsgServerNewBattle)
	message.UUIDBattle = group.uuidBattle.Value()

	group.units = nil
	group.units = append(group.units, units...)

	uuidTokens := make([]uuid.UUID, 0, count)
	for _, unit := range units {
		for i := 0; i < unit.Count(); i++ {
			token := instMatchMgr.uuidBattleTokenGenerator.Gen()
			uuidTokens = append(uuidTokens, token)
			message.UserTokens = append(message.UserTokens, token.Value())
		}
	}
	group.lackCount -= count

	// todo: 通知创建战场
	instAgentGameServerMgr.SendServerProto(group.sourceServerUUID, group.sourceServerID, proto, false)

	for _, unit := range units {
		c := unit.Count()
		unit.MatchSuccess(group.uuidBattle, uuidTokens[:c])
		uuidTokens = uuidTokens[c:]
	}
}

// Enter 一个匹配单元进入此准备组
//	返回值, 是否满员
func (group *readyGroup) Enter(unit iMatchUnit) bool {
	count := unit.Count()

	group.units = append(group.units, unit)
	group.lackCount -= count

	proto := ffProto.ApplyProtoForSend(ffProto.MessageType_ServerBattleUserEnter)
	message := proto.Message().(*ffProto.MsgServerBattleUserEnter)
	message.UUIDBattle = group.uuidBattle.Value()

	uuidTokens := make([]uuid.UUID, 0, count)
	for i := 0; i < count; i++ {
		token := instMatchMgr.uuidBattleTokenGenerator.Gen()
		uuidTokens = append(uuidTokens, token)
		message.UserToken = token.Value()
	}

	// todo: 通知战场内新增玩家
	instAgentGameServerMgr.SendServerProto(group.sourceServerUUID, group.sourceServerID, proto, false)

	unit.MatchSuccess(group.uuidBattle, uuidTokens)

	return group.lackCount == 0
}

// IsFull 是否满员
func (group *readyGroup) IsFull() bool {
	return group.lackCount == 0
}

func newReadyGroup() *readyGroup {
	return &readyGroup{}
}
