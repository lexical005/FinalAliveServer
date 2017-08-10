package cellvalue

import (
	"fmt"
	"strconv"
)

type valueStoreInt32 struct {
	*valueStore

	value int32
}

func (vs *valueStoreInt32) Store(data string) error {
	if len(data) == 0 {
		vs.value = 0
		vs.valueStore.value = vs.value
		return nil
	}

	i64, err := strconv.ParseInt(data, 10, 0)
	if err != nil {
		return fmt.Errorf("ValueStore[%v] not valid number data[%v]", vs.Type(), data)
	}

	if i64 < -2147483648 || i64 > 2147483647 {
		return fmt.Errorf("ValueStore[%v] number outof int32 range data[%v]", vs.Type(), data)
	}

	vs.value = int32(i64)
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtInt32] = func(vt ValueType) ValueStore {
		return &valueStoreInt32{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
