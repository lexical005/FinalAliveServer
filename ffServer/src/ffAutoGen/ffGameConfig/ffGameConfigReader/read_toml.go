package main

import (
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
)

var tomlExcelExportTest *ffGameConfig.ExcelExportTest

func readToml() {
	var err error

	tomlExcelExportTest, err = ffGameConfig.ReadExcelExportTest()
	if err != nil {
		log.RunLogger.Printf("ReadExcelExportTest get error[%v]", err)
	}

}

func init() {
	allRead = append(allRead, readToml)
}
