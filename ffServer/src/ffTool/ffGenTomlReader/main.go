package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// main 将excel导出的go代码文件导入到一个工程, 生成exe, 读取配置文件, 以检测是否正常
func main() {
	defer util.PanicProtect(func(isPanic bool) {
		if isPanic {
			log.RunLogger.Println("异常退出, 以上是错误堆栈")
			<-time.After(time.Hour)
		}
	}, "ffGenTomlReader")

	// 命令行参数解析
	gocodedir := flag.String("gocodedir", "", "golang read toml code directory")
	readername := flag.String("readername", "", "reader directory")
	proto := flag.String("proto", "", "gen go to proto reader")
	csharp := flag.String("csharp", "", "gen csharp code trans")
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
	golangFiles, goFullPathFiles := make([]string, 0, 10), make([]string, 0, 10)
	if fi, err := os.Stat(ffGameConfigPath); err != nil && os.IsExist(err) || fi != nil && fi.IsDir() {
		err := util.Walk(ffGameConfigPath, func(f os.FileInfo) error {
			// 忽略文件夹以及非go文件
			name := f.Name()
			if f.IsDir() || !strings.HasSuffix(name, ".go") {
				return nil
			}

			// 非配置文件
			if name == "Error.go" || name == "Enum.go" || name[0] == strings.ToLower(name)[0] {
				return nil
			}

			golangFiles = append(golangFiles, name[0:len(name)-len(".go")])
			goFullPathFiles = append(goFullPathFiles, filepath.Join(ffGameConfigPath, name))

			return nil
		})

		if err != nil {
			log.RunLogger.Printf("walk gocodedir[%v] get error[%v]", ffGameConfigPath, err)
			return
		}
	}

	_, goCodePackageName := filepath.Split(*gocodedir)

	// go读取toml数据格式的代码 ==> ProtoBuf读取toml数据格式的代码
	if *proto == "proto" {
		// 生成读取文件
		genReadTomlCode(filepath.Join(ffGameConfigPath, *readername, "read_toml.go"), golangFiles, goCodePackageName)

		transGoToProto(
			filepath.Join(ffGameConfigPath, *readername),
			filepath.Join(ffGameConfigPath, *readername, "Config.pb.go"),
			goFullPathFiles,
			goCodePackageName)
	}

	// 客户端charp代码读取转换封装
	if *csharp == "csharp" {
		genProtoCSharpReaderCode(
			filepath.Join("ProtoBuf", "Client"),
			filepath.Join("ProtoBuf", "Client", "Config.cs"))
	}
}
