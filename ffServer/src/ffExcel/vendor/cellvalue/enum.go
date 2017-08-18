package cellvalue

import (
	"ffCommon/log/log"
	"ffCommon/util"

	"github.com/lexical005/toml"
)

type enum struct {
	Enum string
	Desc string
	Name string
}

type fileEnum struct {
	Enum []*enum
}

// 枚举定义
//	EItemType: [GunWeapon,Ammunition,Attachment,MelleeWeapon,Equipment,Consumable,Throwable]
var mapEnums map[string][]string = make(map[string][]string, 16)

func initEnum(tomlFile string) {
	tomlData, err := util.ReadFile(tomlFile)
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

	for _, one := range fileEnum.Enum {
		if listEnums, ok := mapEnums[one.Enum]; ok {
			listEnums = append(listEnums, one.Name)
			mapEnums[one.Enum] = listEnums
		} else {
			listEnums := make([]string, 0, 4)
			listEnums = append(listEnums, one.Name)
			mapEnums[one.Enum] = listEnums
		}
	}
}
