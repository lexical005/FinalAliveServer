package main

import (
	"ffAutoGen/ffClientToml"
	"ffCommon/log/log"
	"ffCommon/util"
)

var tomlExcelExportTest *ffClientToml.ExcelExportTest

func read() {
	defer util.PanicProtect()

	var err error

	tomlExcelExportTest, err = ffClientToml.ReadExcelExportTest()
	if err != nil {
		log.RunLogger.Printf("ReadExcelExportTest get error[%v]", err)
	}

}
