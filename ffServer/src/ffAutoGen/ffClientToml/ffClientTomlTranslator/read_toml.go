package main

import (
	"ffAutoGen/ffClientToml"
	"ffCommon/log/log"
)

var tomlExcelExportTest *ffClientToml.ExcelExportTest

func readToml() {
	var err error

	tomlExcelExportTest, err = ffClientToml.ReadExcelExportTest()
	if err != nil {
		log.RunLogger.Printf("ReadExcelExportTest get error[%v]", err)
	}

}

func init() {
	allRead = append(allRead, readToml)
}
