package cellvalue

import (
	"fmt"
	"strconv"
)

type valueStoreIntsMulti struct {
	*valueStore

	value []int
}

func (vs *valueStoreIntsMulti) Store(data string) error {
	i64, err := strconv.ParseInt(data, 10, 0)
	if err != nil {
		return fmt.Errorf("ValueStore[%v] Invalid int data[%v]", vs.Type(), data)
	}

	vs.value = append(vs.value, int(i64))
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtIntsMulti] = func(vt ValueType) ValueStore {
		return &valueStoreIntsMulti{
			valueStore: &valueStore{
				vt: vt,
			},
			value: make([]int, 0, 1),
		}
	}
}
