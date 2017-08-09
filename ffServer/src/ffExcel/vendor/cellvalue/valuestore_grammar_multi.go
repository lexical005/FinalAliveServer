package cellvalue

import (
	"ffLogic/ffGrammar"
	"fmt"
)

type valueStoreGrammarMulti struct {
	vt ValueType

	grammar *ffGrammar.Grammar

	value string
}

func (vs *valueStoreGrammarMulti) Store(data string) error {
	grammar := &ffGrammar.Grammar{}
	err := grammar.UnmarshalTOML([]byte(data))
	if err != nil {
		return err
	}

	vs.grammar = grammar
	vs.value = data
	return nil
}

func (vs *valueStoreGrammarMulti) Type() string {
	return vs.vt.Type()
}
func (vs *valueStoreGrammarMulti) Value() interface{} {
	return vs.value
}
func (vs *valueStoreGrammarMulti) ValueToml() string {
	return fmt.Sprintf(`"%v"`, vs.value)
}

func (vs *valueStoreGrammarMulti) String() string {
	return fmt.Sprintf("[%v:%v]", vs.Type(), vs.value)
}

func init() {
	mapValueStoreCreator[vtGrammarsMulti] = func(vt ValueType) ValueStore {
		return &valueStoreGrammarMulti{
			vt: vt,
		}
	}
}
