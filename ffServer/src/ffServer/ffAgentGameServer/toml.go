package main

import (
	"ffCommon/util"

	"github.com/lexical005/toml"
)

// 服务器自身描述
type serverInfo struct {
	Channel    string // 渠道
	ServerType string // 服务器类型
	ServerID   int    // 服务器编号
}

// 连接配置
type sessionConfig struct {
	ReadDeadTime          int // ReadDeadTime 读取超时N秒. 为0时, 使用系统默认配置值60
	InitNetEventDataCount int // InitNetEventDataCount 初始创建多少网络事件数据缓存. 为0时, 使用的值为OnlineCount/4. 最小为2
	InitOnlineCount       int // InitOnlineCount 初始创建多少连接缓存, 必须配置. >=2
}

// serveConfig 为用户提供服务的配置
type serveConfig struct {
	// ListenTarget 监听目标
	ListenTarget string

	// ListenAddr 监听地址
	ListenAddr string

	// InitOnlineCount 初始多少同时连接存在
	InitOnlineCount int

	// SendExtraDataType 发送的协议的附加数据类型
	SendExtraDataType string

	// RecvExtraDataType 接收的协议的附加数据类型
	RecvExtraDataType string

	// AcceptNewSessionCache 接受新连接的管道的缓存大小. 影响接受新连接速度.
	AcceptNewSessionCache int

	// SessionNetEventDataCache 用户网络事件管道的缓存大小. 影响处理网络事件的速度.
	SessionNetEventDataCache int

	// SessionSendProtoCache 用户待发送协议管道的缓存大小. 影响发送协议的速度
	SessionSendProtoCache int
}

// 文本日志配置
type fileLoggerConfig struct {
	LoggerType      string // 日志类型
	RelativePath    string // 文本日志的存储相对路径
	FileLenLimit    int    // 单文本日志的大小限制
	RunLogger       bool   // 是否启用运行日志
	RunLoggerPrefix string // 运行日志文件的前缀
}

// 服务器配置
type applicationConfig struct {
	// Server 服务器自身描述
	Server serverInfo

	// Session 连接配置
	Session *sessionConfig

	// ServeUser 服务用户的配置
	ServeUser *serveConfig

	// Logger 日志配置
	Logger *fileLoggerConfig
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
