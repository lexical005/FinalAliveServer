package ffGameConfig

import (
	"ffCommon/util"

	"fmt"

	"github.com/lexical005/toml"
)

// Language excel Language
type Language struct {
	AIName []AIName
}

func (l *Language) String() string {
	result := ""
	result += "AIName"
	for _, row := range l.AIName {
		result += fmt.Sprintf("%v\n", row)
	}

	return result
}

// Name the toml config's name
func (l *Language) Name() string {
	return "Language"
}

// AIName sheet AIName of excel Language
type AIName struct {
	Value string
}

func (ain *AIName) String() string {
	result := "["
	result += fmt.Sprintf("Value:%v,", ain.Value)
	result += "]"
	return result
}

// ReadLanguage read excel Language
func ReadLanguage() (l *Language, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/game/Language.toml")
	if err != nil {
		return
	}

	// 解析
	l = &Language{}
	err = toml.Unmarshal(fileContent, l)
	if err != nil {
		return
	}

	return
}
