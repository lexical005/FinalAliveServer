package main

import (
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
	"ffCommon/util"
)

func main() {
	defer util.PanicProtect()

	var err error

	_, err = ffGameConfig.ReadExcelExportTest()
	if err != nil {
		log.RunLogger.Printf("ReadExcelExportTest get error[%v]", err)
	}

}
