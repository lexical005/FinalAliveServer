package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"path/filepath"
	"sort"

    proto "github.com/golang/protobuf/proto"
)

func transLanguage() {
    message := &Language{}

    // Common
	//CommonKeys := make([]int, 0, len(tomlLanguage.Common)) // 必须使用64位机器
	//CommonKeys := make([]int, 0, len(tomlLanguage.Common)) // 必须使用64位机器
	CommonKeys := make([]string, 0, len(tomlLanguage.Common)) // 必须使用64位机器
	for key := range tomlLanguage.Common {
		//CommonKeys = append(CommonKeys, int(key))
		//CommonKeys = append(CommonKeys, int(key))
		CommonKeys = append(CommonKeys, string(key))
	}
	//sort.Ints(CommonKeys)
	//sort.Ints(CommonKeys)
	sort.Strings(CommonKeys)

	message.CommonKey = make([]string, len(tomlLanguage.Common))
	message.CommonValue = make([]*Language_StCommon, len(tomlLanguage.Common))
	for i, key := range CommonKeys {
		//k := string(key)
		//k := string(key)
		k := string(key)
		v := tomlLanguage.Common[k]

		message.CommonKey[i] = k
		message.CommonValue[i] = &Language_StCommon{
			CN: v.CN,
		}
	}

	// Special
	message.Special = make([]*Language_StSpecial, len(tomlLanguage.Special))
	for k, v := range tomlLanguage.Special {
		message.Special[k] = &Language_StSpecial{
			CN: v.CN,
		}
	}

    // Error
	//ErrorKeys := make([]int, 0, len(tomlLanguage.Error)) // 必须使用64位机器
	//ErrorKeys := make([]int, 0, len(tomlLanguage.Error)) // 必须使用64位机器
	ErrorKeys := make([]string, 0, len(tomlLanguage.Error)) // 必须使用64位机器
	for key := range tomlLanguage.Error {
		//ErrorKeys = append(ErrorKeys, int(key))
		//ErrorKeys = append(ErrorKeys, int(key))
		ErrorKeys = append(ErrorKeys, string(key))
	}
	//sort.Ints(ErrorKeys)
	//sort.Ints(ErrorKeys)
	sort.Strings(ErrorKeys)

	message.ErrorKey = make([]string, len(tomlLanguage.Error))
	message.ErrorValue = make([]*Language_StError, len(tomlLanguage.Error))
	for i, key := range ErrorKeys {
		//k := string(key)
		//k := string(key)
		k := string(key)
		v := tomlLanguage.Error[k]

		message.ErrorKey[i] = k
		message.ErrorValue[i] = &Language_StError{
			CN: v.CN,
		}
	}

    pbBuf := proto.NewBuffer(make([]byte, 0, 1024*10))
    if err := pbBuf.Marshal(message); err != nil {
        log.RunLogger.Printf("transLanguage err[%v]", err)
        return
    }

    util.WriteFile(filepath.Join("ProtoBuf", "Client", "bytes", tomlLanguage.Name()+".bytes"), pbBuf.Bytes())
}

func init() {
    allTrans = append(allTrans, transLanguage)
}
