package ffClientToml

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
	Name     string
	Desc     string
	SceneKey string
	Icon     string
	ItemType ffEnum.EItemType
	SubType  []ffEnum.EItemType
}

func (i *ItemTemplate) String() string {
	result := "["
	result += fmt.Sprintf("Name:%v,", i.Name)
	result += fmt.Sprintf("Desc:%v,", i.Desc)
	result += fmt.Sprintf("SceneKey:%v,", i.SceneKey)
	result += fmt.Sprintf("Icon:%v,", i.Icon)
	result += fmt.Sprintf("ItemType:%v,", i.ItemType)
	result += fmt.Sprintf("SubType:%v,", i.SubType)
	result += "]"
	return result
}

// ReadItem read excel Item
func ReadItem() (i *Item, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/client/Item.toml")
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
