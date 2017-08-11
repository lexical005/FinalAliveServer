package main

import (
	"ffCommon/log/log"
	"ffCommon/net/tcpserver"
	"ffCommon/net/tcpsession"
	"ffProto"

	"math/rand"
	"time"
)

func main() {
	var onlineCount = 2

	err := tcpsession.Init(
		tcpsession.DefaultReadDeadTime,
		onlineCount,
		tcpsession.DefaultInitSessionNetEventDataCount)
	if err != nil {
		log.FatalLogger.Println(err)
		return
	}

	err = tcpserver.Init()
	if err != nil {
		log.FatalLogger.Println(err)
		return
	}

	gServer, err = newServerAgent("127.0.0.1:15101", ffProto.ExtraDataTypeNormal)
	if err != nil {
		log.FatalLogger.Println(err)
		return
	}

	if err = gServer.Start(); err != nil {
		log.FatalLogger.Println(err)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		<-time.After(time.Second * time.Duration(r.Intn(5)+5))

		tcpserver.PrintModule()
		tcpsession.PrintModule()
		ffProto.PrintModule()
	}
}
