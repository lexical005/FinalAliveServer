package ffExcel

import (
	"cellvalue"

	"fmt"
)

type headerLine struct {
	lineDesc     string
	lineName     string
	lineType     cellvalue.ValueType
	lineRequired *valueRequired
	lineHome     *valueHome

	configExport bool
}

// 本列配置，在什么情况下忽略
func (lh *headerLine) ignore() bool {
	return lh.lineType.IsIgnore() || (!lh.lineHome.client && !lh.lineHome.server) || lh.lineName == ""
}

// 本列配置，是否导出到服务端
func (lh *headerLine) exportToServer() bool {
	return !lh.ignore() && lh.lineHome.server && lh.configExport
}

// 本列配置，是否导出到客户端
func (lh *headerLine) exportToClient() bool {
	return !lh.ignore() && lh.lineHome.client && lh.configExport
}

// 本列配置，是不是字典的主键列
func (lh *headerLine) isMapKey() bool {
	return lh.lineName == sheetTypeMapKeyName
}

// 本列配置，是不是grammar列
func (lh *headerLine) isGrammar() bool {
	return lh.lineType.IsGrammar()
}

// 本列配置，是不是必须配值
func (lh *headerLine) isRequired() bool {
	return lh.lineRequired.required
}

// limitByExportConfig 根据配置, 再调整内部数据
func (lh *headerLine) limitByExportConfig(excelName, sheetName string) {
	for _, limit := range exportConfig.ExcelExportLimit {
		if limit.Excel == excelName && limit.Sheet == sheetName {
			// 外界配置文件, 限定是否导出
			export := false
			for _, line := range limit.ExportLines {
				if line == lh.lineName {
					export = true
					break
				}
			}

			// 不导出
			if !export {
				if lh.exportToClient() || lh.exportToServer() {
					// 调整为不导出
					lh.configExport = false
				}
			}

			// 列名重命名
			for i, line := range limit.ExportLinesRenameFrom {
				if line == lh.lineName {
					lh.lineName = limit.ExportLinesRenameTo[i]
					break
				}
			}
		}
	}
}

func (lh *headerLine) String() string {
	return fmt.Sprintf("[[%v][%v][%v][%v][%v] ignore[%v]]",
		lh.lineDesc, lh.lineName, lh.lineType.Type(), lh.lineRequired, lh.lineHome, lh.ignore())
}

func newLineHeader(lineDesc, lineName, lineType, lineRequired, lineHome string) (*headerLine, error) {
	vt, err := cellvalue.NewValueType(lineType)
	if err != nil {
		return nil, err
	}

	home := newValueHome(lineHome)
	if home == nil {
		return nil, fmt.Errorf("invalid lineHome[%v]", lineHome)
	}

	required := newValueRequired(lineRequired)
	if required == nil {
		return nil, fmt.Errorf("invalid lineRequired[%v]", lineRequired)
	}

	return &headerLine{
		lineDesc:     lineDesc,
		lineName:     lineName,
		lineType:     vt,
		lineRequired: required,
		lineHome:     home,

		configExport: true,
	}, nil
}
