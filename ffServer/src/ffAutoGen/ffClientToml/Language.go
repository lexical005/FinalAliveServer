package ffClientToml

import (
	"ffCommon/util"

	"fmt"

	"github.com/lexical005/toml"
)

// Language excel Language
type Language struct {
	Common  map[string]*Common
	Special []*Special
	Error   map[string]*Error
}

func (l *Language) String() string {
	result := ""
	result += "Common"
	for k, v := range l.Common {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	result += "Special"
	for _, row := range l.Special {
		result += fmt.Sprintf("%v\n", row)
	}

	result += "Error"
	for k, v := range l.Error {
		result += fmt.Sprintf("%v:%v\n", k, v)
	}

	return result
}

// Name the toml config's name
func (l *Language) Name() string {
	return "Language"
}

// Common sheet Common of excel Language
type Common struct {
	Value string
}

func (c *Common) String() string {
	result := "["
	result += fmt.Sprintf("Value:%v,", c.Value)
	result += "]"
	return result
}

// Special sheet Special of excel Language
type Special struct {
	Value string
}

func (s *Special) String() string {
	result := "["
	result += fmt.Sprintf("Value:%v,", s.Value)
	result += "]"
	return result
}

// Error sheet Error of excel Language
type Error struct {
	Value string
}

func (e *Error) String() string {
	result := "["
	result += fmt.Sprintf("Value:%v,", e.Value)
	result += "]"
	return result
}

// ReadLanguage read excel Language
func ReadLanguage() (l *Language, err error) {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/client/Language.toml")
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
