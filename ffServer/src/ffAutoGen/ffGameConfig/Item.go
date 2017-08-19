package ffGameConfig

import (
	"ffCommon/util"

	"ffAutoGen/ffEnum"

	"fmt"

	"github.com/lexical005/toml"
)

// Item excel Item
type Item struct {
	ItemTemplate map[int32]ItemTemplate
}

func (i *Item) String() string {
	result := ""
	result += "ItemTemplate"
	for k, v := range i.ItemTemplate {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	return result
}

// Name the toml config's name
func (i *Item) Name() string {
	return "Item"
}

// ItemTemplate sheet ItemTemplate of excel Item
type ItemTemplate struct {
	ItemType ffEnum.EItemType
}

func (i *ItemTemplate) String() string {
	result := "["
	result += fmt.Sprintf("ItemType:%v,", i.ItemType)
	result += "]"
	return result
}

// ReadItem read excel Item
func ReadItem() (i *Item, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/game/Item.toml")
	if err != nil {
		return
	}

	// 解析
	i = &Item{}
	err = toml.Unmarshal(fileContent, i)
	if err != nil {
		return
	}

	return
}
