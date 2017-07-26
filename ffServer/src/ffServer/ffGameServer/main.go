package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"ffLogic/ffGameWorld"
)

func main() {
	util.PanicProtect("ffGameServer.main")

	// 初始化
	err := startup()
	if err != nil {
		log.FatalLogger.Println(err)
		return
	}

	// 创建游戏世界
	world, err = ffGameWorld.NewGameWorld(worldFrame)
	if err != nil {
		log.FatalLogger.Println(err)
		return
	}

	// 启动连接
	if err = agentServerMgr.start(); err != nil {
		log.FatalLogger.Println(err)
		return
	}

	go util.SafeGo(worldFrame.mainLoop)

	// 等待关闭
	select {}
}
