package cellvalue

import (
	"fmt"
	"strconv"
)

type valueStoreInt64 struct {
	*valueStore

	value int
}

func (vs *valueStoreInt64) Store(data string) error {
	i64, err := strconv.ParseInt(data, 10, 0)
	if err != nil {
		return fmt.Errorf("ValueStore[%v] Invalid int data[%v]", vs.Type(), data)
	}

	vs.value = int(i64)
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtInt64] = func(vt ValueType) ValueStore {
		return &valueStoreInt64{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
