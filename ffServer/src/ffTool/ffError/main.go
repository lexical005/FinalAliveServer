package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"ffExcel"

	"bufio"
	"os"
	"path"

	"github.com/lexical005/toml"
)

type errConfig struct {
	Common ffExcel.ExportConfig

	ErrCodeExportDefPath string
}

type errReasonToml struct {
	Error []struct {
		ErrCode string
		CN      string
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			util.PrintPanicStack(err)

			log.RunLogger.Println("\n请处理以上错误, 然后再次执行!")

			inputReader := bufio.NewReader(os.Stdin)
			inputReader.ReadString('\n')
		}
	}()

	// 读取配置文件内容
	fileContent, err := util.ReadFile("config.toml")
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	// 解析配置文件
	errConfig := &errConfig{}
	err = toml.Unmarshal(fileContent, errConfig)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	// 将excel配置转换为程序读取的配置
	err = ffExcel.ExportExcelDir("excel", &errConfig.Common)
	if err != nil {
		log.RunLogger.Println(err)

		log.RunLogger.Println("\n请处理以上错误, 然后再次执行!")

		inputReader := bufio.NewReader(os.Stdin)
		inputReader.ReadString('\n')

		return
	}

	// 读取导出的toml配置
	fileContent, err = util.ReadFile(path.Join("toml", "server", "Error.toml"))
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	// 解析导出的toml配置文件
	errReasonToml := &errReasonToml{}
	err = toml.Unmarshal(fileContent, errReasonToml)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	// 服务端错误枚举
	errGoDef := tomlToGolang(errReasonToml)
	util.WriteFile(path.Join(errConfig.ErrCodeExportDefPath, "error.go"), []byte(errGoDef))

	// 客户端错误枚举
	errCSharpDef := tomlToCSharp(errReasonToml)
	util.WriteFile(path.Join(errConfig.Common.ClientExportCSharpCodePath, "Error.cs"), []byte(errCSharpDef))
}
