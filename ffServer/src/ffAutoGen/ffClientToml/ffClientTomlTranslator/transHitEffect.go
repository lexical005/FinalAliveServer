package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"path/filepath"
	"sort"

	proto "github.com/golang/protobuf/proto"
)

func transHitEffect() {
	message := &HitEffect{}

	// Hit
	//HitKeys := make([]int, 0, len(tomlHitEffect.Hit)) // 必须使用64位机器
	//HitKeys := make([]int, 0, len(tomlHitEffect.Hit)) // 必须使用64位机器
	HitKeys := make([]string, 0, len(tomlHitEffect.Hit)) // 必须使用64位机器
	for key := range tomlHitEffect.Hit {
		//HitKeys = append(HitKeys, int(key))
		//HitKeys = append(HitKeys, int(key))
		HitKeys = append(HitKeys, string(key))
	}
	//sort.Ints(HitKeys)
	//sort.Ints(HitKeys)
	sort.Strings(HitKeys)

	message.HitKey = make([]string, len(tomlHitEffect.Hit))
	message.HitValue = make([]*HitEffect_StHit, len(tomlHitEffect.Hit))
	for k, key := range HitKeys {
		//i := string(key)
		//i := string(key)
		i := string(key)
		v := tomlHitEffect.Hit[i]

		message.HitKey[k] = i
		message.HitValue[k] = &HitEffect_StHit{
			Gun: v.Gun,
		}

	}

	pbBuf := proto.NewBuffer(make([]byte, 0, 1024*10))
	if err := pbBuf.Marshal(message); err != nil {
		log.RunLogger.Printf("transHitEffect err[%v]", err)
		return
	}

	util.WriteFile(filepath.Join("ProtoBuf", "Client", "bytes", tomlHitEffect.Name()+".bytes"), pbBuf.Bytes())
}

func init() {
	allTrans = append(allTrans, transHitEffect)
}
