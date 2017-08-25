package main

import "time"

// readyGroup 准备组, 正在补齐人数, 筹备战场, 等待战场开启
type readyGroup struct {
	units         []iMatchUnit // 涉及的所有匹配单元
	lackCount     int          // 距离满员还差多少玩家
	stopEnterTime time.Time    // 战场最晚开启时间
}

// Init 初始化准备组
func (group *readyGroup) Init(fullCount int, stopEnterTime time.Time) {
	group.lackCount, group.stopEnterTime = fullCount, stopEnterTime
}

// Enter 一个匹配单元进入此准备组
//	返回值, 是否满员
func (group *readyGroup) Enter(unit iMatchUnit) bool {
	group.units = append(group.units, unit)
	group.lackCount -= unit.Count()

	unit.MatchSuccess()

	return group.lackCount == 0
}

// EnterMulti 多个个匹配单元进入此准备组
//	返回值, 是否满员
func (group *readyGroup) EnterMulti(units []iMatchUnit, count int) bool {
	group.units = append(group.units, units...)
	group.lackCount -= count

	for _, unit := range units {
		unit.MatchSuccess()
	}

	return group.lackCount == 0
}

func newReadyGroup() *readyGroup {
	return &readyGroup{}
}
