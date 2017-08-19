package cellvalue

import (
	"strings"
)

type valueStoreStringArray struct {
	*valueStore

	value []string
}

func (vs *valueStoreStringArray) Store(data string, vt ValueType) error {
	if vs.value == nil {
		vs.value = make([]string, 0, 1)
	}

	if strings.HasPrefix(data, "[") && strings.HasSuffix(data, "]") {
		data = data[1 : len(data)-1]
		tmp := strings.Split(data, ",")
		for _, one := range tmp {
			vs.value = append(vs.value, one)
		}
	} else {
		vs.value = append(vs.value, data)
	}

	vs.valueStore.value = vs.value
	return nil
}

func init() {
	basicValueStoreCreator[vtStringArray] = func(vt ValueType) ValueStore {
		return &valueStoreStringArray{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
