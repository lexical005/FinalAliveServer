package ffGameConfig

import (
	"ffCommon/util"

	"fmt"

	"github.com/lexical005/toml"
)

// Asset excel Asset
type Asset struct {
	Actor []*Actor
}

func (a *Asset) String() string {
	result := ""
	result += "Actor"
	for _, row := range a.Actor {
		result += fmt.Sprintf("%v\n", row)
	}

	return result
}

// Name the toml config's name
func (a *Asset) Name() string {
	return "Asset"
}

// Actor sheet Actor of excel Asset
type Actor struct {
	TemplateID int32
}

func (a *Actor) String() string {
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
