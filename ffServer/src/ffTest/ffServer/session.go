package main

import (
	"ffCommon/log/log"
	"net"
)

func handleConn(params ...interface{}) {
	conn, _ := params[0].(net.Conn)

	defer conn.Close()

	log.RunLogger.Println(conn.RemoteAddr().String())
	data := make([]byte, 128)
	for {
		c, err := conn.Read(data)
		if err != nil {
			log.RunLogger.Println(err)
			break
		}
		log.RunLogger.Println(string(data[:c]))
		conn.Write(data)
	}
}
