package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"path/filepath"
	"sort"

    proto "github.com/golang/protobuf/proto"
)

func transExcelExportTest() {
    message := &ExcelExportTest{}

    // VIPmap
	VIPmapKeys := make([]int, 0, len(tomlExcelExportTest.VIPmap)) // 必须使用64位机器
	//VIPmapKeys := make([]int, 0, len(tomlExcelExportTest.VIPmap)) // 必须使用64位机器
	//VIPmapKeys := make([]string, 0, len(tomlExcelExportTest.VIPmap)) // 必须使用64位机器
	for key := range tomlExcelExportTest.VIPmap {
		VIPmapKeys = append(VIPmapKeys, int(key))
		//VIPmapKeys = append(VIPmapKeys, int(key))
		//VIPmapKeys = append(VIPmapKeys, string(key))
	}
	sort.Ints(VIPmapKeys)
	//sort.Ints(VIPmapKeys)
	//sort.Strings(VIPmapKeys)

	message.VIPmapKey = make([]int32, len(tomlExcelExportTest.VIPmap))
	message.VIPmapValue = make([]*ExcelExportTest_StVIPmap, len(tomlExcelExportTest.VIPmap))
	for i, key := range VIPmapKeys {
		k := int32(key)
		//k := int32(key)
		//k := int32(key)
		v := tomlExcelExportTest.VIPmap[k]

		message.VIPmapKey[i] = k
		message.VIPmapValue[i] = &ExcelExportTest_StVIPmap{
			InfoInt32: v.InfoInt32,
			InfoInt64: v.InfoInt64,
			InfoStr: v.InfoStr,
			InfoInt32Array: v.InfoInt32Array,
			InfoInt64Array: v.InfoInt64Array,
			InfoStringArray: v.InfoStringArray,
			ItemClientID: v.ItemClientID,
			Consume: transGrammar(v.Consume),
			EmptyInt32: v.EmptyInt32,
			EmptyInt64: v.EmptyInt64,
			EmptyStr: v.EmptyStr,
			EmptyInt32Array: v.EmptyInt32Array,
			EmptyInt64Array: v.EmptyInt64Array,
			EmptyStringArray: v.EmptyStringArray,
		}
	}

    // VIPstruct
	message.VIPstruct = &ExcelExportTest_StVIPstruct{
			InfoInt32: tomlExcelExportTest.VIPstruct.InfoInt32,
			InfoInt64: tomlExcelExportTest.VIPstruct.InfoInt64,
			InfoStr: tomlExcelExportTest.VIPstruct.InfoStr,
			InfoInt32Array: tomlExcelExportTest.VIPstruct.InfoInt32Array,
			InfoInt64Array: tomlExcelExportTest.VIPstruct.InfoInt64Array,
			InfoStringArray: tomlExcelExportTest.VIPstruct.InfoStringArray,
			ItemClientID: tomlExcelExportTest.VIPstruct.ItemClientID,
			Consume: transGrammar(tomlExcelExportTest.VIPstruct.Consume),
			EmptyInt32: tomlExcelExportTest.VIPstruct.EmptyInt32,
			EmptyInt64: tomlExcelExportTest.VIPstruct.EmptyInt64,
			EmptyStr: tomlExcelExportTest.VIPstruct.EmptyStr,
			EmptyInt32Array: tomlExcelExportTest.VIPstruct.EmptyInt32Array,
			EmptyInt64Array: tomlExcelExportTest.VIPstruct.EmptyInt64Array,
			EmptyStringArray: tomlExcelExportTest.VIPstruct.EmptyStringArray,
	}

	// VIPlist
	message.VIPlist = make([]*ExcelExportTest_StVIPlist, len(tomlExcelExportTest.VIPlist))
	for k, v := range tomlExcelExportTest.VIPlist {
		message.VIPlist[k] = &ExcelExportTest_StVIPlist{
			InfoInt32: v.InfoInt32,
			InfoInt64: v.InfoInt64,
			InfoStr: v.InfoStr,
			InfoInt32Array: v.InfoInt32Array,
			InfoInt64Array: v.InfoInt64Array,
			InfoStringArray: v.InfoStringArray,
			ItemClientID: v.ItemClientID,
			Consume: transGrammar(v.Consume),
			EmptyInt32: v.EmptyInt32,
			EmptyInt64: v.EmptyInt64,
			EmptyStr: v.EmptyStr,
			EmptyInt32Array: v.EmptyInt32Array,
			EmptyInt64Array: v.EmptyInt64Array,
			EmptyStringArray: v.EmptyStringArray,
		}
	}

    pbBuf := proto.NewBuffer(make([]byte, 0, 1024*10))
    if err := pbBuf.Marshal(message); err != nil {
        log.RunLogger.Printf("transExcelExportTest err[%v]", err)
        return
    }

    util.WriteFile(filepath.Join("ProtoBuf", "Client", "bytes", tomlExcelExportTest.Name()+".bytes"), pbBuf.Bytes())
}

func init() {
    allTrans = append(allTrans, transExcelExportTest)
}
