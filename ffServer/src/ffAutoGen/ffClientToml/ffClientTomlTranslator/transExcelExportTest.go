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
	message.VIPmap = make(map[int32]*ExcelExportTest_StVIPmap, len(tomlExcelExportTest.VIPmap))
	for _, key := range VIPmapKeys {
		k := int32(key)
		//k := int32(key)
		//k := int32(key)
		v := tomlExcelExportTest.VIPmap[k]

		message.VIPmap[k] = &ExcelExportTest_StVIPmap{
			InfoInt32: v.InfoInt32,
			InfoInt64: v.InfoInt64,
			InfoStr: v.InfoStr,
			InfoInt32Single: v.InfoInt32Single,
			InfoInt64Single: v.InfoInt64Single,
			InfoStrSingle: v.InfoStrSingle,
			InfoInt32Multi: v.InfoInt32Multi,
			InfoInt64Multi: v.InfoInt64Multi,
			InfoStrMulti: v.InfoStrMulti,
			ItemClientID: v.ItemClientID,
			Consume: transGrammar(v.Consume),
			EmptyInt32: v.EmptyInt32,
			EmptyInt64: v.EmptyInt64,
			EmptyStr: v.EmptyStr,
			EmptyIn32TSingle: v.EmptyIn32tSingle,
			EmptyInt64Single: v.EmptyInt64Single,
			EmptyStrSingle: v.EmptyStrSingle,
			EmptyInt32Multi: v.EmptyInt32Multi,
			EmptyInt64Multi: v.EmptyInt64Multi,
			EmptyStrMulti: v.EmptyStrMulti,
		}
	}

    // VIPstruct
	message.VIPstruct = &ExcelExportTest_StVIPstruct{
			InfoInt32: tomlExcelExportTest.VIPstruct.InfoInt32,
			InfoInt64: tomlExcelExportTest.VIPstruct.InfoInt64,
			InfoStr: tomlExcelExportTest.VIPstruct.InfoStr,
			InfoInt32Single: tomlExcelExportTest.VIPstruct.InfoInt32Single,
			InfoInt64Single: tomlExcelExportTest.VIPstruct.InfoInt64Single,
			InfoStrSingle: tomlExcelExportTest.VIPstruct.InfoStrSingle,
			InfoInt32Multi: tomlExcelExportTest.VIPstruct.InfoInt32Multi,
			InfoInt64Multi: tomlExcelExportTest.VIPstruct.InfoInt64Multi,
			InfoStrMulti: tomlExcelExportTest.VIPstruct.InfoStrMulti,
			ItemClientID: tomlExcelExportTest.VIPstruct.ItemClientID,
			Consume: transGrammar(tomlExcelExportTest.VIPstruct.Consume),
			EmptyInt32: tomlExcelExportTest.VIPstruct.EmptyInt32,
			EmptyInt64: tomlExcelExportTest.VIPstruct.EmptyInt64,
			EmptyStr: tomlExcelExportTest.VIPstruct.EmptyStr,
			EmptyIn32TSingle: tomlExcelExportTest.VIPstruct.EmptyIn32tSingle,
			EmptyInt64Single: tomlExcelExportTest.VIPstruct.EmptyInt64Single,
			EmptyStrSingle: tomlExcelExportTest.VIPstruct.EmptyStrSingle,
			EmptyInt32Multi: tomlExcelExportTest.VIPstruct.EmptyInt32Multi,
			EmptyInt64Multi: tomlExcelExportTest.VIPstruct.EmptyInt64Multi,
			EmptyStrMulti: tomlExcelExportTest.VIPstruct.EmptyStrMulti,
	}

	// VIPlist
	message.VIPlist = make([]*ExcelExportTest_StVIPlist, len(tomlExcelExportTest.VIPlist))
	for k, v := range tomlExcelExportTest.VIPlist {
		message.VIPlist[k] = &ExcelExportTest_StVIPlist{
			InfoInt32: v.InfoInt32,
			InfoInt64: v.InfoInt64,
			InfoStr: v.InfoStr,
			InfoInt32Single: v.InfoInt32Single,
			InfoInt64Single: v.InfoInt64Single,
			InfoStrSingle: v.InfoStrSingle,
			InfoInt32Multi: v.InfoInt32Multi,
			InfoInt64Multi: v.InfoInt64Multi,
			InfoStrMulti: v.InfoStrMulti,
			ItemClientID: v.ItemClientID,
			Consume: transGrammar(v.Consume),
			EmptyInt32: v.EmptyInt32,
			EmptyInt64: v.EmptyInt64,
			EmptyStr: v.EmptyStr,
			EmptyIn32TSingle: v.EmptyIn32tSingle,
			EmptyInt64Single: v.EmptyInt64Single,
			EmptyStrSingle: v.EmptyStrSingle,
			EmptyInt32Multi: v.EmptyInt32Multi,
			EmptyInt64Multi: v.EmptyInt64Multi,
			EmptyStrMulti: v.EmptyStrMulti,
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
