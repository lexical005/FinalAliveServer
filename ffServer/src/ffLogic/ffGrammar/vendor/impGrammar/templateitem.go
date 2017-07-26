package impGrammar

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffLogic/ffDef"
	"fmt"
	"strconv"
)

// 模板物品
type grammarTemplateItem struct {
	id    int
	count int
}

func (g *grammarTemplateItem) Parse(params []string) (err error) {
	// 参数数目
	if len(params) != 2 {
		return fmt.Errorf("impGrammar.grammarTemplateItem.Parse: invalid params count params[%v]", params)
	}

	// 物品模板
	g.id, err = strconv.Atoi(params[0])
	if err != nil || g.id < 1 {
		return fmt.Errorf("impGrammar.grammarTemplateItem.Parse: invalid params item templateid params[%v]", params)
	}

	// 物品数量
	g.count, err = strconv.Atoi(params[1])
	if err != nil || g.count < 1 {
		return fmt.Errorf("impGrammar.grammarTemplateItem.Parse: invalid params item count params[%v]", params)
	}

	return nil
}

// Check 检查
func (g *grammarTemplateItem) Check(account ffDef.IAccount) ffError.Error {
	if g.count <= account.ItemMgr().TemplateCount(g.id) {
		return ffError.ErrNone
	}
	return ffError.ErrTemplateItemLess
}

// Add 添加
func (g *grammarTemplateItem) Add(account ffDef.IAccount) ffError.Error {
	return account.ItemMgr().AddTemplate(g.id, g.count)
}

// Sub 扣除
func (g *grammarTemplateItem) Sub(account ffDef.IAccount) ffError.Error {
	result := g.Check(account)
	if result != ffError.ErrNone {
		return result
	}
	return account.ItemMgr().SubTemplate(g.id, g.count)
}

// Excute 执行
func (g *grammarTemplateItem) Excute(account ffDef.IAccount) ffError.Error {
	log.FatalLogger.Println("grammarTemplateItem.Excute: should not be called")
	return ffError.ErrNone
}

func init() {
	mapImpGrammarCreator["item"] = func(params []string) (IGrammar, error) {
		g := &grammarTemplateItem{}
		if err := g.Parse(params); err != nil {
			return nil, err
		}
		return g, nil
	}
}
