package log

import "runtime"
import "fmt"

type loggerConsoleNormal struct {
	run bool
}

// Printf calls fmt.Printf to print to the loggerConsoleNormal.
// Arguments are handled in the manner of fmt.Printf.
func (l *loggerConsoleNormal) Printf(format string, v ...interface{}) {
	if l.run {
		if len(format) > 0 && format[len(format)-1] != '\n' {
			format += "\n"
		}
		fmt.Printf(format, v...)
	}
}

// Print calls fmt.Print to print to the loggerConsoleNormal.
// Arguments are handled in the manner of fmt.Print.
func (l *loggerConsoleNormal) Print(v ...interface{}) {
	if l.run {
		fmt.Print(v...)
	}
}

// Println calls fmt.Println to print to the loggerConsoleNormal.
// Arguments are handled in the manner of fmt.Println.
func (l *loggerConsoleNormal) Println(v ...interface{}) {
	if l.run {
		fmt.Println(v...)
	}
}

// DumpStack dump caller function call stack
func (l *loggerConsoleNormal) DumpStack() {
	i := 0
	funcName, file, line, ok := runtime.Caller(i)
	for ok {
		fmt.Printf("[%d, %s, %s, %d]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
		i++
		funcName, file, line, ok = runtime.Caller(i)
	}
	fmt.Println()
}

// Stop stop or recover output. fatal not support.
func (l *loggerConsoleNormal) Stop(stop bool) {
	l.run = !stop
}

// Close close output. fatal not support.
func (l *loggerConsoleNormal) Close() {
	// 不再有任何输出
	l.run = false
}

func newLoggerConsoleNormal() Logger {
	return &loggerConsoleNormal{
		run: true,
	}
}
