package main

import (
	"ffCommon/log/log"
	"ffCommon/net/tcpclient"
	"ffCommon/net/tcpsession"
	"ffCommon/util"
	"ffProto"
	"fmt"
	"sync/atomic"
	"time"
)

var waitQuitCount int32
var chApplicationQuit = make(chan struct{}, 1)

func main() {
	defer util.PanicProtect(func(isPanic bool) {
		if isPanic {
			log.RunLogger.Println("异常退出, 以上是错误堆栈")
			<-time.After(time.Hour)
		}
	}, "ffTcpClient")

	// logfile.InitRunLog("", logfile.DefaultLogFileLengthLimit, logfile.DefaultLogFileRunPrefix)
	log.RunLogger = log.NewLoggerEmpty()

	tcpsession.Init(tcpsession.DefaultReadDeadTime, tcpsession.DefaultOnlineCount, tcpsession.DefaultInitSessionNetEventDataCount)
	tcpclient.Init()

	client := &tcpClient{}
	client.start("127.0.0.1:15101")

	select {
	case <-chApplicationQuit:
		break
	}

	log.RunLogger.Printf("closing")
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

			if atomic.LoadInt32(&waitQuitCount) == 0 {
				break quitLoop
			}
		}
	}

	fmt.Println("close complete")
}

func printStatus() {
	tcpsession.PrintModule()
	tcpclient.PrintModule()
	ffProto.PrintModule()
}
