package log

import "runtime"
import "fmt"

type loggerConsoleFatal struct {
}

// Printf calls fmt.Printf to print to the loggerConsoleFatal.
// Arguments are handled in the manner of fmt.Printf.
func (l *loggerConsoleFatal) Printf(format string, v ...interface{}) {
	fmt.Println("Fatal:")
	if len(format) > 0 && format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Printf(format, v...)
	l.DumpStack()
}

// Print calls fmt.Print to print to the loggerConsoleFatal.
// Arguments are handled in the manner of fmt.Print.
func (l *loggerConsoleFatal) Print(v ...interface{}) {
	fmt.Println("Fatal:")
	fmt.Print(v...)
	l.DumpStack()
}

// Println calls fmt.Println to print to the loggerConsoleFatal.
// Arguments are handled in the manner of fmt.Println.
func (l *loggerConsoleFatal) Println(v ...interface{}) {
	fmt.Println("Fatal:")
	fmt.Println(v...)
	l.DumpStack()
}

// DumpStack dump caller function call stack
func (l *loggerConsoleFatal) DumpStack() {
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
func (l *loggerConsoleFatal) Stop(stop bool) {
}

// Close close output. fatal not support.
func (l *loggerConsoleFatal) Close() {
}

func newLoggerConsoleFatal() Logger {
	return &loggerConsoleFatal{}
}
