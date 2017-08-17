package ffClientToml

import (
	"ffCommon/util"
	"ffLogic/ffGrammar"

	"fmt"

	"github.com/lexical005/toml"
)

// ExcelExportTest excel ExcelExportTest
type ExcelExportTest struct {
	VIPmap    map[int32]VIPmap
	VIPstruct VIPstruct
	VIPlist   []VIPlist
}

func (e *ExcelExportTest) String() string {
	result := ""
	result += "VIPmap"
	for k, v := range e.VIPmap {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	result += "VIPstruct"
	result += fmt.Sprintf("%v\n", e.VIPstruct)

	result += "VIPlist"
	for _, row := range e.VIPlist {
		result += fmt.Sprintf("%v\n", row)
	}

	return result
}

// Name the toml config's name
func (e *ExcelExportTest) Name() string {
	return "ExcelExportTest"
}

// VIPmap sheet VIPmap of excel ExcelExportTest
type VIPmap struct {
	InfoInt32        int32
	InfoInt64        int64
	InfoStr          string
	InfoInt32Array   []int32
	InfoInt64Array   []int64
	InfoStringArray  []string
	ItemClientID     int32
	Consume          ffGrammar.Grammar
	EmptyInt32       int32
	EmptyInt64       int64
	EmptyStr         string
	EmptyInt32Array  []int32
	EmptyInt64Array  []int64
	EmptyStringArray []string
}

func (vip *VIPmap) String() string {
	result := "["
	result += fmt.Sprintf("InfoInt32:%v,", vip.InfoInt32)
	result += fmt.Sprintf("InfoInt64:%v,", vip.InfoInt64)
	result += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	result += fmt.Sprintf("InfoInt32Array:%v,", vip.InfoInt32Array)
	result += fmt.Sprintf("InfoInt64Array:%v,", vip.InfoInt64Array)
	result += fmt.Sprintf("InfoStringArray:%v,", vip.InfoStringArray)
	result += fmt.Sprintf("ItemClientID:%v,", vip.ItemClientID)
	result += fmt.Sprintf("Consume:%v,", vip.Consume)
	result += fmt.Sprintf("EmptyInt32:%v,", vip.EmptyInt32)
	result += fmt.Sprintf("EmptyInt64:%v,", vip.EmptyInt64)
	result += fmt.Sprintf("EmptyStr:%v,", vip.EmptyStr)
	result += fmt.Sprintf("EmptyInt32Array:%v,", vip.EmptyInt32Array)
	result += fmt.Sprintf("EmptyInt64Array:%v,", vip.EmptyInt64Array)
	result += fmt.Sprintf("EmptyStringArray:%v,", vip.EmptyStringArray)
	result += "]"
	return result
}

// VIPstruct sheet VIPstruct of excel ExcelExportTest
type VIPstruct struct {
	InfoInt32        int32
	InfoInt64        int64
	InfoStr          string
	InfoInt32Array   []int32
	InfoInt64Array   []int64
	InfoStringArray  []string
	ItemClientID     int32
	Consume          ffGrammar.Grammar
	EmptyInt32       int32
	EmptyInt64       int64
	EmptyStr         string
	EmptyInt32Array  []int32
	EmptyInt64Array  []int64
	EmptyStringArray []string
}

func (vip *VIPstruct) String() string {
	result := "["
	result += fmt.Sprintf("InfoInt32:%v,", vip.InfoInt32)
	result += fmt.Sprintf("InfoInt64:%v,", vip.InfoInt64)
	result += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	result += fmt.Sprintf("InfoInt32Array:%v,", vip.InfoInt32Array)
	result += fmt.Sprintf("InfoInt64Array:%v,", vip.InfoInt64Array)
	result += fmt.Sprintf("InfoStringArray:%v,", vip.InfoStringArray)
	result += fmt.Sprintf("ItemClientID:%v,", vip.ItemClientID)
	result += fmt.Sprintf("Consume:%v,", vip.Consume)
	result += fmt.Sprintf("EmptyInt32:%v,", vip.EmptyInt32)
	result += fmt.Sprintf("EmptyInt64:%v,", vip.EmptyInt64)
	result += fmt.Sprintf("EmptyStr:%v,", vip.EmptyStr)
	result += fmt.Sprintf("EmptyInt32Array:%v,", vip.EmptyInt32Array)
	result += fmt.Sprintf("EmptyInt64Array:%v,", vip.EmptyInt64Array)
	result += fmt.Sprintf("EmptyStringArray:%v,", vip.EmptyStringArray)
	result += "]"
	return result
}

// VIPlist sheet VIPlist of excel ExcelExportTest
type VIPlist struct {
	InfoInt32        int32
	InfoInt64        int64
	InfoStr          string
	InfoInt32Array   []int32
	InfoInt64Array   []int64
	InfoStringArray  []string
	ItemClientID     int32
	Consume          ffGrammar.Grammar
	EmptyInt32       int32
	EmptyInt64       int64
	EmptyStr         string
	EmptyInt32Array  []int32
	EmptyInt64Array  []int64
	EmptyStringArray []string
}

func (vip *VIPlist) String() string {
	result := "["
	result += fmt.Sprintf("InfoInt32:%v,", vip.InfoInt32)
	result += fmt.Sprintf("InfoInt64:%v,", vip.InfoInt64)
	result += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	result += fmt.Sprintf("InfoInt32Array:%v,", vip.InfoInt32Array)
	result += fmt.Sprintf("InfoInt64Array:%v,", vip.InfoInt64Array)
	result += fmt.Sprintf("InfoStringArray:%v,", vip.InfoStringArray)
	result += fmt.Sprintf("ItemClientID:%v,", vip.ItemClientID)
	result += fmt.Sprintf("Consume:%v,", vip.Consume)
	result += fmt.Sprintf("EmptyInt32:%v,", vip.EmptyInt32)
	result += fmt.Sprintf("EmptyInt64:%v,", vip.EmptyInt64)
	result += fmt.Sprintf("EmptyStr:%v,", vip.EmptyStr)
	result += fmt.Sprintf("EmptyInt32Array:%v,", vip.EmptyInt32Array)
	result += fmt.Sprintf("EmptyInt64Array:%v,", vip.EmptyInt64Array)
	result += fmt.Sprintf("EmptyStringArray:%v,", vip.EmptyStringArray)
	result += "]"
	return result
}

// ReadExcelExportTest read excel ExcelExportTest
func ReadExcelExportTest() (e *ExcelExportTest, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/client/ExcelExportTest.toml")
	if err != nil {
		return
	}

	// 解析
	e = &ExcelExportTest{}
	err = toml.Unmarshal(fileContent, e)
	if err != nil {
		return
	}

	return
}
