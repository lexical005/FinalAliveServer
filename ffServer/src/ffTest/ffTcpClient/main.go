package main

import (
	"ffCommon/log/log"
	"ffCommon/net/tcpclient"
	"ffCommon/net/tcpsession"
	"ffCommon/util"
	"fmt"
	"sync/atomic"
	"time"
)

var waitQuitCount int32
var chApplicationQuit = make(chan struct{}, 1)

func main() {
	defer func() {
		util.PanicProtect("main")
		<-time.After(time.Second)
	}()

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
