package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"path/filepath"
	"sort"

	proto "github.com/golang/protobuf/proto"
)

func transItem() {
	message := &Item{}

	// ItemTemplate
	ItemTemplateKeys := make([]int, 0, len(tomlItem.ItemTemplate)) // 必须使用64位机器
	//ItemTemplateKeys := make([]int, 0, len(tomlItem.ItemTemplate)) // 必须使用64位机器
	//ItemTemplateKeys := make([]string, 0, len(tomlItem.ItemTemplate)) // 必须使用64位机器
	for key := range tomlItem.ItemTemplate {
		ItemTemplateKeys = append(ItemTemplateKeys, int(key))
		//ItemTemplateKeys = append(ItemTemplateKeys, int(key))
		//ItemTemplateKeys = append(ItemTemplateKeys, string(key))
	}
	sort.Ints(ItemTemplateKeys)
	//sort.Ints(ItemTemplateKeys)
	//sort.Strings(ItemTemplateKeys)

	message.ItemTemplateKey = make([]int32, len(tomlItem.ItemTemplate))
	message.ItemTemplateValue = make([]*Item_StItemTemplate, len(tomlItem.ItemTemplate))
	for k, key := range ItemTemplateKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlItem.ItemTemplate[i]

		message.ItemTemplateKey[k] = i
		message.ItemTemplateValue[k] = &Item_StItemTemplate{
			Name:        v.Name,
			Desc:        v.Desc,
			SceneKey:    v.SceneKey,
			Icon:        v.Icon,
			Attrs1:      v.Attrs1,
			Attrs1Key:   v.Attrs1Key,
			Attrs1Value: v.Attrs1Value,
			Attrs2Value: v.Attrs2Value,
			Attrs3Value: v.Attrs3Value,
			Attrs6:      v.Attrs6,
			Attrs6Key:   v.Attrs6Key,
			Attrs6Value: v.Attrs6Value,
		}

		message.ItemTemplateValue[k].ItemType = int32(v.ItemType)
		message.ItemTemplateValue[k].SubType = make([]int32, len(v.SubType), len(v.SubType))
		for xx, yy := range v.SubType {
			message.ItemTemplateValue[k].SubType[xx] = int32(yy)
		}
		message.ItemTemplateValue[k].Attrs2Key = make([]int32, len(v.Attrs2Key), len(v.Attrs2Key))
		for xx, yy := range v.Attrs2Key {
			message.ItemTemplateValue[k].Attrs2Key[xx] = int32(yy)
		}
		message.ItemTemplateValue[k].Attrs3Key = make([]int32, len(v.Attrs3Key), len(v.Attrs3Key))
		for xx, yy := range v.Attrs3Key {
			message.ItemTemplateValue[k].Attrs3Key[xx] = int32(yy)
		}
		message.ItemTemplateValue[k].Attrs4Key = make([]int32, len(v.Attrs4Key), len(v.Attrs4Key))
		for xx, yy := range v.Attrs4Key {
			message.ItemTemplateValue[k].Attrs4Key[xx] = int32(yy)
		}
		message.ItemTemplateValue[k].Attrs4Value = make([]int32, len(v.Attrs4Value), len(v.Attrs4Value))
		for xx, yy := range v.Attrs4Value {
			message.ItemTemplateValue[k].Attrs4Value[xx] = int32(yy)
		}
		message.ItemTemplateValue[k].Attrs5Key = make([]int32, len(v.Attrs5Key), len(v.Attrs5Key))
		for xx, yy := range v.Attrs5Key {
			message.ItemTemplateValue[k].Attrs5Key[xx] = int32(yy)
		}
		message.ItemTemplateValue[k].Attrs5Value = make([]int32, len(v.Attrs5Value), len(v.Attrs5Value))
		for xx, yy := range v.Attrs5Value {
			message.ItemTemplateValue[k].Attrs5Value[xx] = int32(yy)
		}
	}

	pbBuf := proto.NewBuffer(make([]byte, 0, 1024*10))
	if err := pbBuf.Marshal(message); err != nil {
		log.RunLogger.Printf("transItem err[%v]", err)
		return
	}

	util.WriteFile(filepath.Join("ProtoBuf", "Client", "bytes", tomlItem.Name()+".bytes"), pbBuf.Bytes())
}

func init() {
	allTrans = append(allTrans, transItem)
}
