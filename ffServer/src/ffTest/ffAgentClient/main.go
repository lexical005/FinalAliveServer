package main

import (
	"ffCommon/log/log"
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

	// err = logfile.Init(logfile.DefaultLogFileRelativePath, logfile.DefaultLogFileLengthLimit, false, "", false, "", false, "")
	// if err != nil {
	// 	log.RunLogger.Println(err)
	// 	return
	// }

	err = agent.start("127.0.0.1:15101", true)
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
