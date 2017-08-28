package ffGameConfig

import (
	"ffCommon/util"

	"fmt"

	"github.com/lexical005/toml"
)

// Asset excel Asset
type Asset struct {
	Assets []*Assets
}

func (a *Asset) String() string {
	result := ""
	result += "Assets"
	for _, row := range a.Assets {
		result += fmt.Sprintf("%v\n", row)
	}

	return result
}

// Name the toml config's name
func (a *Asset) Name() string {
	return "Asset"
}

// Assets sheet Assets of excel Asset
type Assets struct {
	TemplateID int32
}

func (a *Assets) String() string {
	result := "["
	result += fmt.Sprintf("TemplateID:%v,", a.TemplateID)
	result += "]"
	return result
}

// ReadAsset read excel Asset
func ReadAsset() (a *Asset, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/game/Asset.toml")
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
