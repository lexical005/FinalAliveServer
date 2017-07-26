package main

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"time"
)

func triggerPanic(params ...interface{}) {
	<-time.After(time.Second * 2)

	s, _ := params[0].(string)
	log.RunLogger.Println(s)

	a := 10
	b := 100
	c := a / b
	log.RunLogger.Println(b / c)
}

func main() {
	defer util.PanicProtect("main")

	go util.SafeGo(triggerPanic, "triggerPanic")

	// 等待结束
	select {}
}
