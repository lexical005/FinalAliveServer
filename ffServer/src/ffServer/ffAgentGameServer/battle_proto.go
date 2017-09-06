package main

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
	"strconv"
	"strings"
)

// 协议回调函数
//	返回值表明接收到的Proto是否进入了发送逻辑(如果未正确设置返回值, 将导致泄露或者异常)
var mapBattleProtoCallback = map[ffProto.MessageType]func(agent *battleUser, proto *ffProto.Proto) bool{
	ffProto.MessageType_BattlePickProp:      onBattleProtoPickProp,
	ffProto.MessageType_BattleDropBagProp:   onBattleProtoDropBagProp,
	ffProto.MessageType_BattleDropEquipProp: onBattleProtoDropEquipProp,
	ffProto.MessageType_BattleSwitchWeapon:  onBattleProtoSwitchWeapon,

	ffProto.MessageType_BattleRunAway: onBattleProtoRunAway,

	ffProto.MessageType_BattleRoleAction: onBattleProtoRoleAction,

	ffProto.MessageType_BattleRoleShootState: onBattleProtoRoleShootState,
	ffProto.MessageType_BattleRoleShoot:      onBattleProtoRoleShoot,
	ffProto.MessageType_BattleRoleShootHit:   onBattleProtoRoleShootHit,
	ffProto.MessageType_BattleRoleMove:       onBattleProtoRoleMove,
	ffProto.MessageType_BattleRoleEyeRotate:  onBattleProtoRoleEyeRotate,

	ffProto.MessageType_BattleRoleHeal: onBattleProtoRoleHeal,

	ffProto.MessageType_BattleCheat: onBattleProtoCheat,
}

// 战斗协议
func onBattleProto(agent *agentUser, proto *ffProto.Proto) (bool, bool) {
	if ffProto.MessageType_BattleStartSync == proto.ProtoID() {
		return onBattleProtoStartSync(agent, proto), true
	}

	if callback, ok := mapBattleProtoCallback[proto.ProtoID()]; ok {
		return callback(agent.battleUser, proto), true
	}

	return false, false
}

// 开始同步
func onBattleProtoStartSync(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleStartSync)

	uuidBattle := uuid.NewUUID(message.UUIDBattle)
	uuidToken := uuid.NewUUID(message.UUIDToken)
	battle, ok := mapBattle[uuidBattle]

	if !ok || !battle.Enter(agent, uuidBattle, uuidToken) {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	// 反馈
	message.Uniqueid = agent.battleUser.uniqueid
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

// 逃跑
func onBattleProtoRunAway(agent *battleUser, proto *ffProto.Proto) (result bool) {
	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 捡取场景物品
func onBattleProtoPickProp(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattlePickProp)

	// 捡取
	if err := agent.PickProp(message); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}
	message.Itemuniqueid, message.Position = 0, nil

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 丢弃非装备
func onBattleProtoDropBagProp(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleDropBagProp)

	// 战场不存在或者扔失败
	if err := agent.DropBagItem(message.Itemtemplateid, message.ItemData, message.Position); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}
	message.Position = nil

	// 扔成功
	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 丢弃装备
func onBattleProtoDropEquipProp(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleDropEquipProp)

	message.EquipState = &ffProto.StBattleEquipState{}
	if err := agent.DropEquip(message); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}
	message.Position = nil

	// 扔成功
	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 行为
func onBattleProtoRoleAction(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleAction)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	for _, one := range battle.agents {
		if one.uniqueid != agent.uniqueid {
			p := ffProto.ApplyProtoForSend(proto.ProtoID())
			m := p.Message().(*ffProto.MsgBattleRoleAction)
			m.Roleuniqueid = message.Roleuniqueid
			m.Position = message.Position
			m.Action = message.Action
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 射击开火状态
func onBattleProtoRoleShootState(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleShootState)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	for _, one := range battle.agents {
		if one.uniqueid != agent.uniqueid {
			p := ffProto.ApplyProtoForSend(proto.ProtoID())
			m := p.Message().(*ffProto.MsgBattleRoleShootState)
			m.Roleuniqueid = message.Roleuniqueid
			m.State = message.State
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 射击开火
func onBattleProtoRoleShoot(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleShoot)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	Shootid := battle.shootid
	for _, one := range battle.agents {
		if one.uniqueid != agent.uniqueid {
			p := ffProto.ApplyProtoForSend(proto.ProtoID())
			m := p.Message().(*ffProto.MsgBattleRoleShoot)
			m.Roleuniqueid = message.Roleuniqueid
			m.Shootid = Shootid
			m.Position = message.Position
			m.Fireposition = message.Fireposition
			m.EyeField = message.EyeField
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 射击击中
func onBattleProtoRoleShootHit(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleShootHit)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	for _, agent := range battle.agents {
		if agent.uniqueid != message.Roleuniqueid {
			p := ffProto.ApplyProtoForSend(proto.ProtoID())
			m := p.Message().(*ffProto.MsgBattleRoleShootHit)
			m.Roleuniqueid = message.Roleuniqueid
			m.Shootid = message.Shootid
			m.Targetuniqueid = message.Targetuniqueid
			m.Endtag = message.Endtag
			m.Endposition = message.Endposition
			m.Endnormal = message.Endnormal
			ffProto.SendProtoExtraDataNormal(agent, p, false)
		}
	}

	if message.Targetuniqueid != 0 {
		battle.ShootHit(agent, message.Shootid, message.Targetuniqueid)
	}

	// 客户端已进行了预表现, 不需要给客户端返回
	return
}

// 移动
func onBattleProtoRoleMove(agent *battleUser, proto *ffProto.Proto) (result bool) {
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
func onBattleProtoRoleEyeRotate(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleEyeRotate)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		return false
	}

	for _, agent := range battle.agents {
		if agent.uniqueid == message.Roleuniqueid {
			continue
		}

		p := ffProto.ApplyProtoForSend(proto.ProtoID())
		m := p.Message().(*ffProto.MsgBattleRoleEyeRotate)
		m.Roleuniqueid = message.Roleuniqueid
		m.EyeRotate = message.EyeRotate
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

// 治疗
func onBattleProtoRoleHeal(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleHeal)

	// 开始治疗
	if message.State == 1 {
		if err := agent.HealStart(message); err != nil {
			message.Result = ffError.ErrUnknown.Code()
			log.RunLogger.Println(err)
		}
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 切换武器
func onBattleProtoSwitchWeapon(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleSwitchWeapon)

	if err := agent.SwitchWeapon(message); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 战斗作弊指令
func onBattleProtoCheat(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleCheat)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	cmd := strings.ToLower(message.Cmd)
	if cmd == "settle" {
		battle.Settle()
	} else if strings.HasPrefix("item ", " ") {
		tmp := strings.Split(cmd, " ")
		itemtemplate, _ := strconv.Atoi(tmp[1])
		itemdata, _ := strconv.Atoi(tmp[2])
		x, _ := strconv.Atoi(tmp[3])
		y, _ := strconv.Atoi(tmp[4])
		z, _ := strconv.Atoi(tmp[5])
		pos := &ffProto.StVector3{
			X: int64(x),
			Y: int64(y),
			Z: int64(z),
		}
		if battle, ok := mapBattle[agent.uuidBattle]; ok {
			battle.AddSceneProp(int32(itemtemplate), int32(itemdata), pos)
		}
	}

	return
}
