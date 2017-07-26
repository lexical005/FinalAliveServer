package ffDef

// IShopMgr 购物管理接口
type IShopMgr interface {
	// Shop 获取特定商店购物接口
	Shop(templateid int) IShop

	// Init 初始化
	Init()

	// Clear 清除自身, 以待重用
	Clear()
}

// IShop 商店购物接口
type IShop interface {

	// Init 初始化
	Init()

	// Clear 清除自身, 以待重用
	Clear()
}
