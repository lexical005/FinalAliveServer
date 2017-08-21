package log

// LoggerConfig 日志配置
type LoggerConfig struct {
	LoggerType string // 日志类型 console, file
	RunLogger  bool   // 是否启用运行日志, fata日志必然记录

	RelativePath    string // 文本日志的存储相对路径
	FileLenLimit    int    // 单文本日志的大小限制
	RunLoggerPrefix string // 运行日志文件的前缀
}
