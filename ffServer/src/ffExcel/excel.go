package ffExcel

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/tealeg/xlsx"
)

type excel struct {
	name   string
	sheets []*sheet
}

func (e *excel) String() string {
	s := ""
	for _, sheet := range e.sheets {
		s += fmt.Sprintf("\n%v\n", sheet)
	}
	return fmt.Sprintf("excel[%v] sheets:%v", e.name, s)
}

func parseExcel(excelFilePath string) (*excel, error) {
	file, e := xlsx.OpenFile(excelFilePath)
	if e != nil {
		return nil, e
	}

	sheets := make([]*sheet, 0, len(file.Sheets))

	var errResult error
	for _, st := range file.Sheets {
		sheet, err := newSheet(st)
		if err != nil {
			if err != errIgnoreSheetReadme {
				e := fmt.Errorf("excel[%v] sheet[%v] get error[%v]\n", excelFilePath, st.Name, err.Error())
				if errResult == nil {
					errResult = e
				} else {
					errResult = fmt.Errorf("%v%v", err.Error(), e)
				}
			}
			continue
		}
		sheets = append(sheets, sheet)
	}

	_, fileNameWithExt := filepath.Split(excelFilePath)
	fileExt := path.Ext(fileNameWithExt)
	fileName := fileNameWithExt[0 : len(fileNameWithExt)-len(fileExt)]

	return &excel{
		name:   fileName,
		sheets: sheets,
	}, errResult
}
