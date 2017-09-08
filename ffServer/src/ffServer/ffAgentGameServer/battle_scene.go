package main

import (
	"ffAutoGen/ffEnum"
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"math/rand"
	"time"
)

const (
	maxRoleCount = 50

	updateTimerInterval = 20 * time.Microsecond // 逻辑50帧
)

type userProto struct {
	uniqueid int32          // 战场内的唯一标识
	proto    *ffProto.Proto // 协议
}

type battleScene struct {
	uuidBattle uuid.UUID   // 战场标识
	uuidTokens []uuid.UUID // 进入战场用户的凭证

	idProp        int32                           // 战场场景道具编号
	props         map[int32]*ffProto.StBattleProp // 战场道具列表
	preloadBattle map[int32]int32                 // 预加载战斗资源
	preloadScene  map[int32]int32                 // 预加载场景资源

	uniqueids  []int32               // 可用战场用户编号
	agents     map[int32]*battleUser // 用户列表
	aliveCount int32                 // 剩余活着用户数

	shootid int32 // 累计射击编号

	rand *rand.Rand // 随机数

	rolePreparePosition []int32 // 角色出生点

	lastTimer time.Time

	// 用户发来的协议
	chProto chan *userProto
}

func (scene *battleScene) Init(uuidTokens []uint64) {
	scene.shootid = 1

	scene.uuidTokens = make([]uuid.UUID, 0, len(uuidTokens))
	for _, token := range uuidTokens {
		scene.uuidTokens = append(scene.uuidTokens, uuid.NewUUID(token))
	}

	seed := time.Now().Nanosecond()
	scene.rand = rand.New(rand.NewSource(int64(seed)))

	// 角色出生点
	scene.rolePreparePosition = make([]int32, len(instBattleBorn.randRolePreparePositions))
	for i := 0; i < len(scene.rolePreparePosition); i++ {
		scene.rolePreparePosition[i] = int32(i)
	}
	for i := 0; i < len(scene.rolePreparePosition); i++ {
		j := scene.rand.Intn(len(scene.rolePreparePosition) - i)
		last := len(scene.rolePreparePosition) - 1 - i
		scene.rolePreparePosition[last], scene.rolePreparePosition[j] = scene.rolePreparePosition[j], scene.rolePreparePosition[last]
	}

	// 物品
	scene.props = make(map[int32]*ffProto.StBattleProp, 1024)
	for _, config := range ffGameConfig.RandBornData.BornPrepareItem {
		err := instBattleBorn.GenItemPrepareGroup(config, scene.rand, scene)
		if err != nil {
			log.RunLogger.Println(err)
		}
	}

	battleAssetCount, sceneAssetCount := 0, 0
	for _, template := range ffGameConfig.ItemData.ItemTemplate {
		if template.ItemType == ffEnum.EItemTypeGunWeapon {
			battleAssetCount++
			sceneAssetCount++
		} else if template.ItemType == ffEnum.EItemTypeRole {
			battleAssetCount++
		} else {
			sceneAssetCount++
		}
	}
	scene.preloadBattle = make(map[int32]int32, battleAssetCount)
	scene.preloadScene = make(map[int32]int32, sceneAssetCount)

	for _, prop := range scene.props {
		template := ffGameConfig.ItemData.ItemTemplate[prop.Templateid]
		if template.ItemType == ffEnum.EItemTypeGunWeapon {
			if c, ok := scene.preloadBattle[template.AssetID]; ok {
				scene.preloadBattle[template.AssetID] = c + 1
			} else {
				scene.preloadBattle[template.AssetID] = 1
			}
			if c, ok := scene.preloadScene[template.AssetID]; ok {
				scene.preloadScene[template.AssetID] = c + 1
			} else {
				scene.preloadScene[template.AssetID] = 1
			}
		} else {
			if c, ok := scene.preloadScene[template.AssetID]; ok {
				scene.preloadScene[template.AssetID] = c + 1
			} else {
				scene.preloadScene[template.AssetID] = 1
			}
		}
	}
	scene.preloadBattle[1] = 5
	scene.preloadBattle[2] = 5

	scene.agents = make(map[int32]*battleUser, 50)

	scene.uniqueids = make([]int32, 50)
	for i := 0; i < len(scene.uniqueids); i++ {
		scene.uniqueids[i] = int32(4*i) + 1 + scene.rand.Int31n(4) // [4n, 4n+1)
	}
	for i := 0; i < len(scene.uniqueids); i++ {
		j := scene.rand.Intn(len(scene.uniqueids) - i)
		last := len(scene.uniqueids) - 1 - i
		scene.uniqueids[last], scene.uniqueids[j] = scene.uniqueids[j], scene.uniqueids[last]
	}
}

func (scene *battleScene) newMember() *ffProto.StBattleMember {

	uniqueid := scene.uniqueids[len(scene.uniqueids)-1]
	scene.uniqueids = scene.uniqueids[0 : len(scene.uniqueids)-1]

	bornPositionIndex := scene.rolePreparePosition[len(scene.rolePreparePosition)-1]
	scene.rolePreparePosition = scene.rolePreparePosition[:len(scene.rolePreparePosition)-1]
	bornPosition := instBattleBorn.randRolePreparePositions[bornPositionIndex]

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

// 用户进入
func (scene *battleScene) Enter(agent *agentUser, uuidToken uuid.UUID) error {
	for i, token := range scene.uuidTokens {
		if token == uuidToken {
			scene.uuidTokens = append(scene.uuidTokens[:i], scene.uuidTokens[i+1:]...)

			member := scene.newMember()

			battleUser := newBattleUser(agent, scene.uuidBattle, member)
			agent.battleUser = battleUser

			// 新增成员
			for _, one := range scene.agents {
				if one.status != battleStatusRunAway {
					p1 := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleMember)
					m := p1.Message().(*ffProto.MsgBattleMember)
					m.Members = append(m.Members, member)
					ffProto.SendProtoExtraDataNormal(one, p1, false)
				}
			}

			scene.agents[battleUser.uniqueid] = agent.battleUser
			scene.aliveCount++
			log.RunLogger.Printf("Enter agent[%v] uuidToken[%v] success, left uuidTokens[%v]",
				agent.UUID(), uuidToken, scene.uuidTokens)

			return nil
		}
	}
	return fmt.Errorf("Enter agent[%v] can not find uuidToken[%v] or used, valid uuidTokens[%v]",
		agent.UUID(), uuidToken, scene.uuidTokens)
}

// 用户逃跑
func (scene *battleScene) RunAway(agent *battleUser) {
	// for _, one := range scene.agents {
	// 	if one.uniqueid != agent.uniqueid && one.status != battleStatusRunAway {
	// 		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRunAway)
	// 		m := p.Message().(*ffProto.MsgBattleRunAway)
	// 		m.Roleuniqueid = agent.uniqueid
	// 		ffProto.SendProtoExtraDataNormal(one, p, false)
	// 	}
	// }
}

// 被击中
func (scene *battleScene) OnShootHit(agent *battleUser, shootid int32, targetuniqueid int32) {
	if targetuniqueid == 0 {
		return
	}

	if target, ok := scene.agents[targetuniqueid]; ok {
		if target.health > 60 {
			target.health -= 60
		} else if target.health > 0 {
			target.health = 0
		} else {
			return
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

			for _, one := range scene.agents {
				if one.status != battleStatusRunAway {
					p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleDead)
					m := p.Message().(*ffProto.MsgBattleRoleDead)
					m.Roleuniqueid = target.uniqueid
					m.Reason = 0
					m.Sourceuniqueid = agent.uniqueid
					ffProto.SendProtoExtraDataNormal(one, p, false)
				}
			}

			{
				p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleSettle)
				m := p.Message().(*ffProto.MsgBattleSettle)
				m.Rank = scene.aliveCount
				m.RankCount = int32(len(scene.agents))
				m.Health = 0
				m.Kill = target.kill
				ffProto.SendProtoExtraDataNormal(target, p, false)
			}

			scene.aliveCount--

			// 结束
			if scene.aliveCount == 1 {
				scene.Settle()
			}
		}
	}
}

// 结算
func (scene *battleScene) Settle() {
	for _, one := range scene.agents {
		if one.health > 0 && one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleSettle)
			m := p.Message().(*ffProto.MsgBattleSettle)
			m.Rank = scene.aliveCount
			m.RankCount = int32(len(scene.agents))
			m.Health = one.health
			m.Kill = one.kill
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 场景新增物品
func (scene *battleScene) newProp(templateid int32, ItemData int32, position *ffProto.StVector3) *ffProto.StBattleProp {
	prop := &ffProto.StBattleProp{
		Position:   position,
		Uniqueid:   scene.idProp,
		Templateid: templateid,
		ItemData:   ItemData,
	}
	scene.props[scene.idProp] = prop
	scene.idProp++
	return prop
}

// 获取场景物品
func (scene *battleScene) Prop(uniqueid int32) (*ffProto.StBattleProp, error) {
	prop, ok := scene.props[uniqueid]
	if !ok {
		return nil, fmt.Errorf("Prop item uniqueid[%v] not exist", uniqueid)
	}
	return prop, nil
}

// 场景移除物品, 通知所有用户
func (scene *battleScene) RemoveSceneProp(uniqueid int32) {
	// 删除场景物品
	delete(scene.props, uniqueid)

	// 广播场景物品移除
	for _, one := range scene.agents {
		if one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRemoveProp)
			m := p.Message().(*ffProto.MsgBattleRemoveProp)
			m.Uniqueid = uniqueid
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 场景内新增物品, 通知所有用户
func (scene *battleScene) AddSceneProp(itemtemplateid, itemdata int32, position *ffProto.StVector3) {
	// 场景添加物品
	prop := scene.newProp(itemtemplateid, itemdata, position)

	for _, one := range scene.agents {
		if one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleAddProp)
			m := p.Message().(*ffProto.MsgBattleAddProp)
			m.Prop = prop
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 移动
func (scene *battleScene) RoleMove(uniqueid int32, message *ffProto.MsgBattleRoleMove) {
	// 广播用户的装备状态改变
	for _, one := range scene.agents {
		if one.uniqueid != uniqueid && one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleMove)
			m := p.Message().(*ffProto.MsgBattleRoleMove)
			m.Roleuniqueid = message.Roleuniqueid
			m.Move = message.Move
			m.SpeedDocument = message.SpeedDocument
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 视野
func (scene *battleScene) RoleEyeRotate(uniqueid int32, message *ffProto.MsgBattleRoleEyeRotate) {
	// 广播用户的装备状态改变
	for _, one := range scene.agents {
		if one.uniqueid != uniqueid && one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleEyeRotate)
			m := p.Message().(*ffProto.MsgBattleRoleEyeRotate)
			m.Roleuniqueid = message.Roleuniqueid
			m.EyeRotate = message.EyeRotate
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 行为
func (scene *battleScene) RoleAction(uniqueid int32, message *ffProto.MsgBattleRoleAction) {
	// 广播用户的装备状态改变
	for _, one := range scene.agents {
		if one.uniqueid != uniqueid && one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleAction)
			m := p.Message().(*ffProto.MsgBattleRoleAction)
			m.Roleuniqueid = message.Roleuniqueid
			m.Position = message.Position
			m.Action = message.Action
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 射击状态
func (scene *battleScene) RoleShootState(uniqueid int32, message *ffProto.MsgBattleRoleShootState) {
	// 广播用户的装备状态改变
	for _, one := range scene.agents {
		if one.uniqueid != uniqueid && one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleShootState)
			m := p.Message().(*ffProto.MsgBattleRoleShootState)
			m.Roleuniqueid = uniqueid
			m.State = message.State
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 射击
func (scene *battleScene) RoleShoot(uniqueid int32, message *ffProto.MsgBattleRoleShoot) {
	message.Shootid = scene.shootid
	scene.shootid++

	// 广播用户的装备状态改变
	for _, one := range scene.agents {
		if one.uniqueid != uniqueid && one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleShoot)
			m := p.Message().(*ffProto.MsgBattleRoleShoot)
			m.Roleuniqueid = message.Roleuniqueid
			m.Shootid = message.Shootid
			m.Position = message.Position
			m.Fireposition = message.Fireposition
			m.EyeField = message.EyeField
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 射击结果
func (scene *battleScene) RoleShootHit(uniqueid int32, message *ffProto.MsgBattleRoleShootHit) {
	// 广播用户的装备状态改变
	for _, one := range scene.agents {
		if one.uniqueid != uniqueid && one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleShootHit)
			m := p.Message().(*ffProto.MsgBattleRoleShootHit)
			m.Roleuniqueid = message.Roleuniqueid
			m.Shootid = message.Shootid
			m.Targetuniqueid = message.Targetuniqueid
			m.Endtag = message.Endtag
			m.Endposition = message.Endposition
			m.Endnormal = message.Endnormal
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 装备状态改变
func (scene *battleScene) EquipStateChanged(uniqueid int32, equipState *ffProto.StBattleEquipState) {
	// 广播用户的装备状态改变
	for _, one := range scene.agents {
		if one.uniqueid != uniqueid && one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleEquipState)
			m := p.Message().(*ffProto.MsgBattleEquipState)
			m.Roleuniqueid = uniqueid
			m.EquipState = equipState
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 治疗状态改变
func (scene *battleScene) HealStateChanged(uniqueid int32, itemtemplateid int32, state int32) {
	for _, one := range scene.agents {
		if one.uniqueid != uniqueid && one.status != battleStatusRunAway {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleHeal)
			m := p.Message().(*ffProto.MsgBattleRoleHeal)
			m.Roleuniqueid = uniqueid
			m.Itemtemplateid = itemtemplateid
			m.State = state
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}

	// // 开始
	// if state == 1 {
	// 	if agent.healitemtemplateid != 0 {
	// 		return fmt.Errorf("Heal failed, healitemtemplateid[%v] state[%v] conflict with healing healitemtemplateid[%v]",
	// 			healitemtemplateid, state, agent.healitemtemplateid)
	// 	}

	// 	ItemData, ok := agent.items[healitemtemplateid]
	// 	if !ok || ItemData == 0 {
	// 		return fmt.Errorf("Heal failed, healitemtemplateid[%v] state[%v] not own item",
	// 			healitemtemplateid, state)
	// 	}

	// 	template := ffGameConfig.ItemData.ItemTemplate[healitemtemplateid]
	// 	if template.ItemType != ffEnum.EItemTypeConsumable {
	// 		return fmt.Errorf("Heal failed, healitemtemplateid[%v] state[%v] not a consumable item",
	// 			healitemtemplateid, state)
	// 	}

	// 	agent.healitemtemplateid, agent.healTime = healitemtemplateid, time.Now()
	// 	return nil
	// }

	// // 取消
	// if state == 0 {
	// 	agent.healitemtemplateid = 0
	// 	return nil
	// }

	// // 结算
	// if state == 2 {
	// 	if agent.healitemtemplateid == 0 {
	// 		return fmt.Errorf("Heal failed, healitemtemplateid[%v] state[%v] not in healing",
	// 			healitemtemplateid, state)
	// 	}
	// 	agent.healitemtemplateid = 0
	// 	agent.health += 20

	// 	// 血量同步
	// 	{
	// 		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleHealth)
	// 		m := p.Message().(*ffProto.MsgBattleRoleHealth)
	// 		m.Roleuniqueid = agent.uniqueid
	// 		m.Health = agent.health
	// 		ffProto.SendProtoExtraDataNormal(agent, p, false)
	// 	}
	// }
}

func (scene *battleScene) Start() {
	go util.SafeGo(scene.mainLoop, scene.mainLoopEnd)
}

func (scene *battleScene) mainLoop(params ...interface{}) {
	log.RunLogger.Println("battleScene.mainLoop")

	scene.lastTimer = time.Now()

	{
		waitTime := updateTimerInterval
		for {
			<-time.After(waitTime)

			// 时间驱动
			scene.lastTimer = scene.lastTimer.Add(updateTimerInterval)
			waitTime = updateTimerInterval - time.Now().Sub(scene.lastTimer)
			scene.update(updateTimerInterval)
		}
	}
}

func (scene *battleScene) mainLoopEnd(isPanic bool) {
	log.RunLogger.Println("battleScene.mainLoopEnd", isPanic)
}

func (scene *battleScene) update(passTime time.Duration) {
	c := len(scene.chProto)
	for i := 0; i < c; i++ {
		// proto := <-scene.chProto

		// proto.SetCacheDispatched()
	}
}

// // 战斗协议
// func (scene *battleScene) onBattleUserProto(userProto*userProto) (bool, bool) {
// 	if ffProto.MessageType_BattleStartSync == proto.ProtoID() {
// 		return onBattleProtoStartSync(agent, proto), true
// 	}

// 	if callback, ok := mapBattleProtoCallback[proto.ProtoID()]; ok {
// 		return callback(agent.battleUser, proto), true
// 	}

// 	return false, false
// }
