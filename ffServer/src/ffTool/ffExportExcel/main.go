package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"ffExcel"

	"bufio"
	"os"

	"github.com/lexical005/toml"
)

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
	exportConfig := &ffExcel.ExportConfig{}
	err = toml.Unmarshal(fileContent, exportConfig)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	err = ffExcel.ExportExcelDir("excel", exportConfig)
	if err != nil {
		log.RunLogger.Println(err)

		log.RunLogger.Println("\n请处理以上错误, 然后再次执行!")

		inputReader := bufio.NewReader(os.Stdin)
		inputReader.ReadString('\n')

		return
	}
}
