package logfile

import p "ffCommon/pool"
import "ffCommon/log/log"

const (
	logsCacheCount = 16
	logPrefixSize  = 16
	logBufSize     = 256

	// DefaultLogFileLengthLimit 默认日志文件的大小限制
	DefaultLogFileLengthLimit = 1024 * 1024 * 100 // 100MB

	// DefaultLogFileRelativePath 默认日志文件相对启动程序的存储位置
	DefaultLogFileRelativePath = "log" // 日志文件存储相对路径

	// DefaultLogFileFatalPrefix 致命日志文件的前缀
	DefaultLogFileFatalPrefix = "fatal"
)

var lrq *logRequestPool

func init() {
	creator := func() interface{} {
		return &logRequest{
			buf: make([]byte, 0, logBufSize),
		}
	}

	lrq = &logRequestPool{
		pool: p.New("log.lrq.pool", false, creator, 100, 50),
	}
}

// Init 以logfile方式初始化log模块
// 不启用的logger, 将使用空白logger覆盖默认的logger
func Init(
	logFileRelativePath string, logFileLengthLimit int,
	runLog bool, runFilePrefix string,
	fatalLog bool, fatalFilePrefix string) (err error) {

	var runLogger, fatalLogger log.Logger

	if runLog {
		runLogger, err = newFileLoggerNormal(logFileRelativePath, runFilePrefix, logFileLengthLimit)
		if err != nil {
			return err
		}
	} else {
		runLogger = log.NewLoggerEmpty()
	}

	if fatalLog {
		fatalLogger, err = newFileLoggerFatal(logFileRelativePath, fatalFilePrefix, logFileLengthLimit)
		if err != nil {
			return err
		}
	} else {
		fatalLogger = log.NewLoggerEmpty()
	}

	return log.Init(runLogger, fatalLogger)
}

// InitRunLog 以logfile方式初始化log模块RunLogger
func InitRunLog(logFileRelativePath string, logFileLengthLimit int, runFilePrefix string) (err error) {
	runLogger, err := newFileLoggerNormal(logFileRelativePath, runFilePrefix, logFileLengthLimit)
	if err != nil {
		return err
	}

	return log.Init(runLogger, log.FatalLogger)
}

// InitFatalLog 以logfile方式初始化log模块FatalLogger
func InitFatalLog(logFileRelativePath string, logFileLengthLimit int, fatalFilePrefix string) (err error) {
	fatalLogger, err := newFileLoggerFatal(logFileRelativePath, fatalFilePrefix, logFileLengthLimit)
	if err != nil {
		return err
	}

	return log.Init(log.RunLogger, fatalLogger)
}
