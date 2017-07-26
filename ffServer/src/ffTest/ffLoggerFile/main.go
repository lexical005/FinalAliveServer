package main

import (
	"ffCommon/log/log"
	"ffCommon/log/logfile"
)

func main() {
	logfile.Init(logfile.DefaultLogFileRelativePath, logfile.DefaultLogFileLengthLimit, true, "net", true, "run", true, "sql")
	log.RunLogger.Println("good")
}
