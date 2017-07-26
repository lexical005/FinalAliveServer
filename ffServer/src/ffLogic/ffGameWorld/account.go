package ffGameWorld

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffLogic/ffDef"
)

type account struct {
	onceKick util.Once // 用于只执行一次踢出

	// uuid
	uuidAgent   uuid.UUID
	uuidAccount uuid.UUID

	// 数据
	name     string
	intDatas []int

	heroMgr  ffDef.IHeroMgr
	itemMgr  ffDef.IItemMgr
	levelMgr ffDef.ILevelMgr
	shopMgr  ffDef.IShopMgr
}

// Init 初始化
func (ac *account) Init() {
	ac.onceKick.Reset()

	// 重置数据
	for i := 0; i < len(ac.intDatas); i++ {
		ac.intDatas[i] = 0
	}
	ac.name = ""

	// 管理器初始化
	// ac.heroMgr.Init()
	// ac.itemMgr.Init()
	// ac.levelMgr.Init()
	// ac.shopMgr.Init()

	// 上线事件
	managerOfAccount.onEventOnline(ac)
}

// Clear 清除自身, 以待重用
func (ac *account) Clear() {
	// 管理器清空
	ac.heroMgr.Clear()
	ac.itemMgr.Clear()
	ac.levelMgr.Clear()
	ac.shopMgr.Clear()
}

// IntData 帐号通用的数值型数据
// dataType: 无效时, 记录在案, 且返回0
func (ac *account) IntData(dataType int) int {
	if dataType >= 0 && dataType < len(ac.intDatas) {
		return ac.intDatas[dataType]
	}

	log.FatalLogger.Printf("account.IntData: invalid dataType[%d]", dataType)

	return 0
}

// Name 帐号名称
func (ac *account) Name() string {
	return ac.name
}

// HeroMgr 获取英雄管理接口
func (ac *account) HeroMgr() ffDef.IHeroMgr {
	return ac.heroMgr
}

// ItemMgr 获取物品管理接口
func (ac *account) ItemMgr() ffDef.IItemMgr {
	return ac.itemMgr
}

// LevelMgr 获取关卡管理接口
func (ac *account) LevelMgr() ffDef.ILevelMgr {
	return ac.levelMgr
}

// ShopMgr 获取购物管理接口
func (ac *account) ShopMgr() ffDef.IShopMgr {
	return ac.shopMgr
}

// kick 踢下线
func (ac *account) kick(kickReason ffError.Error) {
	ac.onceKick.Do(func() {
		// 离线事件
		managerOfAccount.onEventOffline(ac)

		log.RunLogger.Printf("account.kick: uuidAgent[%x]-uuidAccount[%x] kickReason[%v]", ac.uuidAgent, ac.uuidAccount, kickReason)

		// 自身处理

		// 踢出
		worldFrame.Kick(ac.uuidAgent, true, kickReason)
	})
}

// kick 被踢下线, 不会发送踢出通知协议
func (ac *account) kicked(kickReason ffError.Error, notifyOffline bool) {
	ac.onceKick.Do(func() {
		// 离线事件
		if notifyOffline {
			managerOfAccount.onEventOffline(ac)
		}

		log.RunLogger.Printf("account.kicked: uuidAgent[%x]-uuidAccount[%x] kickReason[%v]", ac.uuidAgent, ac.uuidAccount, kickReason)

		// 自身处理

		// 踢出
		worldFrame.Kick(ac.uuidAgent, false, kickReason)
	})
}

func newAccount() interface{} {
	return &account{}
}
