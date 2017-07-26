package impGrammar

import (
	"ffAutoGen/ffError"
	"ffLogic/ffDef"
	"fmt"
	"strings"
)

// IGrammar grammar的具体实现，需要实现的接口
type IGrammar interface {
	// Parse 根据字符串参数，解析成程序使用的grammar
	Parse(params []string) error

	// Check 检查
	Check(account ffDef.IAccount) ffError.Error

	// Add 添加
	Add(account ffDef.IAccount) ffError.Error

	// Sub 扣除
	Sub(account ffDef.IAccount) ffError.Error

	// Excute 执行
	Excute(account ffDef.IAccount) ffError.Error
}

// grammar具体实现的创建映射
var mapImpGrammarCreator = map[string]func(params []string) (IGrammar, error){}

// Parse 将字符串grammar解析为具体实现
func Parse(strGrammar string) (g IGrammar, err error) {
	if len(strGrammar) < 1 {
		return nil, fmt.Errorf("impGrammar.Parse: invalid strGrammar[%v]", strGrammar)
	}

	// 先移除首尾的空白字符，再移除首尾的"
	strGrammar = strings.TrimSpace(strGrammar)
	if strGrammar[0] == '"' {
		strGrammar = strGrammar[1:]
	}
	if strGrammar[len(strGrammar)-1] == '"' {
		strGrammar = strGrammar[:len(strGrammar)-1]
	}

	result := strings.Split(strGrammar, ",")

	// 无效的字符串
	if len(result) < 2 {
		return nil, fmt.Errorf("impGrammar.Parse: invalid strGrammar[%v]", strGrammar)
	}

	// 无效的Key
	grammarKey := result[0]
	creator, ok := mapImpGrammarCreator[grammarKey]
	if !ok {
		return nil, fmt.Errorf("impGrammar.Parse: invalid grammarKey[%v][%v]", grammarKey, strGrammar)
	}

	return creator(result[1:])
}
