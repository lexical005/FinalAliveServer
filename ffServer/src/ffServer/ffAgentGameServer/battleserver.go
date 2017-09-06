package main

import (
	"ffAutoGen/ffEnum"
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"math/rand"
	"time"
)

var mapBattle = make(map[uuid.UUID]*battle, 1)

const (
	maxRoleCount = 50
)

type battle struct {
	uuidBattle uuid.UUID
	uuidTokens []uuid.UUID

	idProp int32
	Props  map[int32]*ffProto.StBattleProp

	PreloadBattle map[int32]int32
	PreloadScene  map[int32]int32

	agents map[int32]*battleUser

	Members []*ffProto.StBattleMember

	uniqueids []int32

	totalCount int32
	aliveCount int32

	// 累计射击编号
	shootid int32

	// 随机数
	Rand *rand.Rand

	// 角色出生点
	rolePreparePosition []int32

	BattleProps map[int32]*ffProto.StBattleProp
}

func (b *battle) Init(uuidTokens []uint64) {
	b.shootid = 1

	for _, token := range uuidTokens {
		b.uuidTokens = append(b.uuidTokens, uuid.NewUUID(token))
	}

	seed := time.Now().Nanosecond()
	b.Rand = rand.New(rand.NewSource(int64(seed)))

	// 角色出生点
	b.rolePreparePosition = make([]int32, len(instBattleBorn.randRolePreparePositions))
	for i := 0; i < len(b.rolePreparePosition); i++ {
		b.rolePreparePosition[i] = int32(i)
	}
	for i := 0; i < len(b.rolePreparePosition); i++ {
		j := b.Rand.Intn(len(b.rolePreparePosition) - i)
		last := len(b.rolePreparePosition) - 1 - i
		b.rolePreparePosition[last], b.rolePreparePosition[j] = b.rolePreparePosition[j], b.rolePreparePosition[last]
	}

	// 物品
	b.Props = make(map[int32]*ffProto.StBattleProp, 1024)
	for _, config := range ffGameConfig.RandBornData.BornPrepareItem {
		err := instBattleBorn.GenItemPrepareGroup(config, b.Rand, b)
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
	b.PreloadBattle = make(map[int32]int32, battleAssetCount)
	b.PreloadScene = make(map[int32]int32, sceneAssetCount)

	for _, prop := range b.Props {
		template := ffGameConfig.ItemData.ItemTemplate[prop.Templateid]
		if template.ItemType == ffEnum.EItemTypeGunWeapon {
			if c, ok := b.PreloadBattle[template.AssetID]; ok {
				b.PreloadBattle[template.AssetID] = c + 1
			} else {
				b.PreloadBattle[template.AssetID] = 1
			}
			if c, ok := b.PreloadScene[template.AssetID]; ok {
				b.PreloadScene[template.AssetID] = c + 1
			} else {
				b.PreloadScene[template.AssetID] = 1
			}
		} else {
			if c, ok := b.PreloadScene[template.AssetID]; ok {
				b.PreloadScene[template.AssetID] = c + 1
			} else {
				b.PreloadScene[template.AssetID] = 1
			}
		}
	}
	b.PreloadBattle[1] = 5
	b.PreloadBattle[2] = 5

	b.agents = make(map[int32]*battleUser, 50)

	b.Members = make([]*ffProto.StBattleMember, 0, maxRoleCount)

	b.uniqueids = make([]int32, 50)
	for i := 0; i < len(b.uniqueids); i++ {
		b.uniqueids[i] = int32(4*i) + 1 + b.Rand.Int31n(4) // [4n, 4n+1)
	}
	for i := 0; i < len(b.uniqueids); i++ {
		j := b.Rand.Intn(len(b.uniqueids) - i)
		last := len(b.uniqueids) - 1 - i
		b.uniqueids[last], b.uniqueids[j] = b.uniqueids[j], b.uniqueids[last]
	}
}

func (b *battle) newMember() *ffProto.StBattleMember {

	uniqueid := b.uniqueids[len(b.uniqueids)-1]
	b.uniqueids = b.uniqueids[0 : len(b.uniqueids)-1]

	bornPositionIndex := b.rolePreparePosition[len(b.rolePreparePosition)-1]
	b.rolePreparePosition = b.rolePreparePosition[:len(b.rolePreparePosition)-1]
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

func (b *battle) Enter(agent *agentUser, uuidBattle, uuidToken uuid.UUID) bool {
	for i, token := range b.uuidTokens {
		if token == uuidToken {
			b.uuidTokens = append(b.uuidTokens[:i], b.uuidTokens[i+1:]...)

			member := b.newMember()

			battleUser := newBattleUser(agent, uuidBattle, member.Uniqueid)
			agent.battleUser = battleUser

			b.Members = append(b.Members, member)

			// 新增成员
			for _, agent := range b.agents {
				p1 := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleMember)
				m := p1.Message().(*ffProto.MsgBattleMember)
				m.Members = append(m.Members, member)
				ffProto.SendProtoExtraDataNormal(agent, p1, false)
			}

			b.agents[battleUser.uniqueid] = agent.battleUser
			b.totalCount++
			b.aliveCount++
			return true
		}
	}
	return false
}

func (b *battle) Shoot(agent *battleUser, shootid int32) {

}

func (b *battle) ShootHit(agent *battleUser, shootid int32, Targetuniqueid int32) {
	if Targetuniqueid == 0 {
		return
	}

	if target, ok := b.agents[Targetuniqueid]; ok {
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
			if b.aliveCount == 1 {
				b.Settle()
			}
		}
	}
}

// 结算
func (b *battle) Settle() {
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

// 场景新增物品
func (b *battle) newProp(templateid int32, ItemData int32, position *ffProto.StVector3) *ffProto.StBattleProp {
	prop := &ffProto.StBattleProp{
		Position:   position,
		Uniqueid:   b.idProp,
		Templateid: templateid,
		ItemData:   ItemData,
	}
	b.Props[b.idProp] = prop
	b.idProp++
	return prop
}

// 获取场景物品
func (b *battle) Prop(uniqueid int32) (*ffProto.StBattleProp, error) {
	prop, ok := b.Props[uniqueid]
	if !ok {
		return nil, fmt.Errorf("Prop item uniqueid[%v] not exist", uniqueid)
	}
	return prop, nil
}

// 场景移除物品, 通知所有用户
func (b *battle) RemoveSceneProp(uniqueid int32) {
	// 删除场景物品
	delete(b.Props, uniqueid)

	// 广播场景物品移除
	for _, agent := range b.agents {
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRemoveProp)
		m := p.Message().(*ffProto.MsgBattleRemoveProp)
		m.Uniqueid = uniqueid
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}
}

// 场景内新增物品, 通知所有用户
func (b *battle) AddSceneProp(itemtemplateid, dropItemData int32, Position *ffProto.StVector3) {
	// 场景添加物品
	prop := b.newProp(itemtemplateid, dropItemData, Position)

	for _, one := range b.agents {
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleAddProp)
		m := p.Message().(*ffProto.MsgBattleAddProp)
		m.Prop = prop
		ffProto.SendProtoExtraDataNormal(one, p, false)
	}
}

// 装备状态改变
func (b *battle) EquipStateChanged(uniqueid int32, equipState *ffProto.StBattleEquipState) {
	// 广播用户的装备状态改变
	for _, one := range b.agents {
		if one.uniqueid != uniqueid {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleEquipState)
			m := p.Message().(*ffProto.MsgBattleEquipState)
			m.Roleuniqueid = uniqueid
			m.EquipState = equipState
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}
}

// 治疗状态改变
func (b *battle) HealStateChanged(uniqueid int32, itemtemplateid int32, state int32) {
	for _, one := range b.agents {
		if one.uniqueid != uniqueid {
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

func checkBattle(agent *battleUser) (*battle, error) {
	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		return nil, fmt.Errorf("checkBattle failed, agent[%v] uniqueid[%v] uuidBattle[%v]",
			agent.agent.UUID(), agent.uniqueid, agent.uuidBattle)
	}
	return battle, nil
}

func newBattle(uuidBattle uuid.UUID) *battle {
	battle := &battle{
		uuidBattle: uuidBattle,
		uuidTokens: make([]uuid.UUID, 0, 50),
	}
	mapBattle[uuidBattle] = battle
	return battle
}
