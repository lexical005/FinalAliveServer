package ffExcel

import (
	"fmt"
	"strings"

	"github.com/lexical005/xlsx"
)

type sheetHeader struct {
	lines []*headerLine
}

// 判定当前配置的列里面，是否有grammar列
func (h *sheetHeader) hasGrammar() bool {
	for _, line := range h.lines {
		if line.isGrammar() {
			return true
		}
	}
	return false
}

// 判定当前配置的列里面，是否有sheetTypeMapKeyName列
func (h *sheetHeader) hasMapKey() bool {
	for _, line := range h.lines {
		if line.isMapKey() {
			return true
		}
	}
	return false
}

// 仅当hasMapKey返回true时，才应该调用此接口，以获取sheetTypeMapKeyName列的类型
func (h *sheetHeader) mapKeyType() string {
	for _, line := range h.lines {
		if line.isMapKey() {
			return line.lineType.Type()
		}
	}
	return ""
}

func (h *sheetHeader) String() string {
	s := ""
	for _, v := range h.lines {
		tmp := fmt.Sprintf("%v", v)
		tmp = strings.Replace(tmp, "\n", "\\n", -1)
		s += tmp + "\n"
	}
	return s
}

func newSheetHeader(st *xlsx.Sheet, excelName, sheetName string) (*sheetHeader, error) {
	// sheetHeaderRowsCount
	if len(st.Rows) < sheetHeaderRowsCount {
		return nil, fmt.Errorf("invalid sheetHeader rows count")
	}

	// 调整列数一致
	rowDesc, rowName, rowType, rowRequired, rowHome := st.Rows[0].Cells, st.Rows[1].Cells, st.Rows[2].Cells, st.Rows[3].Cells, st.Rows[4].Cells
	countDesc, countName, countType, countRequired, countHome := len(rowDesc), len(rowName), len(rowType), len(rowRequired), len(rowHome)
	linesLimit := countDesc
	if linesLimit > countName {
		linesLimit = countName
	}
	if linesLimit > countType {
		linesLimit = countType
	}
	if linesLimit > countRequired {
		linesLimit = countRequired
	}
	if linesLimit > countHome {
		linesLimit = countHome
	}
	rowDesc, rowName, rowType, rowRequired, rowHome = rowDesc[:linesLimit], rowName[:linesLimit], rowType[:linesLimit], rowRequired[:linesLimit], rowHome[:linesLimit]

	// 解析控制头
	lines := make([]*headerLine, linesLimit, linesLimit)
	for i := 0; i < linesLimit; i++ {
		// 取出内容
		lineDesc, err := rowDesc[i].FormattedValue()
		if err != nil {
			return nil, fmt.Errorf("sheetHeader get cell failed at row[0] line[%v]", i)
		}

		lineName, err := rowName[i].FormattedValue()
		if err != nil {
			return nil, fmt.Errorf("sheetHeader get cell failed at row[1] line[%v]", i)
		}

		lineType, err := rowType[i].FormattedValue()
		if err != nil {
			return nil, fmt.Errorf("sheetHeader get cell failed at row[2] line[%v]", i)
		}

		lineRequired, err := rowType[i].FormattedValue()
		if err != nil {
			return nil, fmt.Errorf("sheetHeader get cell failed at row[3] line[%v]", i)
		}

		lineHome, err := rowHome[i].FormattedValue()
		if err != nil {
			return nil, fmt.Errorf("sheetHeader get cell failed at row[4] line[%v]", i)
		}

		// 多列组合
		var headerLine *headerLine
		for j := 0; j < i; j++ {
			if lines[j].lineName == lineName {
				if !lines[j].lineType.IsMulti() {
					return nil, fmt.Errorf("sheetHeader lineType[%v] not support multi, but lineName[%v] appear at row[2] line[%v:%v]",
						lines[j].lineType.Type(), lineName, j, i)
				}
				if headerLine == nil {
					headerLine = lines[j]
				}
			}
		}

		// 创建新列
		if headerLine == nil {
			headerLine, err = newLineHeader(lineDesc, lineName, lineType, lineRequired, lineHome)
			if err != nil {
				return nil, fmt.Errorf("sheetHeader cell header invalid at row[2] line[%v], reason[%v]", i, err.Error())
			}
			headerLine.limitByExportConfig(excelName, sheetName)
		}

		// 记录
		lines[i] = headerLine
	}

	return &sheetHeader{
		lines: lines,
	}, nil
}
