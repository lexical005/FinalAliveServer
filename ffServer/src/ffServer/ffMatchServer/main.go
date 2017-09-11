package main

import (
	"ffCommon/log/log"
	"ffCommon/net/tcpclient"
	"ffCommon/net/tcpserver"
	"ffCommon/util"
	"ffProto"
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var err error
	defer util.PanicProtect(func(isPanic bool) {
		if isPanic {
			log.RunLogger.Println("异常退出, 以上是错误堆栈")
			<-time.After(time.Hour)
		} else if err != nil {
			util.PrintPanicStack(err)
			log.RunLogger.Println("启动出错, 以上是错误堆栈")
			<-time.After(time.Hour)
		}
	}, "ffMatchServer")

	// 初始化并启动
	err = startup()
	if err != nil {
		return
	}

	// 等待进程关闭通知
	<-chApplicationQuit

	// 等待所有服务关闭
	waitQuit()
}

func waitQuit() {
	closeTime := 0

	// 关闭中
quitLoop:
	for {
		select {
		case <-time.After(time.Second):
			closeTime++
			log.RunLogger.Printf("closing %v", closeTime)
			log.RunLogger.Printf("useragent_server[%s]", instAgentGameServerMgr.Status())
			tcpserver.PrintModule()
			ffProto.PrintModule()

			if atomic.LoadInt32(&waitApplicationQuit) == 0 {
				break quitLoop
			}
		}
	}

	fmt.Println("close complete")
}

func printStatus() {
	tcpclient.PrintModule()
	tcpserver.PrintModule()
	ffProto.PrintModule()
}
