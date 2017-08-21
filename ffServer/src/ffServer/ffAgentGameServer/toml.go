package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/util"

	"github.com/lexical005/toml"
)

// 服务器配置
type applicationConfig struct {
	// Server 服务器自身描述
	Server *base.ServerInfo

	// Session 连接配置
	Session *base.SessionConfig

	// ServeUser 服务用户的配置
	ServeUser *base.ServeConfig

	// ConnectMatchServer 连接MatchServer
	ConnectMatchServer *base.ConnectConfig

	// Logger 日志配置
	Logger *log.LoggerConfig
}

func readToml() error {
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
