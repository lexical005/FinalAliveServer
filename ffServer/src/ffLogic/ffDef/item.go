package ffDef

import "ffAutoGen/ffError"

// IItemMgr 物品管理接口
type IItemMgr interface {
	// TemplateCount 模板物品数量
	TemplateCount(templateid int) int

	// AddTemplate 添加模板道具
	AddTemplate(templateid, count int) ffError.Error

	// SubTemplate 扣除模板道具，外界扣除前，应确保用户持有足够的道具数据，如果不足，则返回失败，并记录在案
	SubTemplate(templateid, count int) ffError.Error

	// Equipment 获取装备
	Equipment(instanceid int) IEquipment

	// CreateEquipment 创建装备
	CreateEquipment(templateid int) (IEquipment, ffError.Error)

	// DestroyEquipment 销毁装备
	DestroyEquipment(instanceid int) ffError.Error

	// Init 初始化
	Init()

	// Clear 清除自身, 以待重用
	Clear()
}

// IEquipment 装备实例接口
type IEquipment interface {
	// IntData 装备实例通用的数值型数据
	// dataType: 无效时, 记录在案, 且返回0
	IntData(dataType int) int

	// TemplateID 获取装备的模板
	TemplateID() int

	// Init 初始化
	Init()

	// Clear 清除自身, 以待重用
	Clear()
}
