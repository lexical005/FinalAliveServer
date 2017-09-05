package ffExcel

import (
	"ffCommon/util"
	"path"
	"strings"
)

var fmtTomlFile = `package ffGameConfig

import (
	"ffCommon/log/log"
)
{AllTomeData}
// ReadAllToml 读取所有toml配置
func ReadAllToml() (result bool) {
	result = true

	var err error
{AllTomlReader}
	initHelper()

	return result
}
`

var fmtTomlVarFile = `
// {ExcelName}Data 配置表{ExcelName}的数据
var {ExcelName}Data *{ExcelName}
`

var fmtTomlReadFile = `
	{ExcelName}Data, err = Read{ExcelName}()
	if err != nil {
		result = false
		log.FatalLogger.Printf("Read{ExcelName} get error[%v]", err)
	}
`

func genReadAllTomlCode(allExcels []*excel) {
	AllTomeData, AllTomlReader := "", ""
	for _, excel := range allExcels {
		// 导出读取toml数据的Go代码
		if excel.exportToServer() && excel.exportType == "config" {
			AllTomeData += strings.Replace(fmtTomlVarFile, "{ExcelName}", excel.name, -1)
			AllTomlReader += strings.Replace(fmtTomlReadFile, "{ExcelName}", excel.name, -1)
		}
	}

	result := strings.Replace(fmtTomlFile, "{AllTomeData}", AllTomeData, -1)
	result = strings.Replace(result, "{AllTomlReader}", AllTomlReader, -1)

	util.WriteFile(path.Join(exportConfig.ServerExportGoCodePath, "read_toml.go"), []byte(result))
}
