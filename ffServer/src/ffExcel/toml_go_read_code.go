package ffExcel

import (
	"fmt"
	"path"
	"strings"
)

// package
var fmtGoPackage = `package %v

`

// import
var fmtGoImport = `
import (
	"ffCommon/util"

	"fmt"

	"github.com/lexical005/toml"
)

`
var fmtGoImportGrammar = `
import (
	"ffCommon/util"
	"ffLogic/ffGrammar"

	"fmt"

	"github.com/lexical005/toml"
)

`

// excel struct define
var fmtGoExcelComment = `// %v excel %v
`
var fmtGoExcelDefStart = `type %v struct {
`
var fmtGoExcelDefFieldList = `   %v []%v
`
var fmtGoExcelDefFieldMap = `   %v map[%v]%v
`
var fmtGoExcelDefFieldStruct = `   %v %v
`
var fmtGoExcelDefEnd = `}

`

// excel String method
var fmtGoExcelMethodStringStart = `func (%v *%v) String() string {`
var fmtGoExcelMethodStringContentStart = `
	s := ""`
var fmtGoExcelMethodStringContentLoopList = `
	s += "%v"
	for _, row := range %v.%v {
		s += fmt.Sprintf("%%v\n", row)
	}
`
var fmtGoExcelMethodStringContentLoopMap = `
	s += "%v"
	for k, v := range %v.%v {
		s += fmt.Sprintf("%%v:%%v\n", k, v)
	}
`
var fmtGoExcelMethodStringContentLoopStruct = `
	s += "%v"
	s += fmt.Sprintf("%%v\n", %v.%v)
`
var fmtGoExcelMethodStringContentEnd = `
	return s`
var fmtGoExcelMethodStringEnd = `
}

`

// excel Name method
var fmtGoExcelMethodName = `// Name the toml config's name
func (%v *%v) Name() string {
	return "%v"
}

`

// excel sheet struct define
var fmtGoSheetComment = `// %v sheet %v of excel %v
`
var fmtGoSheetDefStart = `type %v struct {
`
var fmtGoSheetDefField = `	%v  %v
`
var fmtGoSheetDefEnd = `}

`

// excel sheet String method
var fmtGoSheetMethodStringStart = `func (%v *%v) String() string {`
var fmtGoSheetMethodStringContentStart = `
	s := "["`
var fmtGoSheetMethodStringContentLoop = `
	s += fmt.Sprintf("%v:%%v,", %v.%v)`
var fmtGoSheetMethodStringContentEnd = `
	s += "]"
	return s`
var fmtGoSheetMethodStringEnd = `
}

`

var fmtGoReadFunc = `
// Read%v read excel %v
func Read%v() (%v *%v, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("%v")
	if err != nil {
		return
	}

	// 解析
	%v = &%v{}
	err = toml.Unmarshal(fileContent, %v)
	if err != nil {
		return
	}

	return
}
`

func getShortName(name string) string {
	upperName := strings.ToUpper(name)
	for i := 0; i < len(upperName); i++ {
		if upperName[i] == name[i] {
			continue
		}
		return strings.ToLower(name[:i])
	}

	return strings.ToLower(name)
}

// 得到可读取toml数据的go语言代码
func genTomlDataReadCode(excel *excel, exportConfig *ExportConfig, exportLimit string) string {
	excelName := excel.name
	shortExcelName := getShortName(excelName)

	type sheetLine struct {
		lines       []string
		mapLineType map[string]string
	}

	hasGrammar := false
	excelSheetNames := make([]string, 0, len(excel.sheets))
	excelSheetTypes := make([]int, 0, len(excel.sheets))
	excelSheetMapKeyTypes := make([]string, 0, len(excel.sheets))
	mapExcelSheetInfo := make(map[string]*sheetLine, len(excel.sheets))
	for _, sheet := range excel.sheets {
		tmp := &sheetLine{
			lines:       make([]string, 0, len(sheet.header.lines)),
			mapLineType: make(map[string]string, len(sheet.header.lines)),
		}

		for _, line := range sheet.header.lines {
			if (exportLimit == "server" && line.exportToServer() || exportLimit == "client" && line.exportToClient()) && !line.isMapKey() {
				if _, ok := tmp.mapLineType[line.lineName]; !ok {
					tmp.lines = append(tmp.lines, line.lineName)
					tmp.mapLineType[line.lineName] = line.lineType.Type()
				}
			}
		}
		mapExcelSheetInfo[sheet.name] = tmp

		if len(mapExcelSheetInfo) > 0 {
			excelSheetNames = append(excelSheetNames, sheet.name)
			excelSheetTypes = append(excelSheetTypes, sheet.sheetType)

			if sheet.header.hasMapKey() {
				excelSheetMapKeyTypes = append(excelSheetMapKeyTypes, sheet.header.mapKeyType())
			} else {
				excelSheetMapKeyTypes = append(excelSheetMapKeyTypes, "")
			}
		}

		if sheet.header.hasGrammar() {
			hasGrammar = true
		}
	}

	result := ""

	// package
	if exportLimit == "server" {
		result += fmt.Sprintf(fmtGoPackage, exportConfig.serverPackageName)
	} else if exportLimit == "client" {
		result += fmt.Sprintf(fmtGoPackage, exportConfig.clientPackageName)
	}

	// import
	if hasGrammar {
		result += fmtGoImportGrammar
	} else {
		result += fmtGoImport
	}

	// excel struct define
	result += fmt.Sprintf(fmtGoExcelComment, excelName, excelName)
	result += fmt.Sprintf(fmtGoExcelDefStart, excelName)
	for i, sheetName := range excelSheetNames {
		if excelSheetTypes[i] == sheetTypeList {
			result += fmt.Sprintf(fmtGoExcelDefFieldList, sheetName, sheetName)
		} else if excelSheetTypes[i] == sheetTypeMap {
			result += fmt.Sprintf(fmtGoExcelDefFieldMap, sheetName, excelSheetMapKeyTypes[i], sheetName)
		} else if excelSheetTypes[i] == sheetTypeStruct {
			result += fmt.Sprintf(fmtGoExcelDefFieldStruct, sheetName, sheetName)
		}
	}
	result += fmtGoExcelDefEnd

	// excel String method
	result += fmt.Sprintf(fmtGoExcelMethodStringStart, shortExcelName, excelName)
	result += fmtGoExcelMethodStringContentStart
	for i, sheetName := range excelSheetNames {
		if excelSheetTypes[i] == sheetTypeList {
			result += fmt.Sprintf(fmtGoExcelMethodStringContentLoopList, sheetName, shortExcelName, sheetName)
		} else if excelSheetTypes[i] == sheetTypeMap {
			result += fmt.Sprintf(fmtGoExcelMethodStringContentLoopMap, sheetName, shortExcelName, sheetName)
		} else if excelSheetTypes[i] == sheetTypeStruct {
			result += fmt.Sprintf(fmtGoExcelMethodStringContentLoopStruct, sheetName, shortExcelName, sheetName)
		}
	}
	result += fmtGoExcelMethodStringContentEnd
	result += fmtGoExcelMethodStringEnd

	// excel Name method
	result += fmt.Sprintf(fmtGoExcelMethodName, shortExcelName, excelName, excelName)

	// excel sheet struct define
	for _, sheetName := range excelSheetNames {
		shortSheetName := getShortName(sheetName)

		result += fmt.Sprintf(fmtGoSheetComment, sheetName, sheetName, excelName)
		result += fmt.Sprintf(fmtGoSheetDefStart, sheetName)
		for _, fieldName := range mapExcelSheetInfo[sheetName].lines {
			result += fmt.Sprintf(fmtGoSheetDefField, fieldName, mapExcelSheetInfo[sheetName].mapLineType[fieldName])
		}
		result += fmtGoSheetDefEnd

		// excel sheet String method
		result += fmt.Sprintf(fmtGoSheetMethodStringStart, shortSheetName, sheetName)
		result += fmtGoSheetMethodStringContentStart
		for _, fieldName := range mapExcelSheetInfo[sheetName].lines {
			result += fmt.Sprintf(fmtGoSheetMethodStringContentLoop, fieldName, shortSheetName, fieldName)
		}
		result += fmtGoSheetMethodStringContentEnd
		result += fmtGoSheetMethodStringEnd
	}

	// read excel
	if exportLimit == "server" {
		result += fmt.Sprintf(fmtGoReadFunc,
			excelName, excelName,
			excelName, shortExcelName, excelName,
			path.Join(exportConfig.ServerReadTomlDataPath, fmt.Sprintf("%v.toml", excelName)),
			shortExcelName, excelName,
			shortExcelName)
	} else if exportLimit == "client" {
		result += fmt.Sprintf(fmtGoReadFunc,
			excelName, excelName,
			excelName, shortExcelName, excelName,
			path.Join("toml", "client", fmt.Sprintf("%v.toml", excelName)),
			shortExcelName, excelName,
			shortExcelName)
	}
	return result
}
