package main

import (
	"ffAutoGen/ffEnum"
	"ffAutoGen/ffGameConfig"
	"ffProto"
	"fmt"
)

type itemManager struct {
	items             map[int32]int32 // 拥有的物品数据
	equips            []int32         // 装备信息(武器,防弹衣,头盔)
	activeWeaponIndex int32           // 手上武器下标
}

func (m *itemManager) HasItem(itemtemplateid int32) bool {
	if ItemData, ok := m.items[itemtemplateid]; ok {
		return ItemData > 0
	}
	return false
}

// 修改持有的物品
func (m *itemManager) DropBagItem(itemtemplateid, dropItemData int32) error {
	// 持有的物品判定
	ownItemData, ok := m.items[itemtemplateid]
	if !ok {
		return fmt.Errorf("DropBagItem itemtemplateid[%v] not in bag", itemtemplateid)
	}

	template := ffGameConfig.ItemData.ItemTemplate[itemtemplateid]
	if template.ItemType != ffEnum.EItemTypeArmor { // 非防具类, ownItemData为数量, 判定数量是否足够
		if ownItemData < dropItemData {
			return fmt.Errorf("DropBagItem itemtemplateid[%v] ownItemData[%v] < dropItemData[%v]",
				itemtemplateid, ownItemData, dropItemData)
		}

		// 持有的数量减少
		m.items[itemtemplateid] -= dropItemData

	} else { // 防具类, ownItemData为防具耐久, 判定耐久是否大于0, 大于0时则防具有效
		if ownItemData < 1 {
			return fmt.Errorf("DropBagItem itemtemplateid[%v] not int bag", itemtemplateid)
		}

		// 持有的数量减少
		m.items[itemtemplateid] = 0
	}
	return nil
}

// 修改装备状态, 修改持有的物品
func (m *itemManager) DropEquip(agent *battleUser, equipIndex int32, position *ffProto.StVector3, equipState *ffProto.StBattleEquipState) error {
	// 参数无效
	if equipIndex < 0 || int(equipIndex) >= len(m.equips) || m.equips[equipIndex] == 0 {
		return fmt.Errorf("DropEquip invalid equipIndex[%v] equips[%v]", equipIndex, m.equips)
	}

	// 枪械上的配件卸下
	// 枪上的子弹状态修改

	equipTemplateID := m.equips[equipIndex]

	// 武器位清空
	m.equips[equipIndex] = 0

	// 装备状态
	if equipState != nil {
		// 2 丢弃装备位上的装备(如果装备位正在使用, 则播放丢弃动作, 否则, 背着的装备位上的装备直接消失);(装备位不变);
		equipState.EquipIndex = equipIndex
		equipState.EquipTemplateID = equipTemplateID
		equipState.EquipState = 2
	}

	// 修改持有的物品, 场景内新增扔掉的物品
	template := ffGameConfig.ItemData.ItemTemplate[equipTemplateID]
	if template.ItemType != ffEnum.EItemTypeArmor { // 非防具类, dropItemData为数量
		return agent.DropBagItem(equipTemplateID, 1, position)
	}

	// 防具类, dropItemData为防具耐久
	return agent.DropBagItem(equipTemplateID, m.items[equipTemplateID], position)
}

// 捡取场景里的物品
func (m *itemManager) PickProp(agent *battleUser, prop *ffProto.StBattleProp, message *ffProto.MsgBattlePickProp) error {
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
			if m.equips[0] == 0 {
				// 0号武器位为空
				equipIndex = 0
			} else if m.equips[1] == 0 {
				// 1号武器位为空
				equipIndex = 1
			} else if m.activeWeaponIndex == 0 {
				// 当前使用0号武器位
				equipIndex = 0
			} else if m.activeWeaponIndex == 1 {
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
		if m.equips[equipIndex] != 0 {
			if err := m.DropEquip(agent, equipIndex, message.Position, nil); err != nil {
				return err
			}
		}

		// 装备位状态
		var state int32
		m.equips[equipIndex] = prop.Templateid
		if template.ItemType != ffEnum.EItemTypeArmor {
			if m.activeWeaponIndex == equipIndex { // 武器位不变, 从地上捡起武器拿到手上
				state = 3
			} else if m.equips[m.activeWeaponIndex] == 0 { // 原先手上为空, 从地上捡起武器, 手上切换到捡起武器所在武器位
				m.activeWeaponIndex = equipIndex
				state = 3
			} else { // 武器位不变, 从地上捡起武器背到背部
				state = 4
			}
		} else {
			// 防具状态, 总是为4
			state = 4
		}

		// 捡取导致装备状态改变
		message.EquipState = &ffProto.StBattleEquipState{
			EquipIndex:      equipIndex,
			EquipTemplateID: prop.Templateid,
			EquipState:      state,
		}
	} else {
		message.EquipState = nil
	}

	// 捡成功, 修改持有的物品
	if _, ok := m.items[prop.Templateid]; ok {
		m.items[prop.Templateid] += prop.ItemData
	} else {
		m.items[prop.Templateid] = prop.ItemData
	}
	message.Itemtemplateid, message.ItemData = prop.Templateid, prop.ItemData

	return nil
}

// 切换武器
func (m *itemManager) SwitchWeapon(message *ffProto.MsgBattleSwitchWeapon) error {
	// 参数无效
	if message.EquipIndex < 0 || int(message.EquipIndex) >= len(m.equips) || m.activeWeaponIndex == message.EquipIndex {
		return fmt.Errorf("SwitchWeapon invalid equipIndex[%v] activeWeaponIndex[%v] equips[%v]",
			message.EquipIndex, m.activeWeaponIndex, m.equips)
	}

	// 当前武器位改变
	m.activeWeaponIndex = message.EquipIndex

	// 装备状态改变
	message.EquipState = &ffProto.StBattleEquipState{
		EquipIndex:      message.EquipIndex,
		EquipTemplateID: m.equips[message.EquipIndex],
		EquipState:      5,
	}

	return nil
}

func newItemManager() *itemManager {
	return &itemManager{
		items:  make(map[int32]int32, 16),
		equips: make([]int32, 6),
	}
}
