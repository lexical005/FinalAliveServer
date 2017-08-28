package ffClientToml

import (
	"ffCommon/util"

	"fmt"

	"github.com/lexical005/toml"
)

// Asset excel Asset
type Asset struct {
	Assets map[int32]*Assets
}

func (a *Asset) String() string {
	result := ""
	result += "Assets"
	for k, v := range a.Assets {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	return result
}

// Name the toml config's name
func (a *Asset) Name() string {
	return "Asset"
}

// Assets sheet Assets of excel Asset
type Assets struct {
	BattleDefault string
	HomeDefault   string
}

func (a *Assets) String() string {
	result := "["
	result += fmt.Sprintf("BattleDefault:%v,", a.BattleDefault)
	result += fmt.Sprintf("HomeDefault:%v,", a.HomeDefault)
	result += "]"
	return result
}

// ReadAsset read excel Asset
func ReadAsset() (a *Asset, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/client/Asset.toml")
	if err != nil {
		return
	}

	// 解析
	a = &Asset{}
	err = toml.Unmarshal(fileContent, a)
	if err != nil {
		return
	}

	return
}
