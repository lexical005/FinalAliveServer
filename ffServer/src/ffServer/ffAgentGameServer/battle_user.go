package main

import (
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"time"
)

type roleStatus byte

const (
	roleStatusLoad  roleStatus = iota // 异步加载
	roleStatusAlive                   // 存活中
	roleStatusDead                    // 死亡
	roleStatusLeave                   // 离开(逃跑, 结算(死亡或获胜)后离开)
)

type battleUser struct {
	agent  *agentUser
	status roleStatus

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

	inShootState bool // 在射击状态

	kill int32 // 击杀
}

func (agent *battleUser) String() string {
	return fmt.Sprintf("%p:%v:%v", agent, agent.uniqueid, agent.status)
}

func (agent *battleUser) UUID() uuid.UUID {
	return agent.agent.UUID()
}

// 异步加载完成
func (agent *battleUser) LoadAsyncOver() error {
	log.RunLogger.Printf("battleUser[%v].LoadAsyncOver", agent)

	battle, err := agent.Check(false)
	if err != nil {
		return err
	}

	agent.status = roleStatusAlive

	battle.TryStart()
	return nil
}

// SendProtoExtraDataNormal 发送Proto
//	返回值仅表明请求发送的协议, 是否被添加到待发送管道内, 不代表一定能发送到对端
func (agent *battleUser) SendProtoExtraDataNormal(proto *ffProto.Proto) bool {
	return agent.agent.SendProtoExtraDataNormal(proto)
}

// 检查战斗状态
func (agent *battleUser) Check(checkAlive bool) (*battleScene, error) {
	// 检查活着
	if checkAlive {
		if agent.health < 1 {
			return nil, fmt.Errorf("battleUser[%v].Check not alive", agent)
		}
	}

	// 检查战场
	battle, err := instBattleGameWorld.CheckScene(agent)
	if err != nil {
		return nil, err
	}

	return battle, err
}

// 修改持有的物品, 场景内新增扔掉的物品
func (agent *battleUser) DropBagItem(itemtemplateid, dropItemData int32, position *ffProto.StVector3) error {
	log.RunLogger.Printf("battleUser[%v].DropBagItem itemtemplateid[%v] dropItemData[%v]", agent, itemtemplateid, dropItemData)

	battle, err := agent.Check(true)
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
	log.RunLogger.Printf("battleUser[%v].DropEquip EquipIndex[%v]", agent, message.EquipIndex)

	battle, err := agent.Check(true)
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
	log.RunLogger.Printf("battleUser[%v].PickProp Uniqueid[%v]", agent, message.Itemuniqueid)

	battle, err := agent.Check(true)
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
	log.RunLogger.Printf("battleUser[%v].SwitchWeapon EquipIndex[%v=>%v]", agent, agent.itemManager.activeWeaponIndex, message.EquipIndex)

	battle, err := agent.Check(true)
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
	battle, err := agent.Check(true)
	if err != nil {
		return err
	}

	battle.RoleMove(agent.uniqueid, message)
	return nil
}

// 视野
func (agent *battleUser) RoleEyeRotate(message *ffProto.MsgBattleRoleEyeRotate) error {
	battle, err := agent.Check(true)
	if err != nil {
		return err
	}

	battle.RoleEyeRotate(agent.uniqueid, message)
	return nil
}

// 行为
func (agent *battleUser) RoleAction(message *ffProto.MsgBattleRoleAction) error {
	log.RunLogger.Printf("battleUser[%v].RoleAction Action[%v]", agent, message.Action)

	battle, err := agent.Check(true)
	if err != nil {
		return err
	}

	battle.RoleAction(agent.uniqueid, message)
	return nil
}

// 射击状态
func (agent *battleUser) RoleShootState(message *ffProto.MsgBattleRoleShootState) error {
	log.RunLogger.Printf("battleUser[%v].RoleShootState State[%v]", agent, message.State)

	battle, err := agent.Check(true)
	if err != nil {
		return err
	}

	agent.inShootState = message.State

	battle.RoleShootState(agent.uniqueid, message)
	return nil
}

// 射击
func (agent *battleUser) RoleShoot(message *ffProto.MsgBattleRoleShoot) error {
	log.RunLogger.Printf("battleUser[%v].RoleShoot", agent)

	if !agent.inShootState {
		return fmt.Errorf("battleUser[%v].RoleShoot not in shoot state", agent)
	}

	battle, err := agent.Check(true)
	if err != nil {
		return err
	}

	battle.RoleShoot(agent.uniqueid, message)
	return nil
}

// 射击结果
func (agent *battleUser) RoleShootHit(message *ffProto.MsgBattleRoleShootHit) error {
	log.RunLogger.Printf("battleUser[%v].RoleShootHit Shootid[%v] Targetuniqueid[%v]", agent, message.Shootid, message.Targetuniqueid)

	battle, err := agent.Check(true)
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

// 死亡
func (agent *battleUser) Dead(aliveCount int32, rankCount int32) {
	agent.inShootState = false

	if agent.status != roleStatusLeave {
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleSettle)
		m := p.Message().(*ffProto.MsgBattleSettle)
		m.Rank = aliveCount
		m.RankCount = rankCount
		m.Health = 0
		m.Kill = agent.kill
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}
}

// 治疗开始
func (agent *battleUser) HealStart(message *ffProto.MsgBattleRoleHeal) error {
	log.RunLogger.Printf("battleUser[%v].HealStart Itemtemplateid[%v]", agent, message.Itemtemplateid)

	battle, err := agent.Check(true)
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
	log.RunLogger.Printf("battleUser[%v].HealCancel Itemtemplateid[%v]", agent, message.Itemtemplateid)

	battle, err := agent.Check(true)
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
	log.RunLogger.Printf("battleUser[%v].HealSettle Itemtemplateid[%v]", agent, message.Itemtemplateid)

	battle, err := agent.Check(true)
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
	log.RunLogger.Printf("battleUser[%v].RunAway", agent)

	battle, err := agent.Check(false)
	if err != nil {
		return err
	}

	agent.status = roleStatusLeave

	// 逃跑
	battle.RunAway(agent)

	agent.agent.battleUser = nil
	return nil
}

// 战斗结算后离开
func (agent *battleUser) Leave() error {
	log.RunLogger.Printf("battleUser[%v].Leave", agent)

	agent.status = roleStatusLeave

	agent.agent.battleUser = nil
	return nil
}

func newBattleUser(agent *agentUser, uuidBattle uuid.UUID, member *ffProto.StBattleMember) *battleUser {
	return &battleUser{
		agent:  agent,
		status: roleStatusLoad,

		member:     member,
		uuidBattle: uuidBattle,
		uniqueid:   member.Uniqueid,

		itemManager: newItemManager(),
		healManager: newHealManager(),
		attrManager: newAttrManager(),

		health: 100,
	}
}
