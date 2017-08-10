package ffExcel

import (
	"fmt"
)

// 头部
var fmtGoProtoHeader = `syntax = "proto3";
package main;

message Grammar {
    string grammar = 1;
}
`

var fmtCshapProtoHeader = `syntax = "proto3";
package NConfig;

message Grammar {
    string grammar = 1;
}
`

// 工作表
var fmtProtoExcel = `
message ExcelExportTest {%v
%v
}
`

// 工作表的工作簿数据
var fmtProtoExcelFieldList = "\n    repeated St%v %v = %v;" // repeated StVIPlist VIPlist = 1;
var fmtProtoExcelFieldMap = "\n    map<%v, St%v> %v = %v;"  // map<int32, StVIPmap> VIPmap = 2;
var fmtProtoExcelFieldStruct = "\n    St%v %v = %v;"        // StVIPstruct VIPstruct = 3;

// 工作簿
var fmtProtoSheet = `
    message St%v {%v
	}
`

// 工作簿每一列
var fmtProtoSheetLine = "\n        %v %v = %v;" //         int32 InfoInt = 1;

// 根据toml数据格式, 转换得到Proto定义
func genProtoDefineFromToml(excel *excel, exportLimit string) (goProto, csharpProto string) {
	type sheetLine struct {
		lines       []string
		mapLineType map[string]string
	}

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
					tmp.mapLineType[line.lineName] = line.lineType.ProtoType()
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
	}

	// 头部
	goProto += fmtGoProtoHeader
	csharpProto += fmtCshapProtoHeader

	// 表格内的所有工作簿
	excelSheetVar, excelSheetDef := "", ""
	for i, sheetName := range excelSheetNames {
		sheetFields := ""
		sheetLine := mapExcelSheetInfo[sheetName]
		for j, fieldName := range sheetLine.lines {
			sheetFields += fmt.Sprintf(fmtProtoSheetLine, sheetLine.mapLineType[fieldName], fieldName, j+1)
		}

		// 工作簿的结构定义
		excelSheetDef += fmt.Sprintf(fmtProtoSheet, sheetName, sheetFields)

		// 工作表的工作簿字段变量
		if excelSheetTypes[i] == sheetTypeList {
			excelSheetVar += fmt.Sprintf(fmtProtoExcelFieldList, sheetName, sheetName, i+1)
		} else if excelSheetTypes[i] == sheetTypeMap {
			excelSheetVar += fmt.Sprintf(fmtProtoExcelFieldMap, excelSheetMapKeyTypes[i], sheetName, sheetName, i+1)
		} else if excelSheetTypes[i] == sheetTypeStruct {
			excelSheetVar += fmt.Sprintf(fmtProtoExcelFieldStruct, sheetName, sheetName, i+1)
		}
	}

	// 工作表
	goProto += fmt.Sprintf(fmtProtoExcel, excelSheetDef, excelSheetVar)
	csharpProto += fmt.Sprintf(fmtProtoExcel, excelSheetDef, excelSheetVar)

	return
}
