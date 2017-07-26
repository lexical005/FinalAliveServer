package ffExcel

import (
	"fmt"
	"strings"

	"github.com/tealeg/xlsx"
)

type sheet struct {
	name      string
	sheetType int // sheetTypeList, sheetTypeMap
	header    *sheetHeader
	content   *sheetContent
}

func (s *sheet) String() string {
	headerDesc := "header[[字段描述][字段名称][字段类型][字段归属]]\n"
	contentDesc := "content"
	return fmt.Sprintf("sheet[%v]\n%v%v%v\n%v", s.name, headerDesc, s.header, contentDesc, s.content)
}

func newSheet(st *xlsx.Sheet) (*sheet, error) {
	sheetName := st.Name

	// readme
	if strings.ToLower(sheetName) == "readme" {
		return nil, errIgnoreSheetReadme
	}

	// sheet类型
	sheetType := sheetTypeInvalid
	if strings.HasSuffix(sheetName, sheetTypeListSuffix) {
		sheetType = sheetTypeList
		sheetName = sheetName[0 : len(sheetName)-len(sheetTypeListSuffix)]
	} else if strings.HasSuffix(sheetName, sheetTypeMapSuffix) {
		sheetType = sheetTypeMap
		sheetName = sheetName[0 : len(sheetName)-len(sheetTypeMapSuffix)]
	} else if strings.HasSuffix(sheetName, sheetTypeStructSuffix) {
		sheetType = sheetTypeStruct
		sheetName = sheetName[0 : len(sheetName)-len(sheetTypeStructSuffix)]
	}
	if sheetType == sheetTypeInvalid {
		return nil, errInvalidSheetName
	}

	// header
	header, err := newSheetHeader(st)
	if err != nil {
		return nil, err
	}

	// content
	content, err := newSheetContent(st, header)
	if err != nil {
		return nil, err
	}

	// check header and content
	if sheetType == sheetTypeList {
		if header.hasMapKey() {
			return nil, fmt.Errorf("sheet[%v] with suffix[%v] should not has [%v] line", sheetName, sheetTypeListSuffix, sheetTypeMapKeyName)
		}
	} else if sheetType == sheetTypeMap {
		if !header.hasMapKey() {
			return nil, fmt.Errorf("sheet[%v] with suffix[%v] must has [%v] line", sheetName, sheetTypeMapSuffix, sheetTypeMapKeyName)
		} else if _, ok := sheetTypeMapKeyType[header.mapKeyType()]; !ok {
			return nil, fmt.Errorf("sheet[%v] with suffix[%v] [%v] line type must in %v", sheetName, sheetTypeMapSuffix, sheetTypeMapKeyName, sheetTypeMapKeyType)
		}
	} else if sheetType == sheetTypeStruct {
		if header.hasMapKey() {
			return nil, fmt.Errorf("sheet[%v] with suffix[%v] should not has [%v] line", sheetName, sheetTypeStructSuffix, sheetTypeMapKeyName)
		} else if len(content.rows) != 1 {
			return nil, fmt.Errorf("sheet[%v] with suffix[%v] must 1 row", sheetName, sheetTypeStructSuffix)
		}
	}

	return &sheet{
		name:      sheetName,
		sheetType: sheetType,
		header:    header,
		content:   content,
	}, nil
}
