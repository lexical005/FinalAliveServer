package log

// RunLogger 运行日志
var RunLogger Logger

// FatalLogger 致命日志
var FatalLogger Logger

func init() {
	RunLogger = newLoggerConsoleNormal()
	FatalLogger = newLoggerConsoleFatal()
}

// Init 初始化log模块
func Init(runLogger, faltalLogger Logger) (err error) {
	RunLogger, FatalLogger = runLogger, faltalLogger
	return err
}

var emptyLogger = &loggerEmpty{}

// NewLoggerEmpty return a empty Logger
func NewLoggerEmpty() Logger {
	return emptyLogger
}
