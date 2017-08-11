package main

import (
	"ffCommon/util"
	"fmt"
)

var fmtTomlFile = `package main

import (
	%v
	"ffCommon/log/log"
)

%v
func readToml() {
	var err error
%v
}

func init() {
	allRead = append(allRead, readToml)
}
`

var fmtTomlPackage = `"ffAutoGen/%v"`

var fmtTomlVarFile = `var toml%v *%v.%v
`

var fmtTomlReadFile = `
	toml%v, err = %v.Read%v()
	if err != nil {
		log.RunLogger.Printf("Read%v get error[%%v]", err)
	}
`

func genReadTomlCode(fullpath string, golangFiles []string, goCodePackageName string) {
	packageImport := fmt.Sprintf(fmtTomlPackage, goCodePackageName)

	allVar, readAllFile := "", ""
	for _, filename := range golangFiles {
		allVar += fmt.Sprintf(fmtTomlVarFile, filename, goCodePackageName, filename)
		readAllFile += fmt.Sprintf(fmtTomlReadFile, filename, goCodePackageName, filename, filename)
	}

	result := fmt.Sprintf(fmtTomlFile, packageImport, allVar, readAllFile)

	util.WriteFile(fullpath, []byte(result))
}
