package cellvalue

import (
	"ffLogic/ffGrammar"
	"fmt"
)

type valueStoreGrammar struct {
	vt ValueType

	grammar *ffGrammar.Grammar

	value string
}

func (vs *valueStoreGrammar) Store(data string) error {
	grammar := &ffGrammar.Grammar{}
	err := grammar.UnmarshalTOML([]byte(data))
	if err != nil {
		return err
	}

	vs.grammar = grammar
	vs.value = data
	return nil
}

func (vs *valueStoreGrammar) Type() string {
	return vs.vt.Type()
}
func (vs *valueStoreGrammar) Value() interface{} {
	return vs.value
}
func (vs *valueStoreGrammar) ValueToml() string {
	return fmt.Sprintf(`"%v"`, vs.value)
}

func (vs *valueStoreGrammar) String() string {
	return fmt.Sprintf("[%v:%v]", vs.Type(), vs.value)
}

func init() {
	mapValueStoreCreator[vtGrammar] = func(vt ValueType) ValueStore {
		return &valueStoreGrammar{
			vt: vt,
		}
	}
}
