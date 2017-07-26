package ffExcel

import (
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
)

type sheetContent struct {
	rows []*contentRow
}

func (sc *sheetContent) String() string {
	s := ""
	for _, v := range sc.rows {
		tmp := fmt.Sprintf("%v", v)
		tmp = strings.Replace(tmp, "\n", "\\n", -1)
		if s != "" {
			s += "\n"
		}
		s += tmp
	}
	return s
}

func newSheetContent(st *xlsx.Sheet, header *sheetHeader) (*sheetContent, error) {
	rows := make([]*contentRow, 0, len(st.Rows)-sheetHeaderRowsCount)
	for i := sheetHeaderRowsCount; i < len(st.Rows); i++ {
		// 本行忽略
		if len(st.Rows[i].Cells) == 0 {
			continue
		}

		// 关键字列为空时忽略
		key, err := st.Rows[i].Cells[0].String()
		if err != nil {
			return nil, err
		}

		if len(key) == 0 {
			continue
		}

		// 生成一行数据
		row, err := newContentRow(i, st.Rows[i], header)
		if err != nil {
			return nil, err
		}

		rows = append(rows, row)
	}

	return &sheetContent{
		rows: rows,
	}, nil
}
