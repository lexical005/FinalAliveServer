package main

import (
	"ffAutoGen/ffEnum"
	"ffAutoGen/ffError"
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"strconv"
	"strings"
)

// 协议回调函数
//	返回值表明接收到的Proto是否进入了发送逻辑(如果未正确设置返回值, 将导致泄露或者异常)
var mapBattleProtoCallback = map[ffProto.MessageType]func(agent *battleUser, proto *ffProto.Proto) bool{
	ffProto.MessageType_LoadAsyncOver: onBattleProtoLoadAsyncOver,

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

// 进入战场同步
func onBattleProtoStartSync(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleStartSync)

	uuidBattle := uuid.NewUUID(message.UUIDBattle)
	uuidToken := uuid.NewUUID(message.UUIDToken)
	battle, ok := instBattleGameWorld.mapScene[uuidBattle]

	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Printf("onBattleProtoStartSync uuidBattle[%v] not exist", uuidBattle)
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	if err := battle.Enter(agent, uuidToken); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	// 反馈
	message.Uniqueid = agent.battleUser.uniqueid
	message.PreloadBattle = battle.preloadBattle
	message.PreloadScene = battle.preloadScene
	result = ffProto.SendProtoExtraDataNormal(agent, proto, true)

	// 成员列表
	{
		p1 := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleMember)
		m := p1.Message().(*ffProto.MsgBattleMember)
		for _, agent := range battle.agents {
			m.Members = append(m.Members, agent.member)
		}
		ffProto.SendProtoExtraDataNormal(agent, p1, false)
	}

	// 战场道具
	{
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleProp)
		m := p.Message().(*ffProto.MsgBattleProp)
		for _, prop := range battle.props {
			m.Props = append(m.Props, prop)
		}
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

// 场景异步加载完成
func onBattleProtoLoadAsyncOver(agent *battleUser, proto *ffProto.Proto) (result bool) {
	if err := agent.LoadAsyncOver(); err != nil {
		log.RunLogger.Println(err)
	}
	return false
}

// 逃跑
func onBattleProtoRunAway(agent *battleUser, proto *ffProto.Proto) (result bool) {
	if err := agent.RunAway(); err != nil {
		message, _ := proto.Message().(*ffProto.MsgBattleRunAway)
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}

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

	if err := agent.RoleAction(message); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 射击开火状态
func onBattleProtoRoleShootState(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleShootState)

	if err := agent.RoleShootState(message); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 射击开火
func onBattleProtoRoleShoot(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleShoot)

	if err := agent.RoleShoot(message); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 射击击中
func onBattleProtoRoleShootHit(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleShootHit)

	if err := agent.RoleShootHit(message); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}

	// 客户端已进行了预表现, 不需要给客户端返回
	return
}

// 移动
func onBattleProtoRoleMove(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleMove)

	if err := agent.RoleMove(message); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 视野
func onBattleProtoRoleEyeRotate(agent *battleUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleEyeRotate)

	if err := agent.RoleEyeRotate(message); err != nil {
		log.RunLogger.Println(err)
	}

	// 客户端已进行了预表现, 不需要给客户端返回
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

	battle, ok := instBattleGameWorld.mapScene[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	cmd := strings.ToLower(message.Cmd)
	fmt.Println("onBattleProtoCheat", cmd)
	if cmd == "settle" {
		battle.Settle()
	} else if strings.HasPrefix(cmd, "item ") {
		tmp := strings.Split(cmd, " ")
		itemtemplateid, _ := strconv.Atoi(tmp[1])
		itemdata, _ := strconv.Atoi(tmp[2])
		x, _ := strconv.Atoi(tmp[3])
		y, _ := strconv.Atoi(tmp[4])
		z, _ := strconv.Atoi(tmp[5])
		pos := &ffProto.StVector3{
			X: int64(x),
			Y: int64(y),
			Z: int64(z),
		}

		template, ok := ffGameConfig.ItemData.ItemTemplate[int32(itemtemplateid)]
		if !ok {
			fmt.Println("cheat cmd item, invalid itemtemplateid", itemtemplateid)
			message.Result = ffError.ErrUnknown.Code()
			return ffProto.SendProtoExtraDataNormal(agent, proto, true)
		}

		if template.ItemType == ffEnum.EItemTypeGunWeapon ||
			template.ItemType == ffEnum.EItemTypeAttachment ||
			template.ItemType == ffEnum.EItemTypeConsumable ||
			template.ItemType == ffEnum.EItemTypeThrowable ||
			template.ItemType == ffEnum.EItemTypeMelleeWeapon {
			battle.AddSceneProp(int32(itemtemplateid), 1, pos)

		} else if template.ItemType == ffEnum.EItemTypeAmmunition {
			if itemdata < 1 || itemdata > 300 {
				fmt.Println("cheat cmd item, invalid ammunition itemdata", itemdata)
				message.Result = ffError.ErrUnknown.Code()
				return ffProto.SendProtoExtraDataNormal(agent, proto, true)
			}
			battle.AddSceneProp(int32(itemtemplateid), int32(itemdata), pos)

		} else if template.ItemType == ffEnum.EItemTypeArmor {
			armor := ffGameConfig.ItemData.Armor[int32(itemtemplateid)]
			if itemdata < 1 || int32(itemdata) > armor.Attrs[ffEnum.EAttrDurable] {
				fmt.Println("cheat cmd item, invalid armor itemdata", itemdata)
				message.Result = ffError.ErrUnknown.Code()
				return ffProto.SendProtoExtraDataNormal(agent, proto, true)
			}
			battle.AddSceneProp(int32(itemtemplateid), int32(itemdata), pos)
		}
	}

	return
}
