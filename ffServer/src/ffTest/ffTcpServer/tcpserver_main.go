package main

import (
	"ffCommon/log/log"
	"ffCommon/net/session"
	"ffCommon/net/tcpserver"
)

func main() {
	e := session.Init(120, 2)
	if e != nil {
		log.RunLogger.Println(e)
		return
	}

	s, e := tcpserver.NewServer("127.0.0.1:8765")
	if e != nil {
		log.RunLogger.Println(e)
		return
	}

	e = s.Start(am)
	if e != nil {
		log.RunLogger.Println(e)
		return
	}

	select {}
}
