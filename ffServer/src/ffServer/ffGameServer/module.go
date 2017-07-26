package main

import (
	"ffLogic/ffDef"
)

// appConfig 应用程序配置
var appConfig = &applicationConfig{}

// agentServerMgr 管理与AgentServer的连接
var agentServerMgr = &agentServer{}

// 游戏世界框架支持
var worldFrame = &gameWorldFrame{}

// 游戏世界
var world ffDef.IGameWorld
