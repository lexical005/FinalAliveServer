package ffGameConfig

import (
	"ffCommon/util"

	"ffAutoGen/ffEnum"

	"fmt"

	"github.com/lexical005/toml"
)

// RandBorn excel RandBorn
type RandBorn struct {
	BornPosition []*BornPosition
}

func (r *RandBorn) String() string {
	result := ""
	result += "BornPosition"
	for _, row := range r.BornPosition {
		result += fmt.Sprintf("%v\n", row)
	}

	return result
}

// Name the toml config's name
func (r *RandBorn) Name() string {
	return "RandBorn"
}

// BornPosition sheet BornPosition of excel RandBorn
type BornPosition struct {
	Type      ffEnum.EBornType
	Area      int32
	Group     int32
	Positions []int32
}

func (b *BornPosition) String() string {
	result := "["
	result += fmt.Sprintf("Type:%v,", b.Type)
	result += fmt.Sprintf("Area:%v,", b.Area)
	result += fmt.Sprintf("Group:%v,", b.Group)
	result += fmt.Sprintf("Positions:%v,", b.Positions)
	result += "]"
	return result
}

// ReadRandBorn read excel RandBorn
func ReadRandBorn() (r *RandBorn, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/game/RandBorn.toml")
	if err != nil {
		return
	}

	// 解析
	r = &RandBorn{}
	err = toml.Unmarshal(fileContent, r)
	if err != nil {
		return
	}

	return
}
