package main

import (
	"ffCommon/log/log"
	"ffCommon/net/tcpclient"
	"ffCommon/net/tcpserver"
	"ffCommon/net/tcpsession"
	"ffCommon/util"
	"ffProto"
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	defer func() {
		util.PanicProtect("ffAgentServer.main")

		<-time.After(time.Second)
	}()

	// 初始化
	err := startup()
	if err != nil {
		log.FatalLogger.Println(err)
		return
	}

	// 启动
	if err = agentServerUser.start(appConfig.ServerUser); err != nil {
		log.FatalLogger.Println(err)
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
			log.RunLogger.Printf("useragent_server[%s]", agentServerUser.Status())
			tcpsession.PrintModule()
			tcpserver.PrintModule()
			ffProto.PrintModule()

			if atomic.LoadInt32(&waitServerQuit) == 0 {
				break quitLoop
			}
		}
	}

	fmt.Println("close complete")
}

func printStatus() {
	tcpsession.PrintModule()
	tcpclient.PrintModule()
	tcpserver.PrintModule()
	ffProto.PrintModule()
}
