package main

import (
	"time"
)

// matchGroup 匹配组
type matchGroup struct {
	// ready 准备组, 已创建战场, 正在补齐人数, 等待战场开启
	ready []*readyGroup

	// units 参与匹配的所有单元
	units []iMatchUnit

	// mode 匹配模式
	mode matchMode

	// now 当前时间
	now time.Time

	// funcSetupReadyGroup 匹配函数
	funcSetupReadyGroup func()
}

// init 初始化
func (group *matchGroup) init(mode matchMode) {
	group.mode = mode
	if mode != matchModeDouble {
		group.funcSetupReadyGroup = group.matchNormal
	} else {
		group.funcSetupReadyGroup = group.matchStranger
	}

	group.ready = make([]*readyGroup, 0, 2)
	group.units = make([]iMatchUnit, 0, appConfig.Match.InitMatchCount/4)
}

// JoinMatch 加入匹配
func (group *matchGroup) JoinMatch(unit iMatchUnit) bool {
	if group.mode == matchModeSingle {
		if unit.Count() != matchModeSingleUnitCount {
			return false
		}
	} else if group.mode == matchModeDouble {
		if unit.Count() != matchModeDoubleUnitCount {
			return false
		}
	}

	group.units = append(group.units, unit)

	return true
}

// LeaveMatch 离开匹配
func (group *matchGroup) LeaveMatch(unit iMatchUnit) bool {
	for i, one := range group.units {
		if one == unit {
			group.units = append(group.units[:i], group.units[i+1:]...)
			return true
		}
	}
	return false
}

// Match 尝试匹配
func (group *matchGroup) Match() {
	group.now = time.Now()
	for {
		// 移除已开启的战场
		{
			for i, ready := range group.ready {
				if ready.stopEnterTime.After(group.now) {
					if i > 0 {
						group.ready = group.ready[i:]
					}
					break
				}
			}
		}

		// 补齐不满员的准备组
		{
			for _, ready := range group.ready {
				index := 0

			lackLoop:
				for index < len(group.units) {
					unit := group.units[index]
					if unit.Count() <= ready.lackCount {
						// 从等待列表内移除
						group.units = append(group.units[:index], group.units[index+1:]...)

						// 加入
						if ready.Enter(unit) {
							break lackLoop
						}
					} else {
						index++
					}
				}
			}
		}

		group.funcSetupReadyGroup()
	}
}

// matchNormal 标准匹配, 单人模式, 双人模式
func (group *matchGroup) matchNormal() {
	// 建立新的准备组
	{
		for {
			var ready *readyGroup
			count := 0
			index := 0
			for index < len(group.units) {
				unit := group.units[index]
				count += unit.Count()
				if ready == nil {
					count += unit.Count()
					if count >= appConfig.Match.BattleMinPlayerCount { // 累计人数达成最低人数需求

						// 创建准备组
						ready = instReadyGroupPool.apply()
						ready.Init(appConfig.Match.ExpectMaxPlayerCount,
							group.now.Add(time.Duration(appConfig.Match.StopEnterTime)*time.Second))
						group.ready = append(group.ready, ready)

						// 更新units
						units := group.units[:index+1]
						group.units = group.units[index+1:]
						index = 0

						// 加入
						if full := ready.EnterMulti(units, count); full {
							goto fullReady
						}
					} else {
						index++
					}

				} else {
					if count <= ready.lackCount {
						// 从等待列表内移除, index值不变
						group.units = append(group.units[:index], group.units[index+1:]...)

						// 加入
						if full := ready.Enter(unit); full {
							goto fullReady
						}
					} else {
						index++
					}
				}
				continue

			fullReady:
				// 准备组, 满员, 开始新的一轮
				ready = nil
				index = 0
			}

			// 准备组, 不满员, 开始新的一轮
			if ready != nil {
				ready = nil
				index = 0
			} else {
				break
			}
		}
	}
}

// matchStranger 陌生人组队匹配, 四人模式
func (group *matchGroup) matchStranger() {
	// 建立新的准备组
	{
		for {
			var ready *readyGroup
			count := 0
			index := 0
			for index < len(group.units) {
				unit := group.units[index]
				count += unit.Count()
				if ready == nil {
					count += unit.Count()
					if count >= appConfig.Match.BattleMinPlayerCount { // 累计人数达成最低人数需求

						// 创建准备组
						ready = instReadyGroupPool.apply()
						ready.Init(appConfig.Match.ExpectMaxPlayerCount,
							group.now.Add(time.Duration(appConfig.Match.StopEnterTime)*time.Second))
						group.ready = append(group.ready, ready)

						// 更新units
						units := group.units[:index+1]
						group.units = group.units[index+1:]
						index = 0

						// 加入
						if full := ready.EnterMulti(units, count); full {
							goto fullReady
						}
					} else {
						index++
					}

				} else {
					if count <= ready.lackCount {
						// 从等待列表内移除, index值不变
						group.units = append(group.units[:index], group.units[index+1:]...)

						// 加入
						if full := ready.Enter(unit); full {
							goto fullReady
						}
					} else {
						index++
					}
				}
				continue

			fullReady:
				// 准备组, 满员, 开始新的一轮
				ready = nil
				index = 0
			}

			// 准备组, 不满员, 开始新的一轮
			if ready != nil {
				ready = nil
				index = 0
			} else {
				break
			}
		}
	}
}

func newMatchGroup(mode matchMode) *matchGroup {
	group := &matchGroup{}
	group.init(mode)
	return group
}
