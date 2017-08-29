package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"path/filepath"
	"sort"

	proto "github.com/golang/protobuf/proto"
)

func transAsset() {
	message := &Asset{}

	// Assets
	AssetsKeys := make([]int, 0, len(tomlAsset.Assets)) // 必须使用64位机器
	//AssetsKeys := make([]int, 0, len(tomlAsset.Assets)) // 必须使用64位机器
	//AssetsKeys := make([]string, 0, len(tomlAsset.Assets)) // 必须使用64位机器
	for key := range tomlAsset.Assets {
		AssetsKeys = append(AssetsKeys, int(key))
		//AssetsKeys = append(AssetsKeys, int(key))
		//AssetsKeys = append(AssetsKeys, string(key))
	}
	sort.Ints(AssetsKeys)
	//sort.Ints(AssetsKeys)
	//sort.Strings(AssetsKeys)

	message.AssetsKey = make([]int32, len(tomlAsset.Assets))
	message.AssetsValue = make([]*Asset_StAssets, len(tomlAsset.Assets))
	for k, key := range AssetsKeys {
		i := int32(key)
		//i := int32(key)
		//i := int32(key)
		v := tomlAsset.Assets[i]

		message.AssetsKey[k] = i
		message.AssetsValue[k] = &Asset_StAssets{
			BattleDefault: v.BattleDefault,
			HomeDefault:   v.HomeDefault,
			SceneDefault:  v.SceneDefault,
		}

	}

	pbBuf := proto.NewBuffer(make([]byte, 0, 1024*10))
	if err := pbBuf.Marshal(message); err != nil {
		log.RunLogger.Printf("transAsset err[%v]", err)
		return
	}

	util.WriteFile(filepath.Join("ProtoBuf", "Client", "bytes", tomlAsset.Name()+".bytes"), pbBuf.Bytes())
}

func init() {
	allTrans = append(allTrans, transAsset)
}
