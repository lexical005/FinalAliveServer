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
	InfoInt32Single  []int32
	InfoInt64Single  []int64
	InfoStrSingle    []string
	InfoInt32Multi   []int32
	InfoInt64Multi   []int64
	InfoStrMulti     []string
	ItemClientID     int32
	Consume          ffGrammar.Grammar
	EmptyInt32       int32
	EmptyInt64       int64
	EmptyStr         string
	EmptyIn32tSingle []int32
	EmptyInt64Single []int64
	EmptyStrSingle   []string
	EmptyInt32Multi  []int32
	EmptyInt64Multi  []int64
	EmptyStrMulti    []string
}

func (vip *VIPmap) String() string {
	result := "["
	result += fmt.Sprintf("InfoInt32:%v,", vip.InfoInt32)
	result += fmt.Sprintf("InfoInt64:%v,", vip.InfoInt64)
	result += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	result += fmt.Sprintf("InfoInt32Single:%v,", vip.InfoInt32Single)
	result += fmt.Sprintf("InfoInt64Single:%v,", vip.InfoInt64Single)
	result += fmt.Sprintf("InfoStrSingle:%v,", vip.InfoStrSingle)
	result += fmt.Sprintf("InfoInt32Multi:%v,", vip.InfoInt32Multi)
	result += fmt.Sprintf("InfoInt64Multi:%v,", vip.InfoInt64Multi)
	result += fmt.Sprintf("InfoStrMulti:%v,", vip.InfoStrMulti)
	result += fmt.Sprintf("ItemClientID:%v,", vip.ItemClientID)
	result += fmt.Sprintf("Consume:%v,", vip.Consume)
	result += fmt.Sprintf("EmptyInt32:%v,", vip.EmptyInt32)
	result += fmt.Sprintf("EmptyInt64:%v,", vip.EmptyInt64)
	result += fmt.Sprintf("EmptyStr:%v,", vip.EmptyStr)
	result += fmt.Sprintf("EmptyIn32tSingle:%v,", vip.EmptyIn32tSingle)
	result += fmt.Sprintf("EmptyInt64Single:%v,", vip.EmptyInt64Single)
	result += fmt.Sprintf("EmptyStrSingle:%v,", vip.EmptyStrSingle)
	result += fmt.Sprintf("EmptyInt32Multi:%v,", vip.EmptyInt32Multi)
	result += fmt.Sprintf("EmptyInt64Multi:%v,", vip.EmptyInt64Multi)
	result += fmt.Sprintf("EmptyStrMulti:%v,", vip.EmptyStrMulti)
	result += "]"
	return result
}

// VIPstruct sheet VIPstruct of excel ExcelExportTest
type VIPstruct struct {
	InfoInt32        int32
	InfoInt64        int64
	InfoStr          string
	InfoInt32Single  []int32
	InfoInt64Single  []int64
	InfoStrSingle    []string
	InfoInt32Multi   []int32
	InfoInt64Multi   []int64
	InfoStrMulti     []string
	ItemClientID     int32
	Consume          ffGrammar.Grammar
	EmptyInt32       int32
	EmptyInt64       int64
	EmptyStr         string
	EmptyIn32tSingle []int32
	EmptyInt64Single []int64
	EmptyStrSingle   []string
	EmptyInt32Multi  []int32
	EmptyInt64Multi  []int64
	EmptyStrMulti    []string
}

func (vip *VIPstruct) String() string {
	result := "["
	result += fmt.Sprintf("InfoInt32:%v,", vip.InfoInt32)
	result += fmt.Sprintf("InfoInt64:%v,", vip.InfoInt64)
	result += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	result += fmt.Sprintf("InfoInt32Single:%v,", vip.InfoInt32Single)
	result += fmt.Sprintf("InfoInt64Single:%v,", vip.InfoInt64Single)
	result += fmt.Sprintf("InfoStrSingle:%v,", vip.InfoStrSingle)
	result += fmt.Sprintf("InfoInt32Multi:%v,", vip.InfoInt32Multi)
	result += fmt.Sprintf("InfoInt64Multi:%v,", vip.InfoInt64Multi)
	result += fmt.Sprintf("InfoStrMulti:%v,", vip.InfoStrMulti)
	result += fmt.Sprintf("ItemClientID:%v,", vip.ItemClientID)
	result += fmt.Sprintf("Consume:%v,", vip.Consume)
	result += fmt.Sprintf("EmptyInt32:%v,", vip.EmptyInt32)
	result += fmt.Sprintf("EmptyInt64:%v,", vip.EmptyInt64)
	result += fmt.Sprintf("EmptyStr:%v,", vip.EmptyStr)
	result += fmt.Sprintf("EmptyIn32tSingle:%v,", vip.EmptyIn32tSingle)
	result += fmt.Sprintf("EmptyInt64Single:%v,", vip.EmptyInt64Single)
	result += fmt.Sprintf("EmptyStrSingle:%v,", vip.EmptyStrSingle)
	result += fmt.Sprintf("EmptyInt32Multi:%v,", vip.EmptyInt32Multi)
	result += fmt.Sprintf("EmptyInt64Multi:%v,", vip.EmptyInt64Multi)
	result += fmt.Sprintf("EmptyStrMulti:%v,", vip.EmptyStrMulti)
	result += "]"
	return result
}

// VIPlist sheet VIPlist of excel ExcelExportTest
type VIPlist struct {
	InfoInt32        int32
	InfoInt64        int64
	InfoStr          string
	InfoInt32Single  []int32
	InfoInt64Single  []int64
	InfoStrSingle    []string
	InfoInt32Multi   []int32
	InfoInt64Multi   []int64
	InfoStrMulti     []string
	ItemClientID     int32
	Consume          ffGrammar.Grammar
	EmptyInt32       int32
	EmptyInt64       int64
	EmptyStr         string
	EmptyIn32tSingle []int32
	EmptyInt64Single []int64
	EmptyStrSingle   []string
	EmptyInt32Multi  []int32
	EmptyInt64Multi  []int64
	EmptyStrMulti    []string
}

func (vip *VIPlist) String() string {
	result := "["
	result += fmt.Sprintf("InfoInt32:%v,", vip.InfoInt32)
	result += fmt.Sprintf("InfoInt64:%v,", vip.InfoInt64)
	result += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	result += fmt.Sprintf("InfoInt32Single:%v,", vip.InfoInt32Single)
	result += fmt.Sprintf("InfoInt64Single:%v,", vip.InfoInt64Single)
	result += fmt.Sprintf("InfoStrSingle:%v,", vip.InfoStrSingle)
	result += fmt.Sprintf("InfoInt32Multi:%v,", vip.InfoInt32Multi)
	result += fmt.Sprintf("InfoInt64Multi:%v,", vip.InfoInt64Multi)
	result += fmt.Sprintf("InfoStrMulti:%v,", vip.InfoStrMulti)
	result += fmt.Sprintf("ItemClientID:%v,", vip.ItemClientID)
	result += fmt.Sprintf("Consume:%v,", vip.Consume)
	result += fmt.Sprintf("EmptyInt32:%v,", vip.EmptyInt32)
	result += fmt.Sprintf("EmptyInt64:%v,", vip.EmptyInt64)
	result += fmt.Sprintf("EmptyStr:%v,", vip.EmptyStr)
	result += fmt.Sprintf("EmptyIn32tSingle:%v,", vip.EmptyIn32tSingle)
	result += fmt.Sprintf("EmptyInt64Single:%v,", vip.EmptyInt64Single)
	result += fmt.Sprintf("EmptyStrSingle:%v,", vip.EmptyStrSingle)
	result += fmt.Sprintf("EmptyInt32Multi:%v,", vip.EmptyInt32Multi)
	result += fmt.Sprintf("EmptyInt64Multi:%v,", vip.EmptyInt64Multi)
	result += fmt.Sprintf("EmptyStrMulti:%v,", vip.EmptyStrMulti)
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
