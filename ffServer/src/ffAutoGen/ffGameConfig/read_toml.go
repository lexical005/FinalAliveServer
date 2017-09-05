package ffGameConfig

import (
	"ffCommon/log/log"
)

// ExcelExportTestData 配置表ExcelExportTest的数据
var ExcelExportTestData *ExcelExportTest

// ItemData 配置表Item的数据
var ItemData *Item

// LanguageData 配置表Language的数据
var LanguageData *Language

// RandBornData 配置表RandBorn的数据
var RandBornData *RandBorn

// ReadAllToml 读取所有toml配置
func ReadAllToml() (result bool) {
	result = true

	var err error

	ExcelExportTestData, err = ReadExcelExportTest()
	if err != nil {
		result = false
		log.FatalLogger.Printf("ReadExcelExportTest get error[%v]", err)
	}

	ItemData, err = ReadItem()
	if err != nil {
		result = false
		log.FatalLogger.Printf("ReadItem get error[%v]", err)
	}

	LanguageData, err = ReadLanguage()
	if err != nil {
		result = false
		log.FatalLogger.Printf("ReadLanguage get error[%v]", err)
	}

	RandBornData, err = ReadRandBorn()
	if err != nil {
		result = false
		log.FatalLogger.Printf("ReadRandBorn get error[%v]", err)
	}

	return result
}
