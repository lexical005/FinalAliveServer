package cellvalue

import (
	"encoding/json"
	"fmt"
	"strings"
)

type valueStoreStringArray struct {
	*valueStore

	value []string
}

func (vs *valueStoreStringArray) Store(data string) error {
	if vs.value == nil {
		vs.value = make([]string, 0, 1)
	}

	if strings.HasPrefix(data, "[") && strings.HasSuffix(data, "]") {
		var dataOri []interface{}
		if err := json.Unmarshal([]byte(data), &dataOri); err != nil {
			return fmt.Errorf("ValueStore[%v] Invalid string array data[%v]", vs.GoType(), data)
		}

		for _, s := range dataOri {
			str, _ := s.(string)

			vs.value = append(vs.value, str)
		}

	} else {
		vs.value = append(vs.value, data)
	}

	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtStringArray] = func(vt ValueType) ValueStore {
		return &valueStoreStringArray{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
