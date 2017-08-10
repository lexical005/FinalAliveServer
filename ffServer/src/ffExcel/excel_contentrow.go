package ffExcel

import (
	"cellvalue"
	"fmt"
	"sort"

	"github.com/lexical005/xlsx"
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
		data, err := cell.FormattedValue()
		if err != nil {
			return nil, err
		}

		// 如果是空字符串
		if data == "" {
			// 本列必须配置有效值, 不允许留空白
			if line.isRequired() {
				if !line.lineType.IsMulti() {
					return nil, fmt.Errorf("not allow empty cell at row[%v] line[%v] lineType[%v]",
						rowIndex, index, line.lineType.Type())
				}

				// 留待整行全部解析完成后, 再进行检查(只要有一列配置了值, 列就是有效的)
				continue
			}

			// 可选配置列, 先返回
			continue
		} else if line.lineType.IsString() {
			// 本列配置的是字符串, 且使用者主动配置了空字符串, 则将其转换为程序用的空字符串
			if data == `""` {
				data = ""
			}
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
