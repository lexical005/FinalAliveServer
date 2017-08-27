package main

import (
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
)

var tomlAsset *ffGameConfig.Asset
var tomlExcelExportTest *ffGameConfig.ExcelExportTest
var tomlItem *ffGameConfig.Item
var tomlLanguage *ffGameConfig.Language

func readToml() {
	var err error

	tomlAsset, err = ffGameConfig.ReadAsset()
	if err != nil {
		log.RunLogger.Printf("ReadAsset get error[%v]", err)
	}

	tomlExcelExportTest, err = ffGameConfig.ReadExcelExportTest()
	if err != nil {
		log.RunLogger.Printf("ReadExcelExportTest get error[%v]", err)
	}

	tomlItem, err = ffGameConfig.ReadItem()
	if err != nil {
		log.RunLogger.Printf("ReadItem get error[%v]", err)
	}

	tomlLanguage, err = ffGameConfig.ReadLanguage()
	if err != nil {
		log.RunLogger.Printf("ReadLanguage get error[%v]", err)
	}

}

func init() {
	allRead = append(allRead, readToml)
}
