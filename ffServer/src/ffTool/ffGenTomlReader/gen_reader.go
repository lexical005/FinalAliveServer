package main

import (
	"ffCommon/util"
	"fmt"
)

var fmtFile = `package main

import (
	%v
	"ffCommon/log/log"
	"ffCommon/util"
)

%v
func read() {
	defer util.PanicProtect()

	var err error
%v
}
`

var fmtPackage = `"ffAutoGen/%v"`

var fmtVarFile = `var toml%v *%v.%v
`

var fmtReadFile = `
	toml%v, err = %v.Read%v()
	if err != nil {
		log.RunLogger.Printf("Read%v get error[%%v]", err)
	}
`

func genReadTomlCode(fullpath string, golangFiles []string, packageName string) {
	packageImport := fmt.Sprintf(fmtPackage, packageName)

	allVar, readAllFile := "", ""
	for _, filename := range golangFiles {
		allVar += fmt.Sprintf(fmtVarFile, filename, packageName, filename)
		readAllFile += fmt.Sprintf(fmtReadFile, filename, packageName, filename, filename)
	}

	result := fmt.Sprintf(fmtFile, packageImport, allVar, readAllFile)

	util.WriteFile(fullpath, []byte(result))
}
