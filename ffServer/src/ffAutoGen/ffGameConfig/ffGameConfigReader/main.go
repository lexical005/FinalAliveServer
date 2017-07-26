package main

import (
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
)

func main() {
	var err error
	_, err = ffGameConfig.ReadMall()
	if err != nil {
		log.RunLogger.Printf("ReadMall get error[%v]", err)
	}
}
