package log

import (
	"bytes"
	"fmt"
	"runtime"
)

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

// Stack 获取当前堆栈描述
func Stack() string {
	var buf bytes.Buffer

	buf.WriteString("Stack\n")

	i := 0
	funcName, file, line, ok := runtime.Caller(i)
	for ok {
		buf.WriteString(fmt.Sprintf("[%d, %s, %s, %d]\n", i, runtime.FuncForPC(funcName).Name(), file, line))
		i++
		funcName, file, line, ok = runtime.Caller(i)
	}
	buf.WriteString("\n")

	return buf.String()
}
