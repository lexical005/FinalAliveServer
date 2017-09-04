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

// 修改持有的物品, 场景内新增扔掉的物品
func (b *battle) DropBagItem(agent *agentUser, Itemtemplateid, Itemnumber int32, Position *ffProto.StVector3) error {
	// 持有的物品数量判定
	if number, ok := agent.items[Itemtemplateid]; !ok || number < Itemnumber {
		return fmt.Errorf("DropBagItem invalid Itemtemplateid[%v] number[%v] < Itemnumber[%v]", Itemtemplateid, number, Itemnumber)
	}

	// 持有的数量减少
	agent.items[Itemtemplateid] -= Itemnumber

	// 场景添加物品
	for _, one := range b.agents {
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleAddProp)
		m := p.Message().(*ffProto.MsgBattleAddProp)
		m.Prop = b.NewProp(Itemtemplateid, Itemnumber, Position)
		ffProto.SendProtoExtraDataNormal(one, p, false)
	}
	return nil
}

// 修改武器状态, 修改持有的物品, 场景内新增扔掉的物品
func (b *battle) DropWeaponProp(agent *agentUser, weaponIndex int32, Position *ffProto.StVector3, WeaponState *ffProto.StBattleWeaponState) error {
	// 参数无效
	if weaponIndex < 0 || int(weaponIndex) >= len(agent.weapons) || agent.weapons[weaponIndex] == 0 {
		return fmt.Errorf("DropWeaponProp invalid weaponIndex[%v] weapons[%v]", weaponIndex, agent.weapons)
	}

	// 枪械上的配件卸下
	// 枪上的子弹状态修改

	Itemtemplateid := agent.weapons[weaponIndex]

	// 武器位清空
	agent.weapons[weaponIndex] = 0

	// 武器状态
	if WeaponState != nil {
		// 2 丢弃武器位上的武器(如果武器位正在使用, 则播放丢弃动作, 否则, 背着的武器位上的武器直接消失);
		WeaponState.WeaponIndex = weaponIndex
		WeaponState.WeaponTemplateID = Itemtemplateid
		WeaponState.WeaponState = 2
	}

	// 修改持有的物品, 场景内新增扔掉的物品
	return b.DropBagItem(agent, Itemtemplateid, 1, Position)
}

// 捡取场景里的物品
func (b *battle) PickProp(agent *agentUser, message *ffProto.MsgBattlePickProp) error {
	prop, ok := b.Props[message.Itemuniqueid]
	if !ok {
		return fmt.Errorf("PickProp ItemUniqueID[%v] not exist", message.Itemuniqueid)
	}

	// 枪械和近战武器, 捡取后立即装备
	var equip bool
	var equipIndex int32 = -1
	var gun *ffGameConfig.GunWeapon
	template := ffGameConfig.ItemData.ItemTemplate[prop.Templateid]
	if template.ItemType == ffEnum.EItemTypeGunWeapon {
		gun = ffGameConfig.ItemData.GunWeapon[prop.Templateid]
		equip = true
		if gun.GunWeaponType == ffEnum.EGunWeaponTypePistol {
			// 手枪
			equipIndex = 2
		} else {
			if agent.weapons[0] == 0 {
				// 0号武器位为空
				equipIndex = 0
			} else if agent.weapons[1] == 0 {
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
	}

	if equip {
		/*武器状态
		0 手上武器放回背部(武器位不变);
		1 背部武器拿到手上(武器位不变);
		2 丢弃武器位上的武器(如果武器位正在使用, 则播放丢弃动作, 否则, 背着的武器位上的武器直接消失);(武器位不变);
		3 从地上捡起武器拿到手上(武器位变更为新捡起的武器的武器位);
		4 从地上捡起武器背到背部(武器位不变);
		5 切换武器位(武器位改变)
		*/
		// 修改武器状态, 修改持有的物品, 场景内新增扔掉的物品
		if agent.weapons[equipIndex] != 0 {
			if err := b.DropWeaponProp(agent, equipIndex, message.Position, nil); err != nil {
				return err
			}
		}

		// 武器位状态
		var state int32
		agent.weapons[equipIndex] = prop.Templateid
		if agent.activeWeaponIndex == -1 { // 未使用武器位
			agent.activeWeaponIndex = equipIndex
			state = 3 // 武器位改变, 从地上捡起武器拿到手上
		} else {
			if agent.activeWeaponIndex == equipIndex {
				state = 3 // 武器位不变, 从地上捡起武器拿到手上
			} else {
				if agent.weapons[agent.activeWeaponIndex] == 0 {
					agent.activeWeaponIndex = equipIndex
					state = 3 // 武器位改变, 从地上捡起武器拿到手上
				} else {
					state = 4 // 武器位不变, 从地上捡起武器背到背部
				}
			}
		}

		// 武器状态
		WeaponState := &ffProto.StBattleWeaponState{
			WeaponIndex:      equipIndex,
			WeaponTemplateID: prop.Templateid,
			WeaponState:      state,
		}

		// 捡取导致武器状态改变
		message.WeaponState = WeaponState

		// 广播用户的武器状态改变
		for _, one := range b.agents {
			if one.uniqueid != agent.uniqueid {
				p := ffProto.ApplyProtoForSend(ffProto.MessageType_BattleWeaponState)
				m := p.Message().(*ffProto.MsgBattleWeaponState)
				m.Roleuniqueid = agent.uniqueid
				m.WeaponState = WeaponState
				ffProto.SendProtoExtraDataNormal(one, p, false)
			}
		}
	}

	// 删除场景物品
	delete(b.Props, prop.Uniqueid)

	// 捡成功, 修改持有的物品
	if _, ok = agent.items[prop.Templateid]; ok {
		agent.items[prop.Templateid] += prop.Number
	} else {
		agent.items[prop.Templateid] = prop.Number
	}
	message.Itemtemplateid, message.Itemnumber = prop.Templateid, prop.Number

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
	message.Roleuniqueid, message.Itemuniqueid, message.Position = 0, 0, nil
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
	if err := battle.DropBagItem(agent, message.Itemtemplateid, message.Itemnumber, message.Position); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}
	message.Roleuniqueid = 0
	message.Position = nil

	// 扔成功
	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

func onBattleProtoDropWeaponProp(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleDropWeaponProp)

	// 战场不存在或者扔失败
	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
	}

	if err := battle.DropWeaponProp(agent, message.WeaponIndex, message.Position, message.WeaponState); err != nil {
		message.Result = ffError.ErrUnknown.Code()
		log.RunLogger.Println(err)
	}
	message.Roleuniqueid, message.Position = 0, nil

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
		m.Endtag = message.Endtag
		m.Endposition = message.Endposition
		m.Endnormal = message.Endnormal
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	if message.Targetuniqueid != 0 {
		battle.ShootHit(agent, message.Shootid, message.Targetuniqueid)
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

// 战斗武器状态
func onBattleProtoWeaponState(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgBattleWeaponState)

	battle, ok := mapBattle[agent.uuidBattle]
	if !ok {
		message.Result = ffError.ErrUnknown.Code()
		return ffProto.SendProtoExtraDataNormal(agent, proto, true)
	}

	for _, agent := range battle.agents {
		p := ffProto.ApplyProtoForSend(proto.ProtoID())
		m := p.Message().(*ffProto.MsgBattleWeaponState)
		m.Roleuniqueid = message.Roleuniqueid
		ffProto.SendProtoExtraDataNormal(agent, p, false)
	}

	return
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
