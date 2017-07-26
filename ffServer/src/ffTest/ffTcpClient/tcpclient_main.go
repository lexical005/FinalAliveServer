package main

import (
	"ffCommon/log/log"
	"ffCommon/log/logfile"
	"ffCommon/net/session"

	"math/rand"
	"time"
)

func main() {
	err := session.Init(120, 2)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	err = logfile.Init(false, "", false, "", false, "")
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	err = tc1.start("127.0.0.1:2547", true)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	err = tc2.start("127.0.0.1:2547", true)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		<-time.After(time.Second * time.Duration(r.Intn(5)+5))
		session.PrintModule()
	}
}
