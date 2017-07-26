package main

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"flag"
	"os"
	"strings"
)

func main() {
	// 命令行参数解析
	frm := flag.String("frm", "", "file or directory to compress")
	dst := flag.String("dst", "", "compress to file with suffix")
	flag.Parse()

	if *frm == "" || *dst == "" || strings.Index(*dst, ".") == -1 {
		log.RunLogger.Printf("invalid input params. frm[%s] dst[%s]\n", *frm, *dst)
		return
	} else if fi, _ := os.Stat(*frm); fi == nil {
		// 源文件或源目录必须存在
		log.RunLogger.Printf("invalid frm params(not exist). frm[%s] dst[%s]\n", *frm, *dst)
		return
	}

	// 压缩
	if err := util.CompressZIP(*frm, *dst); err != nil {
		log.RunLogger.Println(err)
	}
}
