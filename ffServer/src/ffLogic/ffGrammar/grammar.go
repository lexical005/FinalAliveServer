package ffGrammar

import (
	"ffAutoGen/ffError"
	"ffLogic/ffDef"
	"impGrammar"
)

// Grammar grammar的具体实现的封装
type Grammar struct {
	grammar impGrammar.IGrammar
}

// UnmarshalTOML toml调用此接口生成Grammar实例, 同时也供ffCommon/excel/vendor/cellvalue解析excel使用
func (g *Grammar) UnmarshalTOML(data []byte) error {
	grammar, err := impGrammar.Parse(string(data))
	if err != nil {
		return err
	}

	g.grammar = grammar
	return nil
}

// Check 检查
func (g *Grammar) Check(account ffDef.IAccount) ffError.Error {
	return g.grammar.Check(account)
}

// Add 添加
func (g *Grammar) Add(account ffDef.IAccount) ffError.Error {
	return g.grammar.Add(account)
}

// Sub 扣除
func (g *Grammar) Sub(account ffDef.IAccount) ffError.Error {
	return g.grammar.Sub(account)
}

// Excute 执行
func (g *Grammar) Excute(account ffDef.IAccount) ffError.Error {
	return g.grammar.Excute(account)
}
