package main

import (
	"ffAutoGen/ffEnum"
	"ffAutoGen/ffError"
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"math/rand"
	"strings"
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
		50101: 1,
		50102: 1,
		50201: 1,
		50202: 1,
	}
	b.PreloadScene = map[int32]int32{
		10301: 1,
		10401: 1,
		50101: 1,
		50102: 1,
		50201: 1,
		50202: 1,
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
	b.NewProp(10402, 1, &ffProto.StVector3{
		X: -4202,
		Y: 155718,
		Z: 2804251,
	})
	b.NewProp(10204, 1, &ffProto.StVector3{
		X: 963,
		Y: 155718,
		Z: 2800251,
	})
	b.NewProp(40201, 1, &ffProto.StVector3{
		X: 3963,
		Y: 155718,
		Z: 2800251,
	})

	b.NewProp(50101, 80, &ffProto.StVector3{
		X: -7414,
		Y: 155718,
		Z: 2803251,
	})
	b.NewProp(50102, 100, &ffProto.StVector3{
		X: -4202,
		Y: 155718,
		Z: 2803251,
	})
	b.NewProp(50201, 200, &ffProto.StVector3{
		X: 963,
		Y: 155718,
		Z: 2803251,
	})
	b.NewProp(50202, 220, &ffProto.StVector3{
		X: 3963,
		Y: 155718,
		Z: 2803251,
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

func (b *battle) NewProp(templateid int32, ItemData int32, position *ffProto.StVector3) *ffProto.StBattleProp {
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

// 修改持有的物品, 场景内新增扔掉的物品
func (b *battle) DropBagItem(agent *agentUser, Itemtemplateid, dropItemData int32, Position *ffProto.StVector3) error {
	// 持有的物品判定
	ownItemData, ok := agent.items[Itemtemplateid]
	if !ok {
		return fmt.Errorf("DropBagItem Itemtemplateid[%v] not in bag", Itemtemplateid)
	}

	template := ffGameConfig.ItemData.ItemTemplate[Itemtemplateid]
	if template.ItemType != ffEnum.EItemTypeArmor { // 非防具类, ownItemData为数量, 判定数量是否足够
		if ownItemData < dropItemData {
			return fmt.Errorf("DropBagItem Itemtemplateid[%v] ownItemData[%v] < dropItemData[%v]",
				Itemtemplateid, ownItemData, dropItemData)
		}

		// 持有的数量减少
		agent.items[Itemtemplateid] -= dropItemData

	} else { // 防具类, ownItemData为防具耐久, 判定耐久是否大于0, 大于0时则防具有效
		if ownItemData < 1 {
			return fmt.Errorf("DropBagItem Itemtemplateid[%v] not int bag", Itemtemplateid)
		}

		// 持有的数量减少
		agent.items[Itemtemplateid] = 0
	}

	// 场景添加物品
	prop := b.NewProp(Itemtemplateid, dropItemData, Position)
	for _, one := range b.agents {
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleAddProp)
		m := p.Message().(*ffProto.MsgBattleAddProp)
		m.Prop = prop
		ffProto.SendProtoExtraDataNormal(one, p, false)
	}
	return nil
}

// 修改装备状态, 修改持有的物品, 场景内新增扔掉的物品
func (b *battle) DropEquipProp(agent *agentUser, equipIndex int32, Position *ffProto.StVector3, EquipState *ffProto.StBattleEquipState) error {
	// 参数无效
	if equipIndex < 0 || int(equipIndex) >= len(agent.equips) || agent.equips[equipIndex] == 0 {
		return fmt.Errorf("DropEquipProp invalid equipIndex[%v] equips[%v]", equipIndex, agent.equips)
	}

	// 枪械上的配件卸下
	// 枪上的子弹状态修改

	EquipTemplateID := agent.equips[equipIndex]

	// 武器位清空
	agent.equips[equipIndex] = 0

	// 装备状态
	if EquipState != nil {
		// 2 丢弃装备位上的装备(如果装备位正在使用, 则播放丢弃动作, 否则, 背着的装备位上的装备直接消失);(装备位不变);
		EquipState.EquipIndex = equipIndex
		EquipState.EquipTemplateID = EquipTemplateID
		EquipState.EquipState = 2
	}

	// 修改持有的物品, 场景内新增扔掉的物品
	template := ffGameConfig.ItemData.ItemTemplate[EquipTemplateID]
	if template.ItemType != ffEnum.EItemTypeArmor { // 非防具类, dropItemData为数量
		return b.DropBagItem(agent, EquipTemplateID, 1, Position)
	}

	// 防具类, dropItemData为防具耐久
	return b.DropBagItem(agent, EquipTemplateID, agent.items[EquipTemplateID], Position)
}

// 捡取场景里的物品
func (b *battle) PickProp(agent *agentUser, message *ffProto.MsgBattlePickProp) error {
	prop, ok := b.Props[message.Itemuniqueid]
	if !ok {
		return fmt.Errorf("PickProp ItemUniqueID[%v] not exist", message.Itemuniqueid)
	}

	// 装备类, 捡取后立即装备
	var equip bool
	var equipIndex int32 = -1
	template := ffGameConfig.ItemData.ItemTemplate[prop.Templateid]
	if template.ItemType == ffEnum.EItemTypeGunWeapon {
		gun := ffGameConfig.ItemData.GunWeapon[prop.Templateid]
		equip = true
		if gun.GunWeaponType == ffEnum.EGunWeaponTypePistol {
			// 手枪
			equipIndex = 2
		} else {
			if agent.equips[0] == 0 {
				// 0号武器位为空
				equipIndex = 0
			} else if agent.equips[1] == 0 {
				// 1号武器位为空
				equipIndex = 1
			} else if agent.activeWeaponIndex == 0 {
				// 当前使用0号武器位
				equipIndex = 0
			} else if agent.activeWeaponIndex == 1 {
				// 当前使用1号武器位
				equipIndex = 1
			} else {
				// 装备到0号武器位
				equipIndex = 0
			}
		}
	} else if template.ItemType == ffEnum.EItemTypeMelleeWeapon {
		// 近战武器
		equip = true
		equipIndex = 3
	} else if template.ItemType == ffEnum.EItemTypeArmor {
		// 防具
		equip = true
		armor := ffGameConfig.ItemData.Armor[prop.Templateid]
		if armor.EArmorType == ffEnum.EArmorTypeVest {
			// 防弹衣位
			equipIndex = 4
		} else {
			// 头盔位
			equipIndex = 5
		}
	}

	if equip {
		/*装备状态
		0 手上武器放回背部(装备位不变);
		1 背部武器拿到手上(装备位不变);
		2 丢弃装备位上的装备(如果装备位正在使用, 则播放丢弃动作, 否则, 背着的装备位上的装备直接消失);(装备位不变);
		3 从地上捡起装备拿到手上(装备位变更为新捡起的装备的装备位)(装备位上有装备时, 由客户端维护删除);
		4 从地上捡起装备背到背部(装备位不变)(装备位上有装备时, 由客户端维护删除);
		5 切换武器位(武器位改变)
		*/
		// 修改装备状态, 修改持有的物品, 场景内新增扔掉的物品
		if agent.equips[equipIndex] != 0 {
			if err := b.DropEquipProp(agent, equipIndex, message.Position, nil); err != nil {
				return err
			}
		}

		// 装备位状态
		var state int32
		agent.equips[equipIndex] = prop.Templateid
		if template.ItemType != ffEnum.EItemTypeArmor {
			if agent.activeWeaponIndex == equipIndex { // 武器位不变, 从地上捡起武器拿到手上
				state = 3
			} else if agent.equips[agent.activeWeaponIndex] == 0 { // 原先手上为空, 从地上捡起武器, 手上切换到捡起武器所在武器位
				agent.activeWeaponIndex = equipIndex
				state = 3
			} else { // 武器位不变, 从地上捡起武器背到背部
				state = 4
			}
		} else {
			// 防具状态, 总是为4
			state = 4
		}

		// 装备状态
		EquipState := &ffProto.StBattleEquipState{
			EquipIndex:      equipIndex,
			EquipTemplateID: prop.Templateid,
			EquipState:      state,
		}

		// 捡取导致装备状态改变
		message.EquipState = EquipState

		// 广播用户的装备状态改变
		for _, one := range b.agents {
			if one.uniqueid != agent.uniqueid {
				p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleEquipState)
				m := p.Message().(*ffProto.MsgBattleEquipState)
				m.Roleuniqueid = agent.uniqueid
				m.EquipState = EquipState
				ffProto.SendProtoExtraDataNormal(one, p, false)
			}
		}
	}

	// 删除场景物品
	delete(b.Props, prop.Uniqueid)

	// 捡成功, 修改持有的物品
	if _, ok = agent.items[prop.Templateid]; ok {
		agent.items[prop.Templateid] += prop.ItemData
	} else {
		agent.items[prop.Templateid] = prop.ItemData
	}
	message.Itemtemplateid, message.ItemData = prop.Templateid, prop.ItemData

	// 广播场景物品移除
	{
		for _, agent := range b.agents {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRemoveProp)
			m := p.Message().(*ffProto.MsgBattleRemoveProp)
			m.Uniqueid = prop.Uniqueid
			ffProto.SendProtoExtraDataNormal(agent, p, false)
		}
	}

	return nil
}

// 切换武器
func (b *battle) SwitchWeapon(agent *agentUser, message *ffProto.MsgBattleSwitchWeapon) error {
	// 参数无效
	if message.EquipIndex < 0 || int(message.EquipIndex) >= len(agent.equips) || agent.activeWeaponIndex == message.EquipIndex {
		return fmt.Errorf("SwitchWeapon invalid equipIndex[%v] activeWeaponIndex[%v] equips[%v]",
			message.EquipIndex, agent.activeWeaponIndex, agent.equips)
	}

	// 当前武器位改变
	agent.activeWeaponIndex = message.EquipIndex

	// 装备状态改变
	EquipState := &ffProto.StBattleEquipState{
		EquipIndex:      message.EquipIndex,
		EquipTemplateID: agent.equips[message.EquipIndex],
		EquipState:      5,
	}
	message.EquipState = EquipState

	// 广播用户的装备状态改变
	for _, one := range b.agents {
		if one.uniqueid != agent.uniqueid {
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleEquipState)
			m := p.Message().(*ffProto.MsgBattleEquipState)
			m.Roleuniqueid = agent.uniqueid
			m.EquipState = EquipState
			ffProto.SendProtoExtraDataNormal(one, p, false)
		}
	}

	return nil
}

func (b *battle) Heal(agent *agentUser, healitemtemplateid int32, state int32) bool {
	if state == 1 {
		if agent.healitemtemplateid != 0 {
			return false
		}
		agent.healitemtemplateid = healitemtemplateid
		return true
	}

	if state == 0 {
		agent.healitemtemplateid = 0
		return true
	}

	if state == 2 {
		if agent.healitemtemplateid == 0 {
			return false
		}
		agent.healitemtemplateid = 0
		agent.health += 20

		// 血量同步
		{
			p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleRoleHealth)
			m := p.Message().(*ffProto.MsgBattleRoleHealth)
			m.Roleuniqueid = agent.uniqueid
			m.Health = agent.health
			ffProto.SendProtoExtraDataNormal(agent, p, false)
		}

		// 扣除道具(根据使用完成, 自行扣除)
		{
			// p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleDropProp)
			// m := p.Message().(*ffProto.MsgBattleDropProp)
			// m.Roleuniqueid = agent.uniqueid
			// m.Itemtemplateid = healitemtemplateid
			// m.Itemnumber = 1
			// ffProto.SendProtoExtraDataNormal(agent, p, false)
		}
	}

	return false
}

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

	// 捡取
	err := battle.PickProp(agent, message)
	message.Itemuniqueid, message.Position = 0, nil
	if err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

func onBattleProtoDropBagProp(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleDropBagProp)

	// 战场不存在或者扔失败
	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
	}
	if err := battle.DropBagItem(agent, message.Itemtemplateid, message.ItemData, message.Position); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}
	message.Position = nil

	// 扔成功
	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

func onBattleProtoDropEquipProp(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleDropEquipProp)

	// 战场不存在或者扔失败
	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
	}

	if err := battle.DropEquipProp(agent, message.EquipIndex, message.Position, message.EquipState); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}
	message.Position = nil

	// 扔成功
	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 行为
func onBattleProtoRoleAction(agent *agentUser, proto *ffProto.Proto) (result bool) {
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
func onBattleProtoRoleShootState(agent *agentUser, proto *ffProto.Proto) (result bool) {
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
func onBattleProtoRoleShoot(agent *agentUser, proto *ffProto.Proto) (result bool) {
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
func onBattleProtoRoleShootHit(agent *agentUser, proto *ffProto.Proto) (result bool) {
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
func onBattleProtoRoleEyeRotate(agent *agentUser, proto *ffProto.Proto) (result bool) {
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
func onBattleProtoRoleHeal(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleRoleHeal)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	if !battle.Heal(agent, message.Itemtemplateid, message.State) {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	for _, agent := range battle.agents {
		p := ffProto.ApplyProtoForSend(proto.ProtoID())
		m := p.Message().(*ffProto.MsgBattleRoleHeal)
		m.Roleuniqueid = message.Roleuniqueid
		m.Itemtemplateid = message.Itemtemplateid
		m.State = message.State
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
}

// 切换武器
func onBattleProtoSwitchWeapon(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleSwitchWeapon)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	err := battle.SwitchWeapon(agent, message)
	if err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

// 战斗作弊指令
func onBattleProtoCheat(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleCheat)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	if strings.ToLower(message.Cmd) == "settle" {
		battle.Settle()
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
