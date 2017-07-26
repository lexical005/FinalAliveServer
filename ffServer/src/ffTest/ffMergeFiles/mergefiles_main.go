package main

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"flag"
	"os"
)

func main() {
	// 命令行参数解析
	dir := flag.String("dir", "", "files to merge")
	flag.Parse()

	if *dir == "" {
		log.RunLogger.Printf("invalid input params. dir[%s]\n", *dir)
		return
	}

	if fi, _ := os.Stat(*dir); fi == nil || !fi.IsDir() {
		log.RunLogger.Printf("dir not exist. dir[%s]\n", *dir)
		return
	}

	// 合并目录内的文件
	if err := util.MergeFiles(*dir); err != nil {
		log.RunLogger.Println(err)
	}
}
