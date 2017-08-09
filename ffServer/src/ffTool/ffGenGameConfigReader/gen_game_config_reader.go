package main

import (
	"ffCommon/util"
	"fmt"
)

var fmtFile = `package main

import (
	"ffAutoGen/ffGameConfig"
	"ffCommon/log/log"
	"ffCommon/util"
)

func main() {
	defer util.PanicProtect()

	var err error
%v
}
`

var fmtOne = `
	_, err = ffGameConfig.Read%v()
	if err != nil {
		log.RunLogger.Printf("Read%v get error[%%v]", err)
	}
`

func genGameConfigReader(ffGameConfigReaderFile string, golangFiles []string) {
	readAllFile := ""
	for _, filename := range golangFiles {
		readAllFile += fmt.Sprintf(fmtOne, filename, filename)
	}

	result := fmt.Sprintf(fmtFile, readAllFile)

	util.WriteFile(ffGameConfigReaderFile, []byte(result))
}
