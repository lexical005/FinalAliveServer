package main

// appConfig 应用程序配置
var appConfig = &applicationConfig{}

// mgrAgentGameServer
var mgrAgentGameServer = &agentGameServerManager{}

// applicationQuit 进程是否要退出
var applicationQuit = false

// chApplicationQuit 用于通知goroutine进程要退出
var chApplicationQuit = make(chan struct{})

// waitServerQuit 等待所有系统退出
var waitServerQuit int32