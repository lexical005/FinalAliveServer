package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

// main 将excel导出的go代码文件导入到一个工程, 生成exe, 读取配置文件, 以检测是否正常
func main() {
	defer util.PanicProtect()

	// 命令行参数解析
	gocodedir := flag.String("gocodedir", "", "golang read toml code directory")
	readername := flag.String("readername", "", "reader directory")
	flag.Parse()

	if *gocodedir == "" || *readername == "" {
		log.RunLogger.Printf("invalid input params: gocodedir[%v] readername[%v]", *gocodedir, *readername)
		return
	}

	ffGameConfigPath, err := filepath.Abs(*gocodedir)
	if err != nil {
		log.RunLogger.Printf("invalid input params: gocodedir[%v] err[%v]", *gocodedir, err)
		return
	}

	// 遍历go文件
	golangFiles := make([]string, 0, 10)
	if fi, err := os.Stat(ffGameConfigPath); err != nil && os.IsExist(err) || fi != nil && fi.IsDir() {
		err := util.Walk(ffGameConfigPath, func(f os.FileInfo) error {
			// 忽略文件夹以及非go文件
			name := f.Name()
			if f.IsDir() || !strings.HasSuffix(name, ".go") {
				return nil
			}

			golangFiles = append(golangFiles, name[0:len(name)-len(".go")])

			return nil
		})

		if err != nil {
			log.RunLogger.Printf("walk gocodedir[%v] get error[%v]", ffGameConfigPath, err)
			return
		}
	}

	_, packageName := filepath.Split(*gocodedir)

	// 生成读取文件
	genReadTomlCode(filepath.Join(ffGameConfigPath, *readername, "read.go"), golangFiles, packageName)
}
