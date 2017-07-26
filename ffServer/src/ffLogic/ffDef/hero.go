package ffDef

// IHeroMgr 英雄管理接口
type IHeroMgr interface {
	// CreateHero 创建英雄
	CreateHero(templateid int)

	// Hero 当前所使用英雄的接口
	Hero() IHero

	// Heros 返回所拥有的英雄列表
	Heros() []IHero

	// Init 初始化
	Init()

	// Clear 清除自身, 以待重用
	Clear()
}

// IHero 英雄数据接口
type IHero interface {
	// IntData 英雄通用的数值型数据
	// dataType: 无效时, 记录在案, 且返回0
	IntData(dataType int) int

	// Name 英雄名称
	Name() string

	// TemplateID 英雄模板
	TemplateID() int

	// Equipments 穿戴的所有装备
	Equipments() []IEquipment

	// Init 初始化
	Init()

	// Clear 清除自身, 以待重用
	Clear()
}
