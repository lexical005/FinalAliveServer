package ffGameConfig

import (
	"ffCommon/util"

	"ffAutoGen/ffEnum"

	"fmt"

	"github.com/lexical005/toml"
)

// Item excel Item
type Item struct {
	ItemTemplate map[int32]*ItemTemplate
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
	ItemType    ffEnum.EItemType
	SubType     []ffEnum.EItemType
	Attrs1      map[int32]int32
	Attrs1Key   []int32
	Attrs1Value []int32
	Attrs2      map[ffEnum.EItemType]int32
	Attrs2Key   []ffEnum.EItemType
	Attrs2Value []int32
	Attrs3      map[ffEnum.EItemType]int32
	Attrs3Key   []ffEnum.EItemType
	Attrs3Value []int32
	Attrs4      map[ffEnum.EItemType]ffEnum.EShootMode
	Attrs4Key   []ffEnum.EItemType
	Attrs4Value []ffEnum.EShootMode
	Attrs5      map[ffEnum.EItemType]ffEnum.EShootMode
	Attrs5Key   []ffEnum.EItemType
	Attrs5Value []ffEnum.EShootMode
	Attrs6      map[string]string
	Attrs6Key   []string
	Attrs6Value []string
}

func (i *ItemTemplate) String() string {
	result := "["
	result += fmt.Sprintf("ItemType:%v,", i.ItemType)
	result += fmt.Sprintf("SubType:%v,", i.SubType)
	result += fmt.Sprintf("Attrs1:%v,", i.Attrs1)
	result += fmt.Sprintf("Attrs1Key:%v,", i.Attrs1Key)
	result += fmt.Sprintf("Attrs1Value:%v,", i.Attrs1Value)
	result += fmt.Sprintf("Attrs2:%v,", i.Attrs2)
	result += fmt.Sprintf("Attrs2Key:%v,", i.Attrs2Key)
	result += fmt.Sprintf("Attrs2Value:%v,", i.Attrs2Value)
	result += fmt.Sprintf("Attrs3:%v,", i.Attrs3)
	result += fmt.Sprintf("Attrs3Key:%v,", i.Attrs3Key)
	result += fmt.Sprintf("Attrs3Value:%v,", i.Attrs3Value)
	result += fmt.Sprintf("Attrs4:%v,", i.Attrs4)
	result += fmt.Sprintf("Attrs4Key:%v,", i.Attrs4Key)
	result += fmt.Sprintf("Attrs4Value:%v,", i.Attrs4Value)
	result += fmt.Sprintf("Attrs5:%v,", i.Attrs5)
	result += fmt.Sprintf("Attrs5Key:%v,", i.Attrs5Key)
	result += fmt.Sprintf("Attrs5Value:%v,", i.Attrs5Value)
	result += fmt.Sprintf("Attrs6:%v,", i.Attrs6)
	result += fmt.Sprintf("Attrs6Key:%v,", i.Attrs6Key)
	result += fmt.Sprintf("Attrs6Value:%v,", i.Attrs6Value)
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

	for _, one := range i.ItemTemplate {
		one.Attrs1 = make(map[int32]int32, len(one.Attrs1Key))
		for index, v := range one.Attrs1Key {
			one.Attrs1[v] = one.Attrs1Value[index]
		}
	}

	for _, one := range i.ItemTemplate {
		one.Attrs2 = make(map[ffEnum.EItemType]int32, len(one.Attrs2Key))
		for index, v := range one.Attrs2Key {
			one.Attrs2[v] = one.Attrs2Value[index]
		}
	}

	for _, one := range i.ItemTemplate {
		one.Attrs3 = make(map[ffEnum.EItemType]int32, len(one.Attrs3Key))
		for index, v := range one.Attrs3Key {
			one.Attrs3[v] = one.Attrs3Value[index]
		}
	}

	for _, one := range i.ItemTemplate {
		one.Attrs4 = make(map[ffEnum.EItemType]ffEnum.EShootMode, len(one.Attrs4Key))
		for index, v := range one.Attrs4Key {
			one.Attrs4[v] = one.Attrs4Value[index]
		}
	}

	for _, one := range i.ItemTemplate {
		one.Attrs5 = make(map[ffEnum.EItemType]ffEnum.EShootMode, len(one.Attrs5Key))
		for index, v := range one.Attrs5Key {
			one.Attrs5[v] = one.Attrs5Value[index]
		}
	}

	for _, one := range i.ItemTemplate {
		one.Attrs6 = make(map[string]string, len(one.Attrs6Key))
		for index, v := range one.Attrs6Key {
			one.Attrs6[v] = one.Attrs6Value[index]
		}
	}

	return
}
