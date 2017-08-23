package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/util"
	"fmt"

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

	// ConnectLogin 连接到http服务器进行登录验证
	ConnectLogin *base.HTTPClientConfig

	// Logger 日志配置
	Logger *log.LoggerConfig
}

// Check 配置检查
func (config *applicationConfig) Check() (err error) {
	err = config.Session.Check()
	if err != nil {
		return
	}

	if config.ConnectMatchServer.KeepAliveInterval < 1 || config.ConnectMatchServer.KeepAliveInterval > config.Session.ReadDeadTime/2 {
		err = fmt.Errorf("invalid KeepAliveInterval[%v], must >0 and <= [%v](SessionConfig.ReadDeadTime/2)",
			config.ConnectMatchServer.KeepAliveInterval, config.Session.ReadDeadTime)
		return
	}

	return nil
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
