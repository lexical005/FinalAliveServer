package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"path/filepath"

	proto "github.com/golang/protobuf/proto"
)

func transAsset() {
	message := &Asset{}

	// Assets
	message.Assets = make([]*Asset_StAssets, len(tomlAsset.Assets))
	for k, v := range tomlAsset.Assets {
		message.Assets[k] = &Asset_StAssets{
			TemplateID:    v.TemplateID,
			BattleDefault: v.BattleDefault,
			HomeDefault:   v.HomeDefault,
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
