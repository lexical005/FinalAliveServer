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
			InfoInt32: int32(v.InfoInt32),
			InfoInt64: int64(v.InfoInt64),
			InfoStr: string(v.InfoStr),
			InfoInt32Array: []int32(v.InfoInt32Array),
			InfoInt64Array: []int64(v.InfoInt64Array),
			InfoStringArray: []string(v.InfoStringArray),
			ItemClientID: int32(v.ItemClientID),
			Consume: transGrammar(v.Consume),
			EmptyInt32: int32(v.EmptyInt32),
			EmptyInt64: int64(v.EmptyInt64),
			EmptyStr: string(v.EmptyStr),
			EmptyInt32Array: []int32(v.EmptyInt32Array),
			EmptyInt64Array: []int64(v.EmptyInt64Array),
			EmptyStringArray: []string(v.EmptyStringArray),
		}
	}

    // VIPstruct
	message.VIPstruct = &ExcelExportTest_StVIPstruct{
			InfoInt32: int32(tomlExcelExportTest.VIPstruct.InfoInt32),
			InfoInt64: int64(tomlExcelExportTest.VIPstruct.InfoInt64),
			InfoStr: string(tomlExcelExportTest.VIPstruct.InfoStr),
			InfoInt32Array: []int32(tomlExcelExportTest.VIPstruct.InfoInt32Array),
			InfoInt64Array: []int64(tomlExcelExportTest.VIPstruct.InfoInt64Array),
			InfoStringArray: []string(tomlExcelExportTest.VIPstruct.InfoStringArray),
			ItemClientID: int32(tomlExcelExportTest.VIPstruct.ItemClientID),
			Consume: transGrammar(tomlExcelExportTest.VIPstruct.Consume),
			EmptyInt32: int32(tomlExcelExportTest.VIPstruct.EmptyInt32),
			EmptyInt64: int64(tomlExcelExportTest.VIPstruct.EmptyInt64),
			EmptyStr: string(tomlExcelExportTest.VIPstruct.EmptyStr),
			EmptyInt32Array: []int32(tomlExcelExportTest.VIPstruct.EmptyInt32Array),
			EmptyInt64Array: []int64(tomlExcelExportTest.VIPstruct.EmptyInt64Array),
			EmptyStringArray: []string(tomlExcelExportTest.VIPstruct.EmptyStringArray),
	}

	// VIPlist
	message.VIPlist = make([]*ExcelExportTest_StVIPlist, len(tomlExcelExportTest.VIPlist))
	for k, v := range tomlExcelExportTest.VIPlist {
		message.VIPlist[k] = &ExcelExportTest_StVIPlist{
			InfoInt32: int32(v.InfoInt32),
			InfoInt64: int64(v.InfoInt64),
			InfoStr: string(v.InfoStr),
			InfoInt32Array: []int32(v.InfoInt32Array),
			InfoInt64Array: []int64(v.InfoInt64Array),
			InfoStringArray: []string(v.InfoStringArray),
			ItemClientID: int32(v.ItemClientID),
			Consume: transGrammar(v.Consume),
			EmptyInt32: int32(v.EmptyInt32),
			EmptyInt64: int64(v.EmptyInt64),
			EmptyStr: string(v.EmptyStr),
			EmptyInt32Array: []int32(v.EmptyInt32Array),
			EmptyInt64Array: []int64(v.EmptyInt64Array),
			EmptyStringArray: []string(v.EmptyStringArray),
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
