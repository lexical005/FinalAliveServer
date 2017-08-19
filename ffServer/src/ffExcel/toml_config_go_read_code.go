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
	{ImportGrammar}
	{ImportEnum}

	"fmt"

	"github.com/lexical005/toml"
)


`

// excel struct define
var fmtGoExcelComment = `// %v excel %v
`
var fmtGoExcelDefStart = `type %v struct {
`
var fmtGoExcelDefFieldList = `   %v []*%v
`
var fmtGoExcelDefFieldMap = `   %v map[%v]*%v
`
var fmtGoExcelDefFieldStruct = `   %v %v
`
var fmtGoExcelDefEnd = `}

`

// excel String method
var fmtGoExcelMethodStringStart = `func (%v *%v) String() string {`
var fmtGoExcelMethodStringContentStart = `
	result := ""`
var fmtGoExcelMethodStringContentLoopList = `
	result += "%v"
	for _, row := range %v.%v {
		result += fmt.Sprintf("%%v\n", row)
	}
`
var fmtGoExcelMethodStringContentLoopMap = `
	result += "%v"
	for k, v := range %v.%v {
		result += fmt.Sprintf("%%v:%%v\n", k, v)
	}
`
var fmtGoExcelMethodStringContentLoopStruct = `
	result += "%v"
	result += fmt.Sprintf("%%v\n", %v.%v)
`
var fmtGoExcelMethodStringContentEnd = `
	return result`
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
	result := "["`
var fmtGoSheetMethodStringContentLoop = `
	result += fmt.Sprintf("%v:%%v,", %v.%v)`
var fmtGoSheetMethodStringContentEnd = `
	result += "]"
	return result`
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
	{MapTrans}
	return
}
`

var fmtMapTransInMap = `
	for _, one := range {ShortExcelName}.{SheetName} {
		one.{MemberName} = make(map[{MapKeyType}]{MapValueType}, len(one.{MemberName}Key))
		for index, v := range one.{MemberName}Key {
			one.{MemberName}[v] = one.{MemberName}Value[index]
		}
	}
`

var fmtMapTransInList = `
	for _, one := range {ShortExcelName}.{SheetName} {
		one.{MemberName} = make(map[{MapKeyType}]{MapValueType}, len(one.{MemberName}Key))
		for index, v := range one.{MemberName}Key {
			one.{MemberName}[v] = one.{MemberName}Value[index]
		}
	}
`

var fmtMapTransInInst = `
	one.{MemberName} = make(map[{MapKeyType}]{MapValueType}, len({ShortExcelName}.{SheetName}.{MemberName}Key))
	for index, v := range {ShortExcelName}.{SheetName}.{MemberName}Key {
		one.{MemberName}[v] = one.{MemberName}Value[index]
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

	MapTrans := ""

	hasGrammar, hasEnum := false, false
	excelSheetNames := make([]string, 0, len(excel.sheets))
	excelSheetTypes := make([]int, 0, len(excel.sheets))
	excelSheetMapKeyTypes := make([]string, 0, len(excel.sheets))
	mapExcelSheetInfo := make(map[string]*sheetLine, len(excel.sheets))
	for _, sheet := range excel.sheets {
		if exportLimit == "server" && !sheet.exportToServer() {
			continue
		} else if exportLimit == "client" && !sheet.exportToClient() {
			continue
		}

		tmp := &sheetLine{
			lines:       make([]string, 0, len(sheet.header.lines)),
			mapLineType: make(map[string]string, len(sheet.header.lines)),
		}

		for _, line := range sheet.header.lines {
			if (exportLimit == "server" && line.exportToServer() || exportLimit == "client" && line.exportToClient()) && !line.isMapKey() {
				if _, ok := tmp.mapLineType[line.lineName]; !ok {
					tmp.lines = append(tmp.lines, line.lineName)
					tmp.mapLineType[line.lineName] = line.lineType.GoType()
					if line.lineType.IsMap() {
						tmp.lines = append(tmp.lines, line.lineName+"Key")
						tmp.mapLineType[line.lineName+"Key"] = line.lineType.MapKeyGoType()

						tmp.lines = append(tmp.lines, line.lineName+"Value")
						tmp.mapLineType[line.lineName+"Value"] = line.lineType.MapValueGoType()

						var s = ""
						if sheet.sheetType == sheetTypeMap {
							s = fmtMapTransInMap
						} else if sheet.sheetType == sheetTypeList {
							s = fmtMapTransInList
						} else {
							s = fmtMapTransInInst
						}
						s = strings.Replace(s, "{ShortExcelName}", shortExcelName, -1)
						s = strings.Replace(s, "{SheetName}", sheet.name, -1)
						s = strings.Replace(s, "{MemberName}", line.lineName, -1)
						s = strings.Replace(s, "{MapKeyType}", line.lineType.MapKeyGoType()[len("[]"):], -1)
						s = strings.Replace(s, "{MapValueType}", line.lineType.MapValueGoType()[len("[]"):], -1)
						MapTrans += s
					}
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
		if sheet.header.hasEnum() {
			hasEnum = true
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
	s := fmtGoImport
	if hasGrammar {
		s = strings.Replace(s, "{ImportGrammar}", `"ffLogic/ffGrammar"`, -1)
	} else {
		s = strings.Replace(s, "{ImportGrammar}", ``, -1)
	}

	if hasEnum {
		s = strings.Replace(s, "{ImportEnum}", `"ffAutoGen/ffEnum"`, -1)
	} else {
		s = strings.Replace(s, "{ImportEnum}", ``, -1)
	}
	result += s

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
	readExcel := ""
	if exportLimit == "server" {
		readExcel = fmt.Sprintf(fmtGoReadFunc,
			excelName, excelName,
			excelName, shortExcelName, excelName,
			path.Join(exportConfig.ServerReadTomlDataPath, fmt.Sprintf("%v.toml", excelName)),
			shortExcelName, excelName,
			shortExcelName)
	} else if exportLimit == "client" {
		readExcel = fmt.Sprintf(fmtGoReadFunc,
			excelName, excelName,
			excelName, shortExcelName, excelName,
			path.Join("toml", "client", fmt.Sprintf("%v.toml", excelName)),
			shortExcelName, excelName,
			shortExcelName)
	}
	readExcel = strings.Replace(readExcel, "{MapTrans}", MapTrans, -1)
	result += readExcel

	return result
}
