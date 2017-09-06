package main

import (
	"ffAutoGen/ffEnum"
	"ffAutoGen/ffGameConfig"
	"ffProto"
	"fmt"
	"math/rand"
)

type bornGroup struct {
	Positions []*ffProto.StVector3
}

type battleBorn struct {
	// randRolePreparePositions 准备期间, 角色随机出生点
	randRolePreparePositions []*ffProto.StVector3

	// randItemPrepareGroups 准备期间, 道具刷新组. key: area*1000+group
	randItemPrepareGroups map[int32]*bornGroup
}

func (b *battleBorn) GenItemPrepareGroup(config *ffGameConfig.BornPrepareItem, Rand *rand.Rand, battle *battle) (err error) {
	key := config.Area*1000 + config.Group

	ItemBase, ok := ffGameConfig.RandBornData.ItemBase[config.ItemBase]
	if !ok {
		err = fmt.Errorf("GenItemPrepareGroup failed, ItemBase[%v] not exist",
			config.ItemBase)
		return
	}

	group, ok := b.randItemPrepareGroups[key]
	if !ok {
		err = fmt.Errorf("GenItemPrepareGroup failed, Area[%v] Group[%v] not exist in randItemPrepareGroups",
			config.Area, config.Group)
		return
	}

	indexGroup := 0
	for i := 0; i < len(ItemBase.Chances); i++ {
		Chance := ItemBase.Chances[i]
		Item := ItemBase.Items[i]
		Number := ItemBase.Numbers[i]
		if Chance == 100 || Rand.Intn(100) < int(Chance) {

			template := ffGameConfig.ItemData.ItemTemplate[Item]
			if template.ItemType != ffEnum.EItemTypeArmor {
				battle.newProp(Item, Number, group.Positions[indexGroup])
			} else {
				armor := ffGameConfig.ItemData.Armor[Item]
				battle.newProp(Item, armor.Attrs[ffEnum.EAttrDurable], group.Positions[indexGroup])
			}

			indexGroup++
			if indexGroup >= len(group.Positions) {
				indexGroup = 0
			}
		}
	}

	return
}

var instBattleBorn = &battleBorn{}

func initBattleBorn() {
	instBattleBorn.randRolePreparePositions = make([]*ffProto.StVector3, 0, maxRoleCount)
	instBattleBorn.randItemPrepareGroups = make(map[int32]*bornGroup, 1024)

	for _, one := range ffGameConfig.RandBornData.BornPosition {
		if one.Type == ffEnum.EBornTypeRandRolePrepare {
			for i := 0; i < len(one.Positions)/3; i++ {
				instBattleBorn.randRolePreparePositions = append(instBattleBorn.randRolePreparePositions, &ffProto.StVector3{
					X: int64(one.Positions[3*i]),
					Y: int64(one.Positions[3*i+1]),
					Z: int64(one.Positions[3*i+2]),
				})
			}
		} else if one.Type == ffEnum.EBornTypeRandItemPrepare {
			key := one.Area*1000 + one.Group
			group, ok := instBattleBorn.randItemPrepareGroups[key]
			if !ok {
				group = &bornGroup{
					Positions: make([]*ffProto.StVector3, 0, len(one.Positions)/3),
				}
				instBattleBorn.randItemPrepareGroups[key] = group
			}
			for i := 0; i < len(one.Positions)/3; i++ {
				group.Positions = append(group.Positions, &ffProto.StVector3{
					X: int64(one.Positions[3*i]),
					Y: int64(one.Positions[3*i+1]),
					Z: int64(one.Positions[3*i+2]),
				})
			}
		}
	}
}
