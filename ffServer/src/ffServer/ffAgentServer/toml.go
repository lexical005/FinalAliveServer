package main

import (
	"ffCommon/util"

	"github.com/lexical005/toml"
)

// server配置
type serverConfig struct {
	ServerType string
	ServerID   int
}

// 连接配置
type sessionConfig struct {
	ReadDeadTime int // 读取超时N秒
	OnlineCount  int // 默认的同时在线连接数
}

// 客户端监听配置
type clientListenConfig struct {
	ClientType string // 客户端类型
	ListenAddr string // 监听地址
}

// 服务端监听配置
type serverListenConfig struct {
	ServerType string // 服务端类型
	ListenAddr string // 监听地址
}

// 文本日志配置
type fileLoggerConfig struct {
	RelativePath    string // 文本日志的存储相对路径
	FileLenLimit    int    // 单文本日志的大小限制
	RunLogger       bool   // 是否启用运行日志
	RunLoggerPrefix string // 运行日志文件的前缀
}

// 服务器配置
type applicationConfig struct {
	ServerConfig serverConfig

	Session sessionConfig

	ClientListen clientListenConfig

	ServerListen serverListenConfig

	Logger fileLoggerConfig
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
