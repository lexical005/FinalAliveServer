package main

import (
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffLogic/ffDef"
	"time"

	"fmt"
)

// gameWorldTimer 定时器
type gameWorldTimer struct {
	uuid uuid.UUID

	stop       bool  // 是否已停止
	nonosecond int64 // 下一次触发时间

	firstOffset time.Duration   // 首次触发偏移，单位: 毫秒
	interval    time.Duration   // 触发间隔，单位: 毫秒
	count       int             // 触发剩余次数. -1 代表不限次数. >0: 代表指定限定次数
	timerFunc   ffDef.TimerFunc // 触发时回调函数
}

// FirstOffset 首次触发偏移
func (gwt *gameWorldTimer) FirstOffset() time.Duration {
	return gwt.firstOffset
}

// Interval 触发间隔
func (gwt *gameWorldTimer) Interval() time.Duration {
	return gwt.interval
}

// LeftCount 剩余触发次数. -1 代表不限次数
func (gwt *gameWorldTimer) LeftCount() int {
	return gwt.count
}

func (gwt *gameWorldTimer) String() string {
	return fmt.Sprintf("gameWorldTimer: firstOffset[%dms] interval[%dms] count[%d]",
		gwt.firstOffset/time.Millisecond, gwt.interval/time.Millisecond, gwt.count)
}

// gameWorldTimeManager 游戏框架--时间管理
type gameWorldTimeManager struct {
	// 当前时间
	now time.Time

	// 定时器：通过2组定时器的机制，来回翻转
	indexTimers int
	timers      [2][]*gameWorldTimer
}

// init 初始化
func (gwtm *gameWorldTimeManager) init() {
	gwtm.now = time.Now()

	gwtm.indexTimers = 0
	gwtm.timers = [2][]*gameWorldTimer{
		make([]*gameWorldTimer, 0, 16),
		make([]*gameWorldTimer, 0, 16),
	}
}

// AddTimer 设置定时器, 单位: 毫秒
//  firstOffset: 首次触发偏移
//  interval: 触发间隔
//  count: 触发次数. -1 代表不限次数. >0: 代表指定限定次数
//  timerFunc: 触发时回调函数
// 返回值: 定时器对象 ITimer
func (gwtm *gameWorldTimeManager) AddTimer(firstOffset time.Duration, interval time.Duration, count int, timerFunc ffDef.TimerFunc) ffDef.ITimer {
	firstOffset *= time.Millisecond
	interval *= time.Millisecond

	timer := &gameWorldTimer{
		uuid: worldFrame.UUID(ffDef.UUIDTimer),

		stop:       false,
		nonosecond: gwtm.now.UnixNano(),

		firstOffset: firstOffset,
		interval:    interval,
		count:       count,
		timerFunc:   timerFunc,
	}

	if firstOffset > 0 {
		timer.nonosecond += firstOffset.Nanoseconds()
	} else {
		timer.nonosecond += interval.Nanoseconds()
	}

	log.RunLogger.Printf("gameWorldTimeManager.AddTimer: timer[%v]\n", timer)

	gwtm.addTimer(timer, gwtm.indexTimers)
	return timer
}

// StopTimer 停止定时器
//  timer: 定时器对象
func (gwtm *gameWorldTimeManager) StopTimer(timer ffDef.ITimer) {
	log.RunLogger.Printf("gameWorldTimeManager.StopTimer: timer[%v]\n", timer)

	for i := 0; i < len(gwtm.timers); i++ {
		for j := 0; j < len(gwtm.timers[i]); j++ {
			if timer == gwtm.timers[i][j] {
				log.RunLogger.Printf("gameWorldTimeManager.StopTimer: found timer[%x]\n", gwtm.timers[i][j].uuid)

				gwtm.timers[i][j].stop = true
				return
			}
		}
	}
}

//
func (gwtm *gameWorldTimeManager) String() string {
	s := "gameWorldTimeManager:\n"
	for i := 0; i < len(gwtm.timers); i++ {
		for j := 0; j < len(gwtm.timers[i]); j++ {
			s += fmt.Sprintf("[%d]-[%v]\n", i+1, gwtm.timers[i][j])
		}
	}
	return s
}

// updateTime 时间更新，更新间隔，设计为worldTimeUpdateInterval
func (gwtm *gameWorldTimeManager) updateTime() {
	// 更新当前时间及实际循环间隔
	now := time.Now()
	nanosecond := now.UnixNano()
	gwtm.now = now

	// 定时器更新
	indexTimers := gwtm.indexTimers
	gwtm.indexTimers = (gwtm.indexTimers + 1) % len(gwtm.timers)
	timers := gwtm.timers[gwtm.indexTimers]

	i, l := 0, len(timers)
	for i = 0; i < l; i++ {
		timer := timers[i]
		// 到达触发时间
		if timer.nonosecond <= nanosecond {
			// 未被停止的定时器
			if !timer.stop {
				// 限定次数
				if timer.count > 0 {
					timer.count--
				}

				// 触发
				timer.timerFunc(timer)

				// 添加到另一列表内
				if timer.count != 0 && !timer.stop {

					// 触发时间加上触发间隔
					timer.nonosecond += timer.interval.Nanoseconds()

					// 添加到另一列表内
					gwtm.addTimer(timer, indexTimers)
				}
			}
		} else {
			break
		}
	}

	//
	if i < l {
		gwtm.timers[gwtm.indexTimers] = timers[i:]
	} else if l > 0 {
		gwtm.timers[gwtm.indexTimers] = timers[:0]
	}
}

// addTimer
func (gwtm *gameWorldTimeManager) addTimer(timer *gameWorldTimer, indexTimers int) {
	timers := gwtm.timers[indexTimers]
	l := len(timers)
	timers = append(timers, timer)
	gwtm.timers[indexTimers] = timers

	for i := 0; i < l; i++ {
		if timer.nonosecond < timers[i].nonosecond || timer.nonosecond == timers[i].nonosecond && timer.uuid < timers[i].uuid {
			copy(timers[i+1:], timers[i:l])
			timers[i] = timer
			return
		}
	}
}
