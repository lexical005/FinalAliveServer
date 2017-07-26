package ffGameConfig

import (
	"ffCommon/util"
	"ffLogic/ffGrammar"

	"fmt"

	"github.com/lexical005/toml"
)

// Mall excel Mall
type Mall struct {
	VIPmap    map[int]VIPmap
	VIPstruct VIPstruct
	VIPlist   []VIPlist
}

func (m Mall) String() string {
	s := ""
	s += "VIPmap"
	for k, v := range m.VIPmap {
		s += fmt.Sprintf("%v:%v\n", k, v)
	}

	s += "VIPstruct"
	s += fmt.Sprintf("%v\n", m.VIPstruct)

	s += "VIPlist"
	for _, row := range m.VIPlist {
		s += fmt.Sprintf("%v\n", row)
	}

	return s
}

// VIPmap sheet VIPmap of excel Mall
type VIPmap struct {
	InfoInt       int
	InfoStr       string
	InfoIntSingle []int
	InfoStrSingle []string
	InfoIntMulti  []int
	InfoStrMulti  []string
	Award         ffGrammar.Grammar
}

func (vip VIPmap) String() string {
	s := "["
	s += fmt.Sprintf("InfoInt:%v,", vip.InfoInt)
	s += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	s += fmt.Sprintf("InfoIntSingle:%v,", vip.InfoIntSingle)
	s += fmt.Sprintf("InfoStrSingle:%v,", vip.InfoStrSingle)
	s += fmt.Sprintf("InfoIntMulti:%v,", vip.InfoIntMulti)
	s += fmt.Sprintf("InfoStrMulti:%v,", vip.InfoStrMulti)
	s += fmt.Sprintf("Award:%v,", vip.Award)
	s += "]"
	return s
}

// VIPstruct sheet VIPstruct of excel Mall
type VIPstruct struct {
	InfoInt       int
	InfoStr       string
	InfoIntSingle []int
	InfoStrSingle []string
	InfoIntMulti  []int
	InfoStrMulti  []string
}

func (vip VIPstruct) String() string {
	s := "["
	s += fmt.Sprintf("InfoInt:%v,", vip.InfoInt)
	s += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	s += fmt.Sprintf("InfoIntSingle:%v,", vip.InfoIntSingle)
	s += fmt.Sprintf("InfoStrSingle:%v,", vip.InfoStrSingle)
	s += fmt.Sprintf("InfoIntMulti:%v,", vip.InfoIntMulti)
	s += fmt.Sprintf("InfoStrMulti:%v,", vip.InfoStrMulti)
	s += "]"
	return s
}

// VIPlist sheet VIPlist of excel Mall
type VIPlist struct {
	InfoInt       int
	InfoStr       string
	InfoIntSingle []int
	InfoStrSingle []string
	InfoIntMulti  []int
	InfoStrMulti  []string
}

func (vip VIPlist) String() string {
	s := "["
	s += fmt.Sprintf("InfoInt:%v,", vip.InfoInt)
	s += fmt.Sprintf("InfoStr:%v,", vip.InfoStr)
	s += fmt.Sprintf("InfoIntSingle:%v,", vip.InfoIntSingle)
	s += fmt.Sprintf("InfoStrSingle:%v,", vip.InfoStrSingle)
	s += fmt.Sprintf("InfoIntMulti:%v,", vip.InfoIntMulti)
	s += fmt.Sprintf("InfoStrMulti:%v,", vip.InfoStrMulti)
	s += "]"
	return s
}

// ReadMall read excel Mall
func ReadMall() (m *Mall, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/game/Mall.toml")
	if err != nil {
		return
	}

	// 解析
	m = &Mall{}
	err = toml.Unmarshal(fileContent, m)
	if err != nil {
		return
	}

	return
}
