package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"net"

	_ "github.com/lexical005/gosproto"
)

func main() {
	ln, err := net.Listen("tcp", ":12315")
	if err != nil {
		log.RunLogger.Println(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.RunLogger.Println(err)
		}

		go util.SafeGo(handleConn, conn)
	}
}
