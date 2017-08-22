package base

// ServeConfig 服务配置
type ServeConfig struct {
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

	// SessionNetEventDataCache 网络事件管道的缓存大小. 影响处理网络事件的速度.
	SessionNetEventDataCache int

	// SessionSendProtoCache 待发送协议管道的缓存大小. 影响发送协议的速度
	SessionSendProtoCache int
}

// ConnectConfig 连接配置
type ConnectConfig struct {
	// ConnectTarget 连接目标
	ConnectTarget string

	// ConnectAddr 连接地址
	ConnectAddr string

	// SendExtraDataType 发送的协议的附加数据类型
	SendExtraDataType string

	// RecvExtraDataType 接收的协议的附加数据类型
	RecvExtraDataType string

	// SessionNetEventDataCache 网络事件管道的缓存大小. 影响处理网络事件的速度.
	SessionNetEventDataCache int

	// SessionSendProtoCache 待发送协议管道的缓存大小. 影响发送协议的速度
	SessionSendProtoCache int
}

// ServerInfo 服务器自身描述
type ServerInfo struct {
	Channel    string // 渠道
	ServerType string // 服务器类型
	ServerID   int32  // 服务器编号
}

// SessionConfig 连接配置
type SessionConfig struct {
	ReadDeadTime          int // ReadDeadTime 读取超时N秒. 为0时, 使用系统默认配置值60
	InitNetEventDataCount int // InitNetEventDataCount 初始创建多少网络事件数据缓存. 为0时, 使用的值为OnlineCount/4. 最小为2
	InitOnlineCount       int // InitOnlineCount 初始创建多少连接缓存, 必须配置. >=2
}

// WebServerConfig WebServer配置
type WebServerConfig struct {
	// ListenAddr 监听地址
	ListenAddr string
}

// HTTPClientConfig WebClient配置
type HTTPClientConfig struct {
	// URL 连接地址
	URL string

	// RequestWorkerCount 使用几个协程发送验证请求
	RequestWorkerCount int
}
