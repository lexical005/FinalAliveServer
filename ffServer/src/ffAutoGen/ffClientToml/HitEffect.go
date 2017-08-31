package ffClientToml

import (
	"ffCommon/util"

	"fmt"

	"github.com/lexical005/toml"
)

// HitEffect excel HitEffect
type HitEffect struct {
	Hit map[string]*Hit
}

func (h *HitEffect) String() string {
	result := ""
	result += "Hit"
	for k, v := range h.Hit {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	return result
}

// Name the toml config's name
func (h *HitEffect) Name() string {
	return "HitEffect"
}

// Hit sheet Hit of excel HitEffect
type Hit struct {
	Gun int32
}

func (h *Hit) String() string {
	result := "["
	result += fmt.Sprintf("Gun:%v,", h.Gun)
	result += "]"
	return result
}

// ReadHitEffect read excel HitEffect
func ReadHitEffect() (h *HitEffect, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/client/HitEffect.toml")
	if err != nil {
		return
	}

	// 解析
	h = &HitEffect{}
	err = toml.Unmarshal(fileContent, h)
	if err != nil {
		return
	}

	return
}
