package log

import "fmt"

type loggerConsoleFatal struct {
}

// Printf calls fmt.Printf to print to the loggerConsoleFatal.
// Arguments are handled in the manner of fmt.Printf.
func (l *loggerConsoleFatal) Printf(format string, v ...interface{}) {
	s := Stack()

	if len(format) > 0 && format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Println("Fatal:")
	fmt.Printf(format, v...)
	fmt.Println(s)

	RunLogger.Printf(format, v...)
	RunLogger.Println(s)
}

// Print calls fmt.Print to print to the loggerConsoleFatal.
// Arguments are handled in the manner of fmt.Print.
func (l *loggerConsoleFatal) Print(v ...interface{}) {
	s := Stack()

	fmt.Println("Fatal:")
	fmt.Print(v...)
	fmt.Println(s)

	RunLogger.Print(v...)
	RunLogger.Println(s)
}

// Println calls fmt.Println to print to the loggerConsoleFatal.
// Arguments are handled in the manner of fmt.Println.
func (l *loggerConsoleFatal) Println(v ...interface{}) {
	s := Stack()

	fmt.Println("Fatal:")
	fmt.Println(v...)
	fmt.Println(s)

	RunLogger.Println(v...)
	RunLogger.Println(s)
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
