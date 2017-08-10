package main

import (
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
	"ffCommon/util"
)

var tomlExcelExportTest *ffGameConfig.ExcelExportTest

func read() {
	defer util.PanicProtect()

	var err error

	tomlExcelExportTest, err = ffGameConfig.ReadExcelExportTest()
	if err != nil {
		log.RunLogger.Printf("ReadExcelExportTest get error[%v]", err)
	}

}
