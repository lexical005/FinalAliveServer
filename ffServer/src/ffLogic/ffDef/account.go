package ffDef

import (
	"ffCommon/uuid"
)

// IAccount 帐号数据接口
type IAccount interface {
	// Init 初始化
	Init()

	// Clear 清除自身, 以待重用
	Clear()

	// UUID 帐号唯一标识
	UUID() uuid.UUID

	// IntData 帐号通用的数值型数据
	// dataType: 无效时, 记录在案, 且返回0
	IntData(dataType int) int

	// Name 帐号名称
	Name() string

	// HeroMgr 获取英雄管理接口
	HeroMgr() IHeroMgr

	// ItemMgr 获取物品管理接口
	ItemMgr() IItemMgr

	// LevelMgr 获取关卡管理接口
	LevelMgr() ILevelMgr

	// ShopMgr 获取购物管理接口
	ShopMgr() IShopMgr
}
