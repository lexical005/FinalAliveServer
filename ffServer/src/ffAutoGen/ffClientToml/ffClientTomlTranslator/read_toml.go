package main

import (
	"ffAutoGen/ffClientToml"
	"ffCommon/log/log"
)

var tomlAsset *ffClientToml.Asset
var tomlExcelExportTest *ffClientToml.ExcelExportTest
var tomlItem *ffClientToml.Item
var tomlLanguage *ffClientToml.Language

func readToml() {
	var err error

	tomlAsset, err = ffClientToml.ReadAsset()
	if err != nil {
		log.RunLogger.Printf("ReadAsset get error[%v]", err)
	}

	tomlExcelExportTest, err = ffClientToml.ReadExcelExportTest()
	if err != nil {
		log.RunLogger.Printf("ReadExcelExportTest get error[%v]", err)
	}

	tomlItem, err = ffClientToml.ReadItem()
	if err != nil {
		log.RunLogger.Printf("ReadItem get error[%v]", err)
	}

	tomlLanguage, err = ffClientToml.ReadLanguage()
	if err != nil {
		log.RunLogger.Printf("ReadLanguage get error[%v]", err)
	}

}

func init() {
	allRead = append(allRead, readToml)
}
