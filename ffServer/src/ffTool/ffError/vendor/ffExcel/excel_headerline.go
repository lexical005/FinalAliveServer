package ffExcel

import (
	"cellvalue"

	"fmt"
)

type headerLine struct {
	lineDesc string
	lineName string
	lineType cellvalue.ValueType
	lineHome valueHome
}

// 本列配置，在什么情况下忽略
func (lh *headerLine) ignore() bool {
	return lh.lineType.IsIgnore() || (!lh.lineHome.client && !lh.lineHome.server) || lh.lineName == ""
}

// 本列配置，是否到处到服务端
func (lh *headerLine) exportToServer() bool {
	return !lh.ignore() && lh.lineHome.server
}

// 本列配置，是否到处到客户端
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

func (lh *headerLine) String() string {
	return fmt.Sprintf("[[%v][%v][%v][%v] ignore[%v]]", lh.lineDesc, lh.lineName, lh.lineType.Type(), lh.lineHome, lh.ignore())
}

func newLineHeader(lineDesc, lineName, lineType, lineHome string) (*headerLine, error) {
	vt, err := cellvalue.NewValueType(lineType)
	if err != nil {
		return nil, err
	}
	home := newValueHome(lineHome)

	return &headerLine{
		lineDesc: lineDesc,
		lineName: lineName,
		lineType: vt,
		lineHome: home,
	}, nil
}
