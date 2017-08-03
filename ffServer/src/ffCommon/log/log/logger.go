package log

// Logger 接口
type Logger interface {
	// Printf Print Println 预期内的日志输出
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})

	// Stop stop or recover output
	Stop(stop bool)

	// Close close output
	Close()
}
