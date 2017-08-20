package main

const (
	listenTargetUser       = "User"
	listenTargetGameServer = "GameServer"
)

// appConfig 应用程序配置
var appConfig = &applicationConfig{}

// mgrUserAgent
var mgrUserAgent = &userAgentManager{}

// // tunnelServerAgent 管理
// var serverAgentMgr = &tunnelServerAgentManager{}

// applicationQuit 进程是否要退出
var applicationQuit = false

// chApplicationQuit 用于通知goroutine进程要退出
var chApplicationQuit = make(chan struct{})

// waitServerQuit 等待所有系统退出
var waitServerQuit int32
