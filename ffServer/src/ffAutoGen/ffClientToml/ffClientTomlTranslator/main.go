package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"ffLogic/ffGrammar"

	proto "github.com/golang/protobuf/proto"
)

func main() {
	defer util.PanicProtect()

	read()

	trans()
}

// 将toml数据转换为pb字节流
func trans() {
	transExcelExportTest()
}

func transGrammar(grammar ffGrammar.Grammar) *Grammar {
	return &Grammar{
		Grammar: grammar.Origin(),
	}
}

func transExcelExportTest() {
	defer util.PanicProtect("transExcelExportTest")

	pbExcelExportTest := &ExcelExportTest{}

	pbExcelExportTest.VIPmap = make(map[int32]*ExcelExportTest_StVIPmap, len(tomlExcelExportTest.VIPmap))
	for k, v := range tomlExcelExportTest.VIPmap {
		pbVIPmap := &ExcelExportTest_StVIPmap{
			InfoInt: int32(v.InfoInt),
			InfoStr: v.InfoStr,
			ItemID:  int32(v.ItemID),
			Award:   transGrammar(v.Award),
		}
		pbExcelExportTest.VIPmap[int32(k)] = pbVIPmap
	}

	pbExcelExportTest.VIPstruct = &ExcelExportTest_StVIPstruct{
		InfoInt: int32(tomlExcelExportTest.VIPstruct.InfoInt),
		InfoStr: tomlExcelExportTest.VIPstruct.InfoStr,
		ItemID:  int32(tomlExcelExportTest.VIPstruct.ItemID),
	}

	pbExcelExportTest.VIPlist = make([]*ExcelExportTest_StVIPlist, len(tomlExcelExportTest.VIPlist), len(tomlExcelExportTest.VIPlist))
	for k, v := range tomlExcelExportTest.VIPlist {
		pbVIPlist := &ExcelExportTest_StVIPlist{
			InfoInt: int32(v.InfoInt),
			InfoStr: v.InfoStr,
			ItemID:  int32(v.ItemID),
		}
		pbExcelExportTest.VIPlist[k] = pbVIPlist
	}

	pbBuf := proto.NewBuffer(make([]byte, 0, 1024*10))
	if err := pbBuf.Marshal(pbExcelExportTest); err != nil {
		log.RunLogger.Printf("transExcelExportTest err[%v]", err)
		return
	}

	util.WriteFile(tomlExcelExportTest.Name()+".bytes", pbBuf.Bytes())
}
