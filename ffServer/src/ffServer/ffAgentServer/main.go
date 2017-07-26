package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
)

func main() {
	util.PanicProtect("ffAgentServer.main")

	// 初始化
	err := startup()
	if err != nil {
		log.FatalLogger.Println(err)
		return
	}

	// 启动
	if err = clientAgentMgr.start(); err != nil {
		log.FatalLogger.Println(err)
		return
	}

	if err = serverAgentMgr.start(); err != nil {
		log.FatalLogger.Println(err)
		return
	}

	// 等待关闭
	select {}
}
