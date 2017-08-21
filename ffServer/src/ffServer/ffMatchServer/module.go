package main

const (
	listenTargetUser       = "User"
	listenTargetGameServer = "GameServer"
)

// appConfig 应用程序配置
var appConfig = &applicationConfig{}

// agentUserSvr
var agentUserSvr = &agentGameServerServer{}

// waitApplicationQuit 等待所有系统退出
var waitApplicationQuit int32

// chApplicationQuit 用于通知goroutine进程要退出
var chApplicationQuit = make(chan struct{})
