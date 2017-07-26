package ffExcel

import "fmt"

const (
	sheetHeaderRowsCount = 4

	sheetTypeInvalid = iota
	sheetTypeList
	sheetTypeMap
	sheetTypeStruct

	sheetTypeListSuffix   = "_list"
	sheetTypeMapSuffix    = "_map"
	sheetTypeStructSuffix = "_struct"

	sheetTypeMapKeyName = "Key"
)

var sheetTypeMapKeyType = map[string]bool{
	"int":    true,
	"string": true,
}

var errIgnoreSheetReadme = fmt.Errorf("ignore sheet readme")
var errInvalidSheetName = fmt.Errorf("sheet name must has suffix in list[%v:%v]", sheetTypeListSuffix, sheetTypeMapSuffix)
