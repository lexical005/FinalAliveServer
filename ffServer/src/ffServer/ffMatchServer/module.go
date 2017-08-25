package main

type matchMode int32

const (
	// matchModeSingle 单人模式
	matchModeSingle matchMode = 1

	// matchModeDouble 双人模式
	matchModeDouble matchMode = 2

	// matchModeFour 四人模式
	matchModeFour matchMode = 4

	// matchModeCount 匹配模式数量
	matchModeCount = 3

	// matchModeSingleUnitCount 单人模式-单元人数
	matchModeSingleUnitCount = 1

	// matchModeDoubleUnitCount 双人模式-单元人数
	matchModeDoubleUnitCount = 2

	// maxTeamMemberCount 队伍最大成员数
	maxTeamMemberCount = 4

	// uuidTeamNone 无队伍
	uuidTeamNone = 0
)

// appConfig 应用程序配置
var appConfig = &applicationConfig{}

// instAgentGameServerMgr AgentGameServer管理器
var instAgentGameServerMgr = &agentGameServerManager{}

// instMatchPlayerMgr 匹配玩家管理器
var instMatchPlayerMgr = &matchPlayerManager{}

// instMatchMgr 匹配管理器
var instMatchMgr = &matchManager{}

// instReadyGroupPool 准备组管理器
var instReadyGroupPool *readyGroupPool

// waitApplicationQuit 等待所有系统退出
var waitApplicationQuit int32

// chApplicationQuit 用于通知goroutine进程要退出
var chApplicationQuit = make(chan struct{})
