package ffGameConfig

import (
	"ffCommon/util"

	"ffAutoGen/ffEnum"

	"fmt"

	"github.com/lexical005/toml"
)

// RandBorn excel RandBorn
type RandBorn struct {
	BornPosition    []*BornPosition
	BornPrepareItem []*BornPrepareItem
	BornBattleItem  []*BornBattleItem
	ItemBase        map[int32]*ItemBase
}

func (r *RandBorn) String() string {
	result := ""
	result += "BornPosition"
	for _, row := range r.BornPosition {
		result += fmt.Sprintf("%v\n", row)
	}

	result += "BornPrepareItem"
	for _, row := range r.BornPrepareItem {
		result += fmt.Sprintf("%v\n", row)
	}

	result += "BornBattleItem"
	for _, row := range r.BornBattleItem {
		result += fmt.Sprintf("%v\n", row)
	}

	result += "ItemBase"
	for k, v := range r.ItemBase {
		result += fmt.Sprintf("%v:%v\n", k, v)
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

// BornPrepareItem sheet BornPrepareItem of excel RandBorn
type BornPrepareItem struct {
	Area     int32
	Group    int32
	ItemBase int32
}

func (b *BornPrepareItem) String() string {
	result := "["
	result += fmt.Sprintf("Area:%v,", b.Area)
	result += fmt.Sprintf("Group:%v,", b.Group)
	result += fmt.Sprintf("ItemBase:%v,", b.ItemBase)
	result += "]"
	return result
}

// BornBattleItem sheet BornBattleItem of excel RandBorn
type BornBattleItem struct {
	Area     int32
	ItemBase int32
	MinCount int32
	MaxCount int32
}

func (b *BornBattleItem) String() string {
	result := "["
	result += fmt.Sprintf("Area:%v,", b.Area)
	result += fmt.Sprintf("ItemBase:%v,", b.ItemBase)
	result += fmt.Sprintf("MinCount:%v,", b.MinCount)
	result += fmt.Sprintf("MaxCount:%v,", b.MaxCount)
	result += "]"
	return result
}

// ItemBase sheet ItemBase of excel RandBorn
type ItemBase struct {
	Chances []int32
	Items   []int32
	Numbers []int32
}

func (i *ItemBase) String() string {
	result := "["
	result += fmt.Sprintf("Chances:%v,", i.Chances)
	result += fmt.Sprintf("Items:%v,", i.Items)
	result += fmt.Sprintf("Numbers:%v,", i.Numbers)
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
