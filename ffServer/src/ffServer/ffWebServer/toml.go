package main

import (
	"ffCommon/util"

	"github.com/lexical005/toml"
)

// 服务器配置
type applicationConfig struct {
	// 网络配置
	Net struct {
		// 连接数限制
		ConnectionLimit int

		// 监听客户端地址
		ListenAddr string
		ListenPort string
		DomainName string
	}
}

func readAppToml() error {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/config.toml")
	if err != nil {
		return err
	}

	// 解析
	err = toml.Unmarshal(fileContent, appConfig)
	if err != nil {
		return err
	}

	return nil
}
