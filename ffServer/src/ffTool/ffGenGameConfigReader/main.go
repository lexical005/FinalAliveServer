package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

// main 将excel导出的go文件导入到一个工程, 生成exe, 读取配置文件, 以检测是否正常
func main() {
	defer util.PanicProtect()

	// 命令行参数解析
	dir := flag.String("dir", "", "golang ffGameConfig directory")
	flag.Parse()

	if *dir == "" {
		log.RunLogger.Printf("invalid input params: dir[%v]", *dir)
		return
	}

	ffGameConfigPath, err := filepath.Abs(*dir)
	if err != nil {
		log.RunLogger.Printf("invalid input params: dir[%v] err[%v]", *dir, err)
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
			log.RunLogger.Printf("walk dir[%v] get error[%v]", ffGameConfigPath, err)
			return
		}
	}

	// 生成读取文件
	ffGameConfigReaderFile := filepath.Join(ffGameConfigPath, "ffGameConfigReader", "main.go")
	genGameConfigReader(ffGameConfigReaderFile, golangFiles)
}
