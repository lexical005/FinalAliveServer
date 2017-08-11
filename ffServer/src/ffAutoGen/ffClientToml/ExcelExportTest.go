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
	s := ""
	s += "VIPmap"
	for k, v := range e.VIPmap {
		s += fmt.Sprintf("%v:%v\n", k, v)
	}

	s += "VIPstruct"
	s += fmt.Sprintf("%v\n", e.VIPstruct)

	s += "VIPlist"
	for _, row := range e.VIPlist {
		s += fmt.Sprintf("%v\n", row)
	}

	return s
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
	s := "["
	s += fmt.Sprintf("InfoInt32:%v,", vip.InfoInt32)
	s += fmt.Sprintf("InfoInt64:%v,", vip.InfoInt64)
	s += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	s += fmt.Sprintf("InfoInt32Single:%v,", vip.InfoInt32Single)
	s += fmt.Sprintf("InfoInt64Single:%v,", vip.InfoInt64Single)
	s += fmt.Sprintf("InfoStrSingle:%v,", vip.InfoStrSingle)
	s += fmt.Sprintf("InfoInt32Multi:%v,", vip.InfoInt32Multi)
	s += fmt.Sprintf("InfoInt64Multi:%v,", vip.InfoInt64Multi)
	s += fmt.Sprintf("InfoStrMulti:%v,", vip.InfoStrMulti)
	s += fmt.Sprintf("ItemClientID:%v,", vip.ItemClientID)
	s += fmt.Sprintf("Consume:%v,", vip.Consume)
	s += fmt.Sprintf("EmptyInt32:%v,", vip.EmptyInt32)
	s += fmt.Sprintf("EmptyInt64:%v,", vip.EmptyInt64)
	s += fmt.Sprintf("EmptyStr:%v,", vip.EmptyStr)
	s += fmt.Sprintf("EmptyIn32tSingle:%v,", vip.EmptyIn32tSingle)
	s += fmt.Sprintf("EmptyInt64Single:%v,", vip.EmptyInt64Single)
	s += fmt.Sprintf("EmptyStrSingle:%v,", vip.EmptyStrSingle)
	s += fmt.Sprintf("EmptyInt32Multi:%v,", vip.EmptyInt32Multi)
	s += fmt.Sprintf("EmptyInt64Multi:%v,", vip.EmptyInt64Multi)
	s += fmt.Sprintf("EmptyStrMulti:%v,", vip.EmptyStrMulti)
	s += "]"
	return s
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
	s := "["
	s += fmt.Sprintf("InfoInt32:%v,", vip.InfoInt32)
	s += fmt.Sprintf("InfoInt64:%v,", vip.InfoInt64)
	s += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	s += fmt.Sprintf("InfoInt32Single:%v,", vip.InfoInt32Single)
	s += fmt.Sprintf("InfoInt64Single:%v,", vip.InfoInt64Single)
	s += fmt.Sprintf("InfoStrSingle:%v,", vip.InfoStrSingle)
	s += fmt.Sprintf("InfoInt32Multi:%v,", vip.InfoInt32Multi)
	s += fmt.Sprintf("InfoInt64Multi:%v,", vip.InfoInt64Multi)
	s += fmt.Sprintf("InfoStrMulti:%v,", vip.InfoStrMulti)
	s += fmt.Sprintf("ItemClientID:%v,", vip.ItemClientID)
	s += fmt.Sprintf("Consume:%v,", vip.Consume)
	s += fmt.Sprintf("EmptyInt32:%v,", vip.EmptyInt32)
	s += fmt.Sprintf("EmptyInt64:%v,", vip.EmptyInt64)
	s += fmt.Sprintf("EmptyStr:%v,", vip.EmptyStr)
	s += fmt.Sprintf("EmptyIn32tSingle:%v,", vip.EmptyIn32tSingle)
	s += fmt.Sprintf("EmptyInt64Single:%v,", vip.EmptyInt64Single)
	s += fmt.Sprintf("EmptyStrSingle:%v,", vip.EmptyStrSingle)
	s += fmt.Sprintf("EmptyInt32Multi:%v,", vip.EmptyInt32Multi)
	s += fmt.Sprintf("EmptyInt64Multi:%v,", vip.EmptyInt64Multi)
	s += fmt.Sprintf("EmptyStrMulti:%v,", vip.EmptyStrMulti)
	s += "]"
	return s
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
	s := "["
	s += fmt.Sprintf("InfoInt32:%v,", vip.InfoInt32)
	s += fmt.Sprintf("InfoInt64:%v,", vip.InfoInt64)
	s += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	s += fmt.Sprintf("InfoInt32Single:%v,", vip.InfoInt32Single)
	s += fmt.Sprintf("InfoInt64Single:%v,", vip.InfoInt64Single)
	s += fmt.Sprintf("InfoStrSingle:%v,", vip.InfoStrSingle)
	s += fmt.Sprintf("InfoInt32Multi:%v,", vip.InfoInt32Multi)
	s += fmt.Sprintf("InfoInt64Multi:%v,", vip.InfoInt64Multi)
	s += fmt.Sprintf("InfoStrMulti:%v,", vip.InfoStrMulti)
	s += fmt.Sprintf("ItemClientID:%v,", vip.ItemClientID)
	s += fmt.Sprintf("Consume:%v,", vip.Consume)
	s += fmt.Sprintf("EmptyInt32:%v,", vip.EmptyInt32)
	s += fmt.Sprintf("EmptyInt64:%v,", vip.EmptyInt64)
	s += fmt.Sprintf("EmptyStr:%v,", vip.EmptyStr)
	s += fmt.Sprintf("EmptyIn32tSingle:%v,", vip.EmptyIn32tSingle)
	s += fmt.Sprintf("EmptyInt64Single:%v,", vip.EmptyInt64Single)
	s += fmt.Sprintf("EmptyStrSingle:%v,", vip.EmptyStrSingle)
	s += fmt.Sprintf("EmptyInt32Multi:%v,", vip.EmptyInt32Multi)
	s += fmt.Sprintf("EmptyInt64Multi:%v,", vip.EmptyInt64Multi)
	s += fmt.Sprintf("EmptyStrMulti:%v,", vip.EmptyStrMulti)
	s += "]"
	return s
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
