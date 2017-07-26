package ffExcel

import (
	"cellvalue"
	"fmt"
	"sort"

	"github.com/tealeg/xlsx"
)

type contentRow struct {
	rowData map[string]cellvalue.ValueStore
}

func (cr *contentRow) String() string {
	keys := make([]string, 0, len(cr.rowData))
	for k := range cr.rowData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	s := ""
	for _, k := range keys {
		tmp := fmt.Sprintf("%v:%v", k, cr.rowData[k])
		if s != "" {
			s += "," + tmp
		} else {
			s += tmp
		}
	}
	return fmt.Sprintf("[%v]", s)
}

func newContentRow(rowIndex int, row *xlsx.Row, header *sheetHeader) (*contentRow, error) {
	rowData := make(map[string]cellvalue.ValueStore)
	for index, line := range header.lines {
		// 本列配置被忽略
		if line.ignore() {
			continue
		}

		// cell内容未配置
		if index >= len(row.Cells) {
			continue
		}

		// 取出cell内容
		cell := row.Cells[index]
		data, err := cell.String()
		if err != nil {
			return nil, err
		}

		// 只有允许多列配置时，才允许cell内容为空
		if data == "" && !line.lineType.IsMulti() && !line.lineType.IsString() {
			return nil, fmt.Errorf("only multi or string value type support cell content empty. row[%v] index[%v] type[%v]", rowIndex, index, line.lineType.Type())
		}

		// 存储值
		vs, ok := rowData[line.lineName]
		if !ok {
			vs, err = cellvalue.NewValueStore(line.lineType)
			if err != nil {
				return nil, err
			}
		}
		err = vs.Store(data)
		if err != nil {
			return nil, err
		}
		rowData[line.lineName] = vs
	}

	return &contentRow{
		rowData: rowData,
	}, nil
}
