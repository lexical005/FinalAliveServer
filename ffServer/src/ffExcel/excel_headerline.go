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
}

// 本列配置，在什么情况下忽略
func (lh *headerLine) ignore() bool {
	return lh.lineType.IsIgnore() || (!lh.lineHome.client && !lh.lineHome.server) || lh.lineName == ""
}

// 本列配置，是否导出到服务端
func (lh *headerLine) exportToServer() bool {
	return !lh.ignore() && lh.lineHome.server
}

// 本列配置，是否导出到客户端
func (lh *headerLine) exportToClient() bool {
	return !lh.ignore() && lh.lineHome.client
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
	}, nil
}
