package main

import (
	"ffCommon/log/log"
	"ffLogic/ffDef"
	"time"

	"testing"
)

var timeMgr *gameWorldTimeManager

func timerCallback1(timer ffDef.ITimer) {
	log.RunLogger.Println("timerCallback1", timer)
}

func timerCallback2(timer ffDef.ITimer) {
	log.RunLogger.Println("timerCallback2", timer)
}

func goTime() {
	timeMgr.AddTimer(0, 33, 100, timerCallback1)
	timeMgr.AddTimer(1, 33, 100, timerCallback2)
	log.RunLogger.Println(timeMgr)

	t := 0
	for {
		select {
		case <-time.After(worldTimeUpdateInterval):
			timeMgr.updateTime()
		}

		t++
		if t%1010 == 0 {
			log.RunLogger.Println(timeMgr)
		}
	}
}

func Test_Timer1(t *testing.T) {
	timeMgr = &gameWorldTimeManager{}
	timeMgr.init()

	go goTime()

	select {
	case <-time.After(time.Second * 5):
	}

	log.RunLogger.Println("Test_Timer1 passed")
}
