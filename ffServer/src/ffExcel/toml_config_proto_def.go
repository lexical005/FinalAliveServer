package ffExcel

import (
	"fmt"
	"strconv"
	"strings"
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
message {MessageName} {{SubMessages}
{MessageMembers}
}
`

// 工作表的工作簿数据
var fmtProtoExcelFieldList = "\n    repeated St{StructType} {StructType} = {FieldIndex};"          // repeated StVIPlist VIPlist = 1;
var fmtProtoExcelFieldMap = "\n    map<{MapKeyType}, St{StructType}> {StructType} = {FieldIndex};" // map<int32, StVIPmap> VIPmap = 2;
var fmtProtoExcelFieldMapKey = "\n    repeated {MapKeyType} {StructType}Key = {FieldIndex};"       // repeated int32 VIPmapKey = 2;
var fmtProtoExcelFieldMapValue = "\n    repeated St{StructType} {StructType}Value = {FieldIndex};" // repeated StVIPmap VIPmapValue = 3;
var fmtProtoExcelFieldStruct = "\n    St{StructType} {StructType} = {FieldIndex};"                 // StVIPstruct VIPstruct = 4;

// 工作簿
var fmtProtoSheet = `
    message St%v {%v
	}
`

// 工作簿每一列
var fmtProtoSheetLine = "\n        %v %v = %v;" //         int32 InfoInt = 1;

// 根据toml数据格式, 转换得到Proto定义
func genProtoDefineFromToml(allExcels []*excel, exportLimit string) (goProto, csharpProto string) {
	type sheetLine struct {
		lines       []string
		mapLineType map[string]string
	}

	// 头部
	goProto += fmtGoProtoHeader
	csharpProto += fmtCshapProtoHeader

	for _, excel := range allExcels {
		if !excel.exportToClient() {
			continue
		}

		excelSheetNames := make([]string, 0, len(excel.sheets))
		excelSheetTypes := make([]int, 0, len(excel.sheets))
		excelSheetMapKeyTypes := make([]string, 0, len(excel.sheets))
		mapExcelSheetInfo := make(map[string]*sheetLine, len(excel.sheets))
		for _, sheet := range excel.sheets {
			if !sheet.exportToClient() {
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

		// 表格内的所有工作簿
		excelSheetVarGo, excelSheetVarCsharp, excelSheetDef := "", "", ""
		protoGoIndex, protoCSharpIndex := 0, 0
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
				protoGoIndex++
				field := strings.Replace(fmtProtoExcelFieldList, "{StructType}", sheetName, -1)
				field = strings.Replace(field, "{FieldIndex}", strconv.Itoa(protoGoIndex), -1)
				excelSheetVarGo += field

				protoCSharpIndex++
				field = strings.Replace(fmtProtoExcelFieldList, "{StructType}", sheetName, -1)
				field = strings.Replace(field, "{FieldIndex}", strconv.Itoa(protoCSharpIndex), -1)
				excelSheetVarCsharp += field
			} else if excelSheetTypes[i] == sheetTypeMap {
				protoGoIndex++
				field := strings.Replace(fmtProtoExcelFieldMapKey, "{StructType}", sheetName, -1)
				field = strings.Replace(field, "{FieldIndex}", strconv.Itoa(protoGoIndex), -1)
				field = strings.Replace(field, "{MapKeyType}", excelSheetMapKeyTypes[i], -1)
				excelSheetVarGo += field

				protoGoIndex++
				field = strings.Replace(fmtProtoExcelFieldMapValue, "{StructType}", sheetName, -1)
				field = strings.Replace(field, "{FieldIndex}", strconv.Itoa(protoGoIndex), -1)
				excelSheetVarGo += field

				protoCSharpIndex++
				field = strings.Replace(fmtProtoExcelFieldMapKey, "{StructType}", sheetName, -1)
				field = strings.Replace(field, "{FieldIndex}", strconv.Itoa(protoCSharpIndex), -1)
				field = strings.Replace(field, "{MapKeyType}", excelSheetMapKeyTypes[i], -1)
				excelSheetVarCsharp += field

				protoCSharpIndex++
				field = strings.Replace(fmtProtoExcelFieldMapValue, "{StructType}", sheetName, -1)
				field = strings.Replace(field, "{FieldIndex}", strconv.Itoa(protoCSharpIndex), -1)
				excelSheetVarCsharp += field

			} else if excelSheetTypes[i] == sheetTypeStruct {
				protoGoIndex++
				field := strings.Replace(fmtProtoExcelFieldStruct, "{StructType}", sheetName, -1)
				field = strings.Replace(field, "{FieldIndex}", strconv.Itoa(protoGoIndex), -1)
				excelSheetVarGo += field

				protoCSharpIndex++
				field = strings.Replace(fmtProtoExcelFieldStruct, "{StructType}", sheetName, -1)
				field = strings.Replace(field, "{FieldIndex}", strconv.Itoa(protoCSharpIndex), -1)
				excelSheetVarCsharp += field
			}
		}

		// 工作表
		s := strings.Replace(fmtProtoExcel, "{MessageName}", excel.name, -1)
		s = strings.Replace(s, "{SubMessages}", excelSheetDef, -1)
		s = strings.Replace(s, "{MessageMembers}", excelSheetVarGo, -1)
		goProto += s

		s = strings.Replace(fmtProtoExcel, "{MessageName}", excel.name, -1)
		s = strings.Replace(s, "{SubMessages}", excelSheetDef, -1)
		s = strings.Replace(s, "{MessageMembers}", excelSheetVarCsharp, -1)
		csharpProto += s
	}

	return
}
