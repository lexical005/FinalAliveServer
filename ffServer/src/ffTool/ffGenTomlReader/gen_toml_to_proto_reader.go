package main

import (
	"ffCommon/util"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var fmtTransPackage = `package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"path/filepath"

    proto "github.com/golang/protobuf/proto"
)
`

var fmtTransInit = `
func init() {
    allTrans = append(allTrans, trans{FileName})
}
`

var fmtTransFuncMain = `
func trans{FileName}() {
    message := &{FileName}{}
%v
    pbBuf := proto.NewBuffer(make([]byte, 0, 1024*10))
    if err := pbBuf.Marshal(message); err != nil {
        log.RunLogger.Printf("trans{FileName} err[%%v]", err)
        return
    }

    util.WriteFile(filepath.Join("ProtoBuf", "Client", "bytes", toml{FileName}.Name()+".bytes"), pbBuf.Bytes())
}
`

var fmtTransStructMap = `
    // {StructName}
	message.{StructName} = make(map[int32]*{FileName}_St{StructName}, len(toml{FileName}.{StructName}))
	for k, v := range toml{FileName}.{StructName} {
		message.{StructName}[k] = &{FileName}_St{StructName}{%v
		}
	}
`

var fmtTransStructStruct = `
    // {StructName}
	message.{StructName} = &{FileName}_St{StructName}{%v
	}
`

var fmtTransStructList = `
    // {StructName}
	message.{StructName} = make([]*{FileName}_St{StructName}, len(toml{FileName}.{StructName}))
	for k, v := range toml{FileName}.{StructName} {
		message.{StructName}[k] = &{FileName}_St{StructName}{%v
		}
	}
`

var fmtTransMemberNormalMap = "\n\t\t\t%v: v.%v,"
var fmtTransMemberNormalList = "\n\t\t\t%v: v.%v,"
var fmtTransMemberNormalStruct = "\n\t\t\t%v: toml{FileName}.{StructName}.%v,"
var fmtTransMemberGrammarMap = "\n\t\t\t%v: transGrammar(v.%v),"
var fmtTransMemberGrammarList = "\n\t\t\t%v: transGrammar(v.%v),"
var fmtTransMemberGrammarStruct = "\n\t\t\t%v: transGrammar(toml{FileName}.{StructName}.%v),"

// 正则表达式说明
// http://www.cnblogs.com/golove/p/3269099.html
var regexpStruct = regexp.MustCompile(`type\s+([\w]+)\s+struct\s+{\n(?s)(.+?)\}`)
var regexpStructVar = regexp.MustCompile(`\s*([\w]+)\s+([\w\[\]\.\*]+)`)

// 结构体定义
type structDef struct {
	name      string   // 结构体自身的定义
	vars      []string // 结构体的成员变量的名称
	types     []string // 结构体的成员变量的类型
	lowerVars []string // 结构体的成员变量的小写名称
}

// 文件内的所有结构体定义
type fileStructDef struct {
	name string
	defs []*structDef
}

// 解析文件内的结构体定义
func getFileDef(content string, filename string) *fileStructDef {
	// 捕获所有的结构体定义
	result1 := regexpStruct.FindAllStringSubmatch(content, -1)
	fileDef := &fileStructDef{
		name: filename,
		defs: make([]*structDef, len(result1), len(result1)),
	}

	for i, one := range result1 {
		// 结构体名称
		nameStruct := one[1]
		// 结构体成员
		allStructVars := one[2]

		// 捕获成员定义
		result2 := regexpStructVar.FindAllStringSubmatch(allStructVars, -1)

		structDef := &structDef{
			name:      nameStruct,
			vars:      make([]string, len(result2), len(result2)),
			lowerVars: make([]string, len(result2), len(result2)),
			types:     make([]string, len(result2), len(result2)),
		}

		//
		for j, two := range result2 {
			varName := two[1]
			varType := two[2]
			structDef.vars[j] = varName
			structDef.lowerVars[j] = strings.ToLower(varName)
			structDef.types[j] = varType
		}

		fileDef.defs[i] = structDef
	}

	return fileDef
}

// 读取文件
func fileContent(fullpath string) (filename string, content string) {
	// 获取文件名, 不含扩展
	filename = filepath.Base(fullpath)
	fileext := filepath.Ext(fullpath)
	filename = filename[0 : len(filename)-len(fileext)]

	// 读取文件内容
	bytes, _ := util.ReadFile(fullpath)
	return filename, string(bytes)
}

// 生成转换代码
func genTransCode(saveFullDir string, protoFileDef, tomlFileDef *fileStructDef) {

	getStructDef := func(fileDef *fileStructDef, key string) *structDef {
		for _, structDef := range fileDef.defs {
			if key == strings.ToLower(structDef.name) {
				return structDef
			}
		}
		return nil
	}

	tomlMainStructDef := getStructDef(tomlFileDef, strings.ToLower(tomlFileDef.name))

	result := ""

	mainContent := strings.Replace(fmtTransFuncMain, "{FileName}", tomlFileDef.name, -1)

	allStructs := ""
	for _, tomlDef := range tomlFileDef.defs {
		if strings.ToLower(tomlDef.name) == strings.ToLower(tomlFileDef.name) {
			continue
		}

		mainStructVarType := "struct"
		for j, name := range tomlMainStructDef.vars {
			if name == tomlDef.name {
				if strings.HasPrefix(tomlMainStructDef.types[j], "map[") {
					mainStructVarType = "map"
				} else if strings.HasPrefix(tomlMainStructDef.types[j], "[]") {
					mainStructVarType = "list"
				}
				break
			}
		}

		protoStructDef := getStructDef(protoFileDef, strings.ToLower(tomlFileDef.name)+"_st"+strings.ToLower(tomlDef.name))

		members := ""
		for j := 0; j < len(tomlDef.vars); j++ {
			if tomlDef.types[j] != "ffGrammar.Grammar" {
				if mainStructVarType == "map" {
					members += fmt.Sprintf(fmtTransMemberNormalMap, protoStructDef.vars[j], tomlDef.vars[j])
				} else if mainStructVarType == "list" {
					members += fmt.Sprintf(fmtTransMemberNormalList, protoStructDef.vars[j], tomlDef.vars[j])
				} else {
					t := strings.Replace(fmtTransMemberNormalStruct, "{FileName}", tomlFileDef.name, -1)
					t = strings.Replace(t, "{StructName}", tomlDef.name, -1)
					members += fmt.Sprintf(t, protoStructDef.vars[j], tomlDef.vars[j])
				}
			} else {
				if mainStructVarType == "map" {
					members += fmt.Sprintf(fmtTransMemberGrammarMap, protoStructDef.vars[j], tomlDef.vars[j])
				} else if mainStructVarType == "list" {
					members += fmt.Sprintf(fmtTransMemberGrammarList, protoStructDef.vars[j], tomlDef.vars[j])
				} else {
					t := strings.Replace(fmtTransMemberGrammarStruct, "{FileName}", tomlFileDef.name, -1)
					t = strings.Replace(t, "{StructName}", tomlDef.name, -1)
					members += fmt.Sprintf(t, protoStructDef.vars[j], tomlDef.vars[j])
				}
			}
		}

		var structs string
		if mainStructVarType == "map" {
			structs = strings.Replace(fmtTransStructMap, "{FileName}", tomlFileDef.name, -1)
		} else if mainStructVarType == "list" {
			structs = strings.Replace(fmtTransStructList, "{FileName}", tomlFileDef.name, -1)
		} else {
			structs = strings.Replace(fmtTransStructStruct, "{FileName}", tomlFileDef.name, -1)
		}
		structs = strings.Replace(structs, "{StructName}", tomlDef.name, -1)
		structs = fmt.Sprintf(structs, members)

		allStructs += structs
	}
	mainContent = fmt.Sprintf(mainContent, allStructs)

	result += fmtTransPackage
	result += mainContent
	result += strings.Replace(fmtTransInit, "{FileName}", tomlFileDef.name, -1)

	util.WriteFile(filepath.Join(saveFullDir, "trans"+tomlFileDef.name+".go"), []byte(result))
}

// 转换
func transGoToProto(saveFullDir string, protoFilePath string, goFullPathFiles []string, packageName string) {
	// Proto的go代码
	var protoFileDef *fileStructDef
	{
		filename, content := fileContent(protoFilePath)
		protoFileDef = getFileDef(string(content), filename)

		fmt.Printf("%v:\n", filename)
		for _, v := range protoFileDef.defs {
			fmt.Printf("%v:%v\n%q\n%q\n\n", v.name, len(v.vars), v.vars, v.types)
		}
		fmt.Printf("\n\n")
	}

	// 读取toml的go代码
	for _, fullpath := range goFullPathFiles {
		filename, content := fileContent(fullpath)
		tomlFileDef := getFileDef(string(content), filename)

		fmt.Printf("%v:\n", filename)
		for _, v := range tomlFileDef.defs {
			fmt.Printf("%v:%v\n%q\n%q\n\n", v.name, len(v.vars), v.vars, v.types)
		}
		fmt.Printf("\n\n")

		genTransCode(saveFullDir, protoFileDef, tomlFileDef)
	}
}
