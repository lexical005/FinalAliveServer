package ffExcel

import "fmt"

func checkToml(excel *excel) (err error) {
	// check excel name, sheet name, field name
	shortExcelName := getShortName(excel.name)
	if shortExcelName == "" {
		err = fmt.Errorf("excel name[%v] must start with upper char", excel.name)
	}

	for _, sheet := range excel.sheets {
		shortSheetName := getShortName(sheet.name)
		if shortSheetName == "" {
			if err != nil {
				err = fmt.Errorf("%v\nsheet name[%v] must start with upper char in excel[%v]", err, sheet.name, excel.name)
			} else {
				err = fmt.Errorf("sheet name[%v] must start with upper char in excel[%v]", sheet.name, excel.name)
			}
		}

		for _, line := range sheet.header.lines {
			if !line.ignore() && !line.isMapKey() {
				shortFieldName := getShortName(line.lineName)
				if shortFieldName == "" {
					if err != nil {
						err = fmt.Errorf("%v\nfield name[%v] must start with upper char in excel[%v] sheet[%v]", err, line.lineName, excel.name, sheet.name)
					} else {
						err = fmt.Errorf("field name[%v] must start with upper char in excel[%v] sheet[%v]", line.lineName, excel.name, sheet.name)
					}
				}
			}
		}
	}
	return
}

func genToml(excel *excel, exportConfig *ExportConfig) (tomlDataServerReader string, tomlData string) {
	tomlDataServerReader = genTomlDataReadCode(excel, exportConfig, "server")
	tomlData = genTomlData(excel, exportConfig, "server")
	return
}
