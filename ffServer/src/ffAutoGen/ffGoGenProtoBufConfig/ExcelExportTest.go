package ffGameConfig


import (
	"ffCommon/util"
	"ffLogic/ffGrammar"

	"fmt"

	"github.com/lexical005/toml"
)

// ExcelExportTest excel ExcelExportTest
type ExcelExportTest struct {
   VIPmap map[int]VIPmap
   VIPstruct VIPstruct
   VIPlist []VIPlist
}

func (e ExcelExportTest) String() string {
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

// VIPmap sheet VIPmap of excel ExcelExportTest
type VIPmap struct {
	InfoInt  int
	InfoStr  string
	ItemID  int
	Award  ffGrammar.Grammar
}

func (vip VIPmap) String() string {
	s := "["
	s += fmt.Sprintf("InfoInt:%v,", vip.InfoInt)
	s += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	s += fmt.Sprintf("ItemID:%v,", vip.ItemID)
	s += fmt.Sprintf("Award:%v,", vip.Award)
	s += "]"
	return s
}

// VIPstruct sheet VIPstruct of excel ExcelExportTest
type VIPstruct struct {
	InfoInt  int
	InfoStr  string
	ItemID  int
}

func (vip VIPstruct) String() string {
	s := "["
	s += fmt.Sprintf("InfoInt:%v,", vip.InfoInt)
	s += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	s += fmt.Sprintf("ItemID:%v,", vip.ItemID)
	s += "]"
	return s
}

// VIPlist sheet VIPlist of excel ExcelExportTest
type VIPlist struct {
	InfoInt  int
	InfoStr  string
	ItemID  int
}

func (vip VIPlist) String() string {
	s := "["
	s += fmt.Sprintf("InfoInt:%v,", vip.InfoInt)
	s += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	s += fmt.Sprintf("ItemID:%v,", vip.ItemID)
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
