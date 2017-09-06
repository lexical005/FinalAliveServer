package main

import (
	"ffAutoGen/ffEnum"
	"ffAutoGen/ffGameConfig"
	"fmt"
	"time"
)

type healManager struct {
	itemtemplateid     int32     // 正在使用的治疗物品模板id
	healActionEndTime  time.Time // heal action 结束的时间
	KeepRecoverEndTime time.Time // 持续恢复结束的时间
	KeepRecoverTimer   time.Time // 上次恢复的时间
}

// 治疗开始
func (mgr *healManager) HealStart(agent *battleUser, itemtemplateid int32) error {
	if mgr.itemtemplateid != 0 {
		return fmt.Errorf("agent[%v] HealStart failed, itemtemplateid[%v]conflict with healing itemtemplateid[%v]",
			agent.UUID(), itemtemplateid, mgr.itemtemplateid)
	}

	// 拥有物品
	if !agent.itemManager.HasItem(itemtemplateid) {
		return fmt.Errorf("agent[%v] HealStart failed, not own itemtemplateid[%v]", agent.UUID(), itemtemplateid)
	}

	// 物品是治疗消耗品
	template := ffGameConfig.ItemData.ItemTemplate[itemtemplateid]
	if template.ItemType != ffEnum.EItemTypeConsumable {
		return fmt.Errorf("agent[%v] HealStart failed, itemtemplateid[%v] not a consumable item", agent.UUID(), itemtemplateid)
	}

	consumable := ffGameConfig.ItemData.Consumable[itemtemplateid]

	// 使用时血量限制
	if !agent.attrManager.IsHealthLessThan(consumable.UseHpLimit) {
		return fmt.Errorf("agent[%v] HealStart failed, itemtemplateid[%v] UseHpLimit[%v] not less than health[%v]",
			agent.UUID(), itemtemplateid, consumable.UseHpLimit, agent.attrManager.health)
	}

	// 记录
	mgr.itemtemplateid, mgr.healActionEndTime = itemtemplateid, time.Now().Add(time.Duration(consumable.UseTime)*time.Millisecond)
	return nil
}

// 治疗中断
func (mgr *healManager) HealCancel(agent *battleUser) error {
	if mgr.itemtemplateid == 0 {
		return fmt.Errorf("agent[%v] HealCancel failed, not in healing", agent.UUID())
	}

	// 记录
	mgr.itemtemplateid = 0
	return nil
}

// 治疗结算
func (mgr *healManager) HealSettle(agent *battleUser) error {
	if mgr.itemtemplateid == 0 {
		return fmt.Errorf("agent[%v] HealSettle failed, not in healing", agent.UUID())
	}

	// 记录
	mgr.itemtemplateid = 0
	return nil
}

func newHealManager() *healManager {
	return &healManager{}
}
