package ffExcel

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"os/exec"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/lexical005/toml"
)

var fmtGoEnumFile = `package ffEnum

import (
	"fmt"
)
{AllEnumType}
`

var fmtGoInternalEnum = `
	internal{EnumType}{EnumKey}    {EnumType} = {EnumType}({EnumValue}) // {EnumDesc}`

var fmtGoInternalEnumInfo = `
	&internal{EnumType}Info{
		value: internal{EnumType}{EnumKey},
		toml:  "{EnumType}.{EnumKey}",
		desc:  "{EnumDesc}",
	},`

var fmtGoMapInternalEnumInfo = `
	all{EnumType}Info[int(internal{EnumType}{EnumKey})].toml: all{EnumType}Info[int(internal{EnumType}{EnumKey})],`

var fmtGoOneEnum = `
// {EnumType} {EnumType}
type {EnumType} int32

const ({AllInternalEnum}
)

type internal{EnumType}Info struct {
	value {EnumType}
	toml  string
	desc  string
}

var all{EnumType}Info = []*internal{EnumType}Info{{AllInternalEnumInfo}
}

var mapCodeTo{EnumType}Info = map[string]*internal{EnumType}Info{{AllMapInternalEnumInfo}
}

// UnmarshalText implements encoding.TextUnmarshaler
func (e *{EnumType}) UnmarshalText(data []byte) error {
	key := string(data)
	v, ok := mapCodeTo{EnumType}Info[key]
	if !ok {
		return fmt.Errorf("{EnumType}.UnmarshalText failed: invalid {EnumType}[%v]", key)
	}
	*e = v.value
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (e {EnumType}) MarshalText() ([]byte, error) {
	return []byte(all{EnumType}Info[e].toml), nil
}

func (e {EnumType}) String() string {
	return all{EnumType}Info[e].toml
}`

type enum struct {
	Enum string
	Desc string
	Name string
}
type fileEnum struct {
	Enum []*enum
}

func genGoEnum(dataFilePath string, fileEnum *fileEnum) {
	mapEnums := make(map[string][]*enum, len(fileEnum.Enum))
	for _, one := range fileEnum.Enum {
		if listEnums, ok := mapEnums[one.Enum]; ok {
			listEnums = append(listEnums, one)
			mapEnums[one.Enum] = listEnums
		} else {
			listEnums := make([]*enum, 0, 4)
			listEnums = append(listEnums, one)
			mapEnums[one.Enum] = listEnums
		}
	}

	allEnumType := make([]string, 0, 4)
	for enumType := range mapEnums {
		allEnumType = append(allEnumType, enumType)
	}
	sort.Strings(allEnumType)

	AllEnumType := ""
	for _, enumType := range allEnumType {
		listEnums := mapEnums[enumType]

		AllInternalEnum := ""
		for i, one := range listEnums {
			s := strings.Replace(fmtGoInternalEnum, "{EnumType}", enumType, -1)
			s = strings.Replace(s, "{EnumKey}", one.Name, -1)
			s = strings.Replace(s, "{EnumValue}", strconv.Itoa(i), -1)
			s = strings.Replace(s, "{EnumDesc}", one.Desc, -1)
			AllInternalEnum += s
		}

		AllInternalEnumInfo := ""
		for _, one := range listEnums {
			s := strings.Replace(fmtGoInternalEnumInfo, "{EnumType}", enumType, -1)
			s = strings.Replace(s, "{EnumKey}", one.Name, -1)
			s = strings.Replace(s, "{EnumDesc}", one.Desc, -1)
			AllInternalEnumInfo += s
		}

		AllMapInternalEnumInfo := ""
		for _, one := range listEnums {
			s := strings.Replace(fmtGoMapInternalEnumInfo, "{EnumType}", enumType, -1)
			s = strings.Replace(s, "{EnumKey}", one.Name, -1)
			s = strings.Replace(s, "{EnumDesc}", one.Desc, -1)
			AllMapInternalEnumInfo += s
		}

		s := strings.Replace(fmtGoOneEnum, "{EnumType}", enumType, -1)
		s = strings.Replace(s, "{AllInternalEnum}", AllInternalEnum, -1)
		s = strings.Replace(s, "{AllInternalEnumInfo}", AllInternalEnumInfo, -1)
		s = strings.Replace(s, "{AllMapInternalEnumInfo}", AllMapInternalEnumInfo, -1)
		AllEnumType += s
	}

	result := strings.Replace(fmtGoEnumFile, "{AllEnumType}", AllEnumType, -1)

	util.WriteFile(dataFilePath, []byte(result))
	log.RunLogger.Println(dataFilePath)
	exec.Command("go", "fmt", "ffAutoGen/ffEnum").Output()
}

var fmtCSharpFile = `namespace NConfig
{{AllEnum}
}
`

var fmtCSharpEnum = `
	// {EnumType}
	public static class {EnumType}
	{{AllEnumKey}
	}
`

var fmtCSharpEnumKey = `
		// {EnumDesc}
		public const int {EnumKey} = {EnumValue};`

func genCSharpEnum(dataFilePath string, fileEnum *fileEnum) {
	mapEnums := make(map[string][]*enum, len(fileEnum.Enum))
	for _, one := range fileEnum.Enum {
		if listEnums, ok := mapEnums[one.Enum]; ok {
			listEnums = append(listEnums, one)
			mapEnums[one.Enum] = listEnums
		} else {
			listEnums := make([]*enum, 0, 4)
			listEnums = append(listEnums, one)
			mapEnums[one.Enum] = listEnums
		}
	}

	allEnumType := make([]string, 0, 4)
	for enumType := range mapEnums {
		allEnumType = append(allEnumType, enumType)
	}
	sort.Strings(allEnumType)

	AllEnum := ""
	for _, enumType := range allEnumType {
		listEnums := mapEnums[enumType]

		AllEnumKey := ""
		for i, one := range listEnums {
			s := strings.Replace(fmtCSharpEnumKey, "{EnumKey}", one.Name, -1)
			s = strings.Replace(s, "{EnumValue}", strconv.Itoa(i), -1)
			s = strings.Replace(s, "{EnumDesc}", one.Desc, -1)
			AllEnumKey += s
		}

		s := strings.Replace(fmtCSharpEnum, "{EnumType}", enumType, -1)
		s = strings.Replace(s, "{AllEnumKey}", AllEnumKey, -1)
		AllEnum += s
	}

	result := strings.Replace(fmtCSharpFile, "{AllEnum}", AllEnum, -1)

	util.WriteFile(dataFilePath, []byte(result))
	log.RunLogger.Println(dataFilePath)
}

func genEnum(excel *excel) {
	// Server
	if excel.exportServerGoCodePath != "" && exportConfig.hasGoEnv {
		tomlData, err := util.ReadFile(path.Join("toml", "server", "Enum.toml"))
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		// 解析导出的toml配置文件
		fileEnum := &fileEnum{}
		err = toml.Unmarshal(tomlData, fileEnum)
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		genGoEnum(path.Join(excel.exportServerGoCodePath, "Enum.go"), fileEnum)
	}

	// Client
	if excel.exportClientCSharpCodePath != "" {
		tomlData, err := util.ReadFile(path.Join("toml", "client", "Enum.toml"))
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		// 解析导出的toml配置文件
		fileEnum := &fileEnum{}
		err = toml.Unmarshal(tomlData, fileEnum)
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		genCSharpEnum(path.Join(excel.exportClientCSharpCodePath, "Enum.cs"), fileEnum)
	}
}
