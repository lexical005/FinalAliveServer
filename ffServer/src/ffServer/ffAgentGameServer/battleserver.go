package main

import (
	"ffAutoGen/ffEnum"
	"ffAutoGen/ffError"
	"ffCommon/uuid"
	"ffProto"
	"math/rand"
)

var mapBattle = make(map[uuid.UUID]*battle, 1)

type battle struct {
	uuidBattle uuid.UUID
	uuidTokens []uuid.UUID

	idProp int32
	Props  map[int32]*ffProto.StBattleProp

	PreloadBattle map[int32]int32
	PreloadScene  map[int32]int32

	agents map[int32]*agentUser

	Members []*ffProto.StBattleMember

	uniqueids []int32

	bornPosition []*ffProto.StVector3

	totalCount int32
	aliveCount int32

	// 累计射击编号
	shootid int32
}

func (b *battle) Init(uuidTokens []uint64) {
	b.shootid = 1

	for _, token := range uuidTokens {
		b.uuidTokens = append(b.uuidTokens, uuid.NewUUID(token))
	}

	b.PreloadBattle = map[int32]int32{
		1:     1,
		2:     1,
		10301: 1,
		10401: 1,
		1101:  2,
		1201:  2,
		1301:  2,
		1401:  2,
		1501:  2,
		1601:  2,
		1701:  2,
	}
	b.PreloadScene = map[int32]int32{
		10301: 1,
		10401: 1,
	}

	b.agents = make(map[int32]*agentUser, 50)

	b.Members = make([]*ffProto.StBattleMember, 0, 50)

	b.Props = make(map[int32]*ffProto.StBattleProp, 50)
	b.NewProp(10301, 1, &ffProto.StVector3{
		X: -7414,
		Y: 155718,
		Z: 2800251,
	})
	b.NewProp(10401, 1, &ffProto.StVector3{
		X: -4202,
		Y: 155718,
		Z: 2800251,
	})

	b.bornPosition = make([]*ffProto.StVector3, 0, 50)
	b.bornPosition = append(b.bornPosition, &ffProto.StVector3{
		X: -708,
		Y: 155718,
		Z: 2803553,
	})
	b.bornPosition = append(b.bornPosition, &ffProto.StVector3{
		X: -708,
		Y: 155718,
		Z: 2812334,
	})

	b.uniqueids = make([]int32, 50)
	for i := 0; i < len(b.uniqueids); i++ {
		b.uniqueids[i] = int32(4*i) + 1 + rand.Int31n(4) // [4n, 4n+1)
	}
	for i := 0; i < len(b.uniqueids); i++ {
		j := rand.Intn(len(b.uniqueids) - i)
		last := len(b.uniqueids) - 1 - i
		b.uniqueids[last], b.uniqueids[j] = b.uniqueids[j], b.uniqueids[last]
	}

	// if proto.ProtoID() == ffProto.MessageType_MatchResult {
	// 	message := proto.Message().(*ffProto.MsgMatchResult)
	// 	if message.Result == 0 {
	// 		go util.SafeGo(func(params ...interface{}) {
	// 			<-time.After(time.Second * 30)

	// 			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleSettle)
	// 			m := p.Message().(*ffProto.MsgBattleSettle)
	// 			m.Rank = 1
	// 			m.RankCount = 1
	// 			m.Kill = 0
	// 			m.RankScore = 1
	// 			m.KillScore = 0
	// 			m.Health = 100
	// 			ffProto.SendProtoExtraDataNormal(agent, p, false)
	// 		}, nil)
	// 	}
	// }
}

func (b *battle) newMember() *ffProto.StBattleMember {

	uniqueid := b.uniqueids[len(b.uniqueids)-1]
	b.uniqueids = b.uniqueids[0 : len(b.uniqueids)-1]

	bornPosition := b.bornPosition[len(b.bornPosition)-1]
	b.bornPosition = b.bornPosition[0 : len(b.bornPosition)-1]

	member := &ffProto.StBattleMember{
		Position: bornPosition,
		Sight:    &ffProto.StVector3{},
		Uniqueid: uniqueid,
		Datas: map[int32]int32{
			int32(ffEnum.EActorAttrActor): 1,
		},
	}
	return member
}

func (b *battle) AddToken(token uint64) {
	b.uuidTokens = append(b.uuidTokens, uuid.NewUUID(token))
}

func (b *battle) RemoveToken(token uint64) {
	t := uuid.NewUUID(token)
	for i, id := range b.uuidTokens {
		if id == t {
			b.uuidTokens = append(b.uuidTokens[:i], b.uuidTokens[i+1:]...)
			break
		}
	}
}

func (b *battle) Enter(agent *agentUser, UUIDBattle, UUIDToken uuid.UUID) bool {
	for i, token := range b.uuidTokens {
		if token == UUIDToken {
			b.uuidTokens = append(b.uuidTokens[:i], b.uuidTokens[i+1:]...)

			member := b.newMember()

			agent.uuidBattle, agent.uniqueid, agent.health = UUIDBattle, member.Uniqueid, 100

			b.Members = append(b.Members, member)

			// 新增成员
			for _, agent := range b.agents {
				p1 := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleMember)
				m := p1.Message().(*ffProto.MsgBattleMember)
				m.Members = append(m.Members, member)
				ffProto.SendProtoExtraDataNormal(agent, p1, false)
			}

			b.agents[agent.uniqueid] = agent
			b.totalCount++
			b.aliveCount++
			return true
		}
	}
	return false
}

func (b *battle) Shoot(agent *agentUser, shootid int32) {

}

func (b *battle) ShootHit(agent *agentUser, shootid int32, Targetuniqueid int32) {
	if Targetuniqueid == 0 {
		return
	}

	if target, ok := b.agents[Targetuniqueid]; ok {
		if target.health > 60 {
			target.health -= 60
		} else {
			target.health = 0
		}

		// 血量同步
		{
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleHealth)
			m := p.Message().(*ffProto.MsgBattleRoleHealth)
			m.Roleuniqueid = target.uniqueid
			m.Health = target.health
			ffProto.SendProtoExtraDataNormal(target, p, false)
		}

		// 死亡
		if target.health == 0 {
			agent.kill++

			for _, agent := range b.agents {
				p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleDead)
				m := p.Message().(*ffProto.MsgBattleRoleDead)
				m.Roleuniqueid = target.uniqueid
				m.Reason = 0
				m.Sourceuniqueid = agent.uniqueid
				ffProto.SendProtoExtraDataNormal(agent, p, false)
			}

			{
				p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleSettle)
				m := p.Message().(*ffProto.MsgBattleSettle)
				m.Rank = b.aliveCount
				m.RankCount = b.totalCount
				m.Health = 0
				m.Kill = target.kill
				ffProto.SendProtoExtraDataNormal(target, p, false)
			}

			b.aliveCount--

			// 结束
			if b.aliveCount == 0 {
				for _, agent := range b.agents {
					if agent.health > 0 {
						p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleSettle)
						m := p.Message().(*ffProto.MsgBattleSettle)
						m.Rank = b.aliveCount
						m.RankCount = b.totalCount
						m.Health = agent.health
						m.Kill = agent.kill
						ffProto.SendProtoExtraDataNormal(agent, p, false)
					}
				}
			}
		}
	}
}

func (b *battle) NewProp(templateid int32, number int32, position *ffProto.StVector3) *ffProto.StBattleProp {
	prop := &ffProto.StBattleProp{
		Position:   position,
		Uniqueid:   b.idProp,
		Templateid: templateid,
		Number:     number,
	}
	b.Props[b.idProp] = prop
	b.idProp++
	return prop
}

func (b *battle) Prop(uniqueid int32) (*ffProto.StBattleProp, bool) {
	prop, ok := b.Props[uniqueid]
	return prop, ok
}

func (b *battle) PickProp(uniqueid int32) {
	delete(b.Props, uniqueid)
}

// func (b *battle) SendProtoToAllExcept(proto *ffProto.Proto) {
// 	message := proto.Message()
// 	for _, agent := range b.agents {
// 		p := ffProto.ApplyProtoForSend(proto.ProtoID())
// 		m := p.Message()
// 		m.
// 		ffProto.SendProtoExtraDataNormal(agent, p, false)
// 	}
// }

func onBattleProtoStartSync(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleStartSync)

	UUIDBattle := uuid.NewUUID(message.UUIDBattle)
	battle, ok := mapBattle[UUIDBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	UUIDToken := uuid.NewUUID(message.UUIDToken)
	if !battle.Enter(agent, UUIDBattle, UUIDToken) {
		message.Result = ffError.ErrUnknown.Code()
	}

	// 反馈
	message.Uniqueid = agent.uniqueid
	message.PreloadBattle = battle.PreloadBattle
	message.PreloadScene = battle.PreloadScene
	result = ffProto.SendProtoExtraDataNormal(agent, proto, true)

	// 成员列表
	{
		p1 := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleMember)
		m := p1.Message().(*ffProto.MsgBattleMember)
		m.Members = battle.Members
		ffProto.SendProtoExtraDataNormal(agent, p1, false)
	}

	// 战场道具
	{
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleProp)
		m := p.Message().(*ffProto.MsgBattleProp)
		for _, prop := range battle.Props {
			m.Props = append(m.Props, prop)
		}

		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

func onBattleProtoRunAway(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRunAway)
	message.Result = ffError.ErrNone.Code()
	ffProto.SendProtoExtraDataNormal(agent, proto, true)

	return
}

func onBattleProtoPickProp(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattlePickProp)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	prop, ok := battle.Prop(message.Itemuniqueid)
	if ok {
		battle.PickProp(prop.Uniqueid)
	} else {
		message.Result = ffError.ErrUnknown.Code()
	}
	message.Roleuniqueid, message.Itemuniqueid = 0, 0
	message.Itemtemplateid, message.Itemnumber = prop.Templateid, prop.Number

	if ok {
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRemoveProp)
		m := p.Message().(*ffProto.MsgBattleRemoveProp)
		m.Uniqueid = prop.Uniqueid
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

func onBattleProtoDropProp(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleDropProp)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	// 扔的物品的判定

	message.Roleuniqueid = 0
	result = ffProto.SendProtoExtraDataNormal(agent, proto, true)

	{
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleAddProp)
		m := p.Message().(*ffProto.MsgBattleAddProp)
		m.Prop = battle.NewProp(message.Itemtemplateid, message.Itemnumber, message.Position)
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

// 行为
func onBattleProtoRoleAction(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleAction)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	for _, agent := range battle.agents {
		p := ffProto.ApplyProtoForSend(proto.ProtoID())
		m := p.Message().(*ffProto.MsgBattleRoleAction)
		m.Roleuniqueid = message.Roleuniqueid
		m.Position = message.Position
		m.Action = message.Action
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

// 射击开火状态
func onBattleProtoRoleShootState(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleShootState)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	for _, agent := range battle.agents {
		p := ffProto.ApplyProtoForSend(proto.ProtoID())
		m := p.Message().(*ffProto.MsgBattleRoleShootState)
		m.Roleuniqueid = message.Roleuniqueid
		m.State = message.State
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

// 射击开火
func onBattleProtoRoleShoot(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleShoot)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	Shootid := battle.shootid
	for _, agent := range battle.agents {
		p := ffProto.ApplyProtoForSend(proto.ProtoID())
		m := p.Message().(*ffProto.MsgBattleRoleShoot)
		m.Roleuniqueid = message.Roleuniqueid
		m.Shootid = Shootid
		m.Position = message.Position
		m.Fireposition = message.Fireposition
		m.EyeField = message.EyeField
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

// 射击击中
func onBattleProtoRoleShootHit(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleShootHit)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	for _, agent := range battle.agents {
		if agent.uniqueid == message.Roleuniqueid {
			continue
		}

		p := ffProto.ApplyProtoForSend(proto.ProtoID())
		m := p.Message().(*ffProto.MsgBattleRoleShootHit)
		m.Roleuniqueid = message.Roleuniqueid
		m.Shootid = message.Shootid
		m.Targetuniqueid = message.Targetuniqueid
		m.Endposition = message.Endposition
		m.Endnormal = message.Endnormal
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

// 移动
func onBattleProtoRoleMove(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleMove)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	for _, agent := range battle.agents {
		p := ffProto.ApplyProtoForSend(proto.ProtoID())
		m := p.Message().(*ffProto.MsgBattleRoleMove)
		m.Roleuniqueid = message.Roleuniqueid
		m.Move = message.Move
		m.SpeedDocument = message.SpeedDocument
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

// 视野
func onBattleProtoRoleEyeField(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleEyeField)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		return false
	}

	for _, agent := range battle.agents {
		p := ffProto.ApplyProtoForSend(proto.ProtoID())
		m := p.Message().(*ffProto.MsgBattleRoleEyeField)
		m.Roleuniqueid = message.Roleuniqueid
		m.EyeField = message.EyeField
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

func newBattle(uuidBattle uuid.UUID) *battle {
	battle := &battle{
		uuidBattle: uuidBattle,
		uuidTokens: make([]uuid.UUID, 0, 50),
	}
	mapBattle[uuidBattle] = battle
	return battle
}
