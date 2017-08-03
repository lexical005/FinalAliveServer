package log

type loggerEmpty struct {
}

// Printf calls l.logAsync to print to the loggerEmpty.
// Arguments are handled in the manner of fmt.Printf.
func (l *loggerEmpty) Printf(format string, v ...interface{}) {
}

// Print calls l.logAsync to print to the loggerEmpty.
// Arguments are handled in the manner of fmt.Print.
func (l *loggerEmpty) Print(v ...interface{}) {
}

// Println calls l.logAsync to print to the loggerEmpty.
// Arguments are handled in the manner of fmt.Println.
func (l *loggerEmpty) Println(v ...interface{}) {
}

// Stop stop or recover output. fatal not support.
func (l *loggerEmpty) Stop(stop bool) {
}

// Close close output. fatal not support.
func (l *loggerEmpty) Close() {
}
