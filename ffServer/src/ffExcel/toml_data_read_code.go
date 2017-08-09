package ffExcel

import (
	"fmt"
	"path"
	"strings"
)

// package
var fmtPackage = `package %v

`

// import
var fmtImport = `
import (
	"ffCommon/util"

	"fmt"

	"github.com/lexical005/toml"
)

`
var fmtImportGrammar = `
import (
	"ffCommon/util"
	"ffLogic/ffGrammar"

	"fmt"

	"github.com/lexical005/toml"
)

`

// excel struct define
var fmtExcelComment = `// %v excel %v
`
var fmtExcelDefStart = `type %v struct {
`
var fmtExcelDefFieldList = `   %v []%v
`
var fmtExcelDefFieldMap = `   %v map[%v]%v
`
var fmtExcelDefFieldStruct = `   %v %v
`
var fmtExcelDefEnd = `}

`

// excel String method
var fmtExcelMethodStringStart = `func (%v %v) String() string {`
var fmtExcelMethodStringContentStart = `
	s := ""`
var fmtExcelMethodStringContentLoopList = `
	s += "%v"
	for _, row := range %v.%v {
		s += fmt.Sprintf("%%v\n", row)
	}
`
var fmtExcelMethodStringContentLoopMap = `
	s += "%v"
	for k, v := range %v.%v {
		s += fmt.Sprintf("%%v:%%v\n", k, v)
	}
`
var fmtExcelMethodStringContentLoopStruct = `
	s += "%v"
	s += fmt.Sprintf("%%v\n", %v.%v)
`
var fmtExcelMethodStringContentEnd = `
	return s`
var fmtExcelMethodStringEnd = `
}

`

// excel sheet struct define
var fmtSheetComment = `// %v sheet %v of excel %v
`
var fmtSheetDefStart = `type %v struct {
`
var fmtSheetDefField = `	%v  %v
`
var fmtSheetDefEnd = `}

`

// excel sheet String method
var fmtSheetMethodStringStart = `func (%v %v) String() string {`
var fmtSheetMethodStringContentStart = `
	s := "["`
var fmtSheetMethodStringContentLoop = `
	s += fmt.Sprintf("%v:%%v,", %v.%v)`
var fmtSheetMethodStringContentEnd = `
	s += "]"
	return s`
var fmtSheetMethodStringEnd = `
}

`

var fmtReadFunc = `
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
	result += fmt.Sprintf(fmtPackage, exportConfig.packageName)

	// import
	if hasGrammar {
		result += fmtImportGrammar
	} else {
		result += fmtImport
	}

	// excel struct define
	result += fmt.Sprintf(fmtExcelComment, excelName, excelName)
	result += fmt.Sprintf(fmtExcelDefStart, excelName)
	for i, sheetName := range excelSheetNames {
		if excelSheetTypes[i] == sheetTypeList {
			result += fmt.Sprintf(fmtExcelDefFieldList, sheetName, sheetName)
		} else if excelSheetTypes[i] == sheetTypeMap {
			result += fmt.Sprintf(fmtExcelDefFieldMap, sheetName, excelSheetMapKeyTypes[i], sheetName)
		} else if excelSheetTypes[i] == sheetTypeStruct {
			result += fmt.Sprintf(fmtExcelDefFieldStruct, sheetName, sheetName)
		}
	}
	result += fmtExcelDefEnd

	// excel String method
	result += fmt.Sprintf(fmtExcelMethodStringStart, shortExcelName, excelName)
	result += fmtExcelMethodStringContentStart
	for i, sheetName := range excelSheetNames {
		if excelSheetTypes[i] == sheetTypeList {
			result += fmt.Sprintf(fmtExcelMethodStringContentLoopList, sheetName, shortExcelName, sheetName)
		} else if excelSheetTypes[i] == sheetTypeMap {
			result += fmt.Sprintf(fmtExcelMethodStringContentLoopMap, sheetName, shortExcelName, sheetName)
		} else if excelSheetTypes[i] == sheetTypeStruct {
			result += fmt.Sprintf(fmtExcelMethodStringContentLoopStruct, sheetName, shortExcelName, sheetName)
		}
	}
	result += fmtExcelMethodStringContentEnd
	result += fmtExcelMethodStringEnd

	// excel sheet struct define
	for _, sheetName := range excelSheetNames {
		shortSheetName := getShortName(sheetName)

		result += fmt.Sprintf(fmtSheetComment, sheetName, sheetName, excelName)
		result += fmt.Sprintf(fmtSheetDefStart, sheetName)
		for _, fieldName := range mapExcelSheetInfo[sheetName].lines {
			result += fmt.Sprintf(fmtSheetDefField, fieldName, mapExcelSheetInfo[sheetName].mapLineType[fieldName])
		}
		result += fmtSheetDefEnd

		// excel sheet String method
		result += fmt.Sprintf(fmtSheetMethodStringStart, shortSheetName, sheetName)
		result += fmtSheetMethodStringContentStart
		for _, fieldName := range mapExcelSheetInfo[sheetName].lines {
			result += fmt.Sprintf(fmtSheetMethodStringContentLoop, fieldName, shortSheetName, fieldName)
		}
		result += fmtSheetMethodStringContentEnd
		result += fmtSheetMethodStringEnd
	}

	// read excel
	if exportLimit == "server" {
		result += fmt.Sprintf(fmtReadFunc,
			excelName, excelName,
			excelName, shortExcelName, excelName,
			path.Join(exportConfig.ServerReadTomlDataPath, fmt.Sprintf("%v.toml", excelName)),
			shortExcelName, excelName,
			shortExcelName)
	} else if exportLimit == "client" {
		result += fmt.Sprintf(fmtReadFunc,
			excelName, excelName,
			excelName, shortExcelName, excelName,
			path.Join("toml", "client", fmt.Sprintf("%v.toml", excelName)),
			shortExcelName, excelName,
			shortExcelName)
	}
	return result
}
