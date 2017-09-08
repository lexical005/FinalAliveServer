package main

import (
	"ffCommon/uuid"
	"ffProto"
	"time"
)

type battleStatus byte

const (
	battleStatusAlive battleStatus = iota
	battleStatusDead
	battleStatusRunAway
)

type battleUser struct {
	agent  *agentUser
	status battleStatus

	member     *ffProto.StBattleMember
	uuidBattle uuid.UUID // 战场
	uniqueid   int32     // 战场内的唯一标识

	itemManager *itemManager // 物品管理
	healManager *healManager // 治疗管理
	attrManager *attrManager // 属性管理

	health      int32 // 血量
	bodyDefence int32 // 身体防御
	headDefence int32 // 头部防御

	healitemtemplateid int32     // 正在使用的治疗物品模板id
	healTime           time.Time // 开始heal的时间

	kill int32 // 击杀
}

func (agent *battleUser) UUID() uuid.UUID {
	return agent.agent.UUID()
}

// SendProtoExtraDataNormal 发送Proto
//	返回值仅表明请求发送的协议, 是否被添加到待发送管道内, 不代表一定能发送到对端
func (agent *battleUser) SendProtoExtraDataNormal(proto *ffProto.Proto) bool {
	return agent.agent.SendProtoExtraDataNormal(proto)
}

// 修改持有的物品, 场景内新增扔掉的物品
func (agent *battleUser) DropBagItem(itemtemplateid, dropItemData int32, position *ffProto.StVector3) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	if err := agent.itemManager.DropBagItem(itemtemplateid, dropItemData); err != nil {
		return err
	}

	battle.AddSceneProp(itemtemplateid, dropItemData, position)
	return nil
}

// 修改装备状态, 修改持有的物品, 场景内新增扔掉的物品
func (agent *battleUser) DropEquip(message *ffProto.MsgBattleDropEquipProp) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	// 丢弃装备
	if err = agent.itemManager.DropEquip(agent, message.EquipIndex, message.Position, message.EquipState); err != nil {
		return err
	}

	// 装备状态改变了
	battle.EquipStateChanged(agent.uniqueid, message.EquipState)

	return nil
}

// 捡取场景里的物品
func (agent *battleUser) PickProp(message *ffProto.MsgBattlePickProp) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	prop, err := battle.Prop(message.Itemuniqueid)
	if err != nil {
		return err
	}

	// 捡物品
	err = agent.itemManager.PickProp(agent, prop, message)
	if err != nil {
		return err
	}

	// 场景移除物品
	battle.RemoveSceneProp(prop.Uniqueid)

	// 装备状态改变了
	if message.EquipState != nil {
		battle.EquipStateChanged(agent.uniqueid, message.EquipState)
	}

	return nil
}

// 切换武器
func (agent *battleUser) SwitchWeapon(message *ffProto.MsgBattleSwitchWeapon) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	// 切换武器
	err = agent.itemManager.SwitchWeapon(message)
	if err != nil {
		return err
	}

	// 装备状态改变了
	battle.EquipStateChanged(agent.uniqueid, message.EquipState)

	return nil
}

// 移动
func (agent *battleUser) RoleMove(message *ffProto.MsgBattleRoleMove) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	battle.RoleMove(agent.uniqueid, message)
	return nil
}

// 视野
func (agent *battleUser) RoleEyeRotate(message *ffProto.MsgBattleRoleEyeRotate) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	battle.RoleEyeRotate(agent.uniqueid, message)
	return nil
}

// 行为
func (agent *battleUser) RoleAction(message *ffProto.MsgBattleRoleAction) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	battle.RoleAction(agent.uniqueid, message)
	return nil
}

// 射击状态
func (agent *battleUser) RoleShootState(message *ffProto.MsgBattleRoleShootState) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	battle.RoleShootState(agent.uniqueid, message)
	return nil
}

// 射击
func (agent *battleUser) RoleShoot(message *ffProto.MsgBattleRoleShoot) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	battle.RoleShoot(agent.uniqueid, message)
	return nil
}

// 射击结果
func (agent *battleUser) RoleShootHit(message *ffProto.MsgBattleRoleShootHit) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	// 广播击中表现
	battle.RoleShootHit(agent.uniqueid, message)

	// 处理击中逻辑
	if message.Targetuniqueid != 0 {
		battle.OnShootHit(agent, message.Shootid, message.Targetuniqueid)
	}

	return nil
}

// 治疗开始
func (agent *battleUser) HealStart(message *ffProto.MsgBattleRoleHeal) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	// 治疗开始
	if err = agent.healManager.HealStart(agent, message.Itemtemplateid); err != nil {
		return err
	}

	// 治疗状态改变
	battle.HealStateChanged(agent.uniqueid, message.Itemtemplateid, message.State)
	return nil
}

// 治疗中断
func (agent *battleUser) HealCancel(message *ffProto.MsgBattleRoleHeal) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	// 治疗中断
	if err = agent.healManager.HealCancel(agent); err != nil {
		return err
	}

	// 治疗状态改变
	battle.HealStateChanged(agent.uniqueid, 0, message.State)
	return nil
}

// 治疗结算
func (agent *battleUser) HealSettle(message *ffProto.MsgBattleRoleHeal) error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	// 治疗结算
	if err = agent.healManager.HealSettle(agent); err != nil {
		return err
	}

	// 治疗状态改变
	battle.HealStateChanged(agent.uniqueid, 0, message.State)
	return nil
}

// 逃跑(主动逃跑, 断线)
func (agent *battleUser) RunAway() error {
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return err
	}

	agent.status = battleStatusRunAway

	// 逃跑
	battle.RunAway(agent)
	return nil
}

func newBattleUser(agent *agentUser, uuidBattle uuid.UUID, member *ffProto.StBattleMember) *battleUser {
	return &battleUser{
		agent:  agent,
		status: battleStatusAlive,

		member:     member,
		uuidBattle: uuidBattle,
		uniqueid:   member.Uniqueid,

		itemManager: newItemManager(),
		healManager: newHealManager(),
		attrManager: newAttrManager(),

		health: 100,
	}
}
