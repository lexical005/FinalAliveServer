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
	for i, key := range ItemTemplateKeys {
		k := int32(key)
		//k := int32(key)
		//k := int32(key)
		v := tomlItem.ItemTemplate[k]

		message.ItemTemplateKey[i] = k
		message.ItemTemplateValue[i] = &Item_StItemTemplate{
			Name: string(v.Name),
			Desc: string(v.Desc),
			SceneKey: string(v.SceneKey),
			Icon: string(v.Icon),
			ItemType: int32(v.ItemType),
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
