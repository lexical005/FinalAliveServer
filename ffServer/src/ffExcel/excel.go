package ffExcel

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/lexical005/xlsx"
)

type excel struct {
	name       string
	sheets     []*sheet
	exportType string
}

// exportToServer 本工作表是否需要导出到服务端
func (e *excel) exportToServer() bool {
	for _, sheet := range e.sheets {
		if sheet.exportToServer() {
			return true
		}
	}
	return false
}

// exportToClient 本工作表是否需要导出到客户端
func (e *excel) exportToClient() bool {
	for _, sheet := range e.sheets {
		if sheet.exportToClient() {
			return true
		}
	}
	return false
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

	_, fileNameWithExt := filepath.Split(excelFilePath)
	fileExt := path.Ext(fileNameWithExt)
	fileName := fileNameWithExt[0 : len(fileNameWithExt)-len(fileExt)]

	sheets := make([]*sheet, 0, len(file.Sheets))

	var errResult error
	for _, st := range file.Sheets {
		sheet, err := newSheet(st, fileName)
		if err != nil {
			if err != errIgnoreSheetReadme {
				e := fmt.Errorf("excel[%v] sheet[%v] get error[%v]", excelFilePath, st.Name, err.Error())
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

	return &excel{
		name:       fileName,
		sheets:     sheets,
		exportType: "config",
	}, errResult
}
