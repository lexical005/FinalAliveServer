package cellvalue

import (
	"encoding/json"
	"fmt"
)

type valueStoreStringsSingle struct {
	*valueStore

	value []string
}

func (vs *valueStoreStringsSingle) Store(data string) error {
	var dataOri []interface{}
	if err := json.Unmarshal([]byte(data), &dataOri); err != nil {
		return fmt.Errorf("ValueStore[%v] Invalid string array data[%v]", vs.Type(), data)
	}

	value := make([]string, len(dataOri), len(dataOri))
	for i, s := range dataOri {
		str, _ := s.(string)

		value[i] = str
	}

	vs.value = value
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtStringsSingle] = func(vt ValueType) ValueStore {
		return &valueStoreStringsSingle{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
