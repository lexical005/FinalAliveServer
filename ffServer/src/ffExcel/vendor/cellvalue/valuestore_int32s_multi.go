package cellvalue

import (
	"fmt"
	"strconv"
)

type valueStoreInt32sMulti struct {
	*valueStore

	value []int32
}

func (vs *valueStoreInt32sMulti) Store(data string) error {
	i64, err := strconv.ParseInt(data, 10, 0)
	if err != nil {
		return fmt.Errorf("ValueStore[%v] not valid number data[%v]", vs.Type(), data)
	}

	if i64 < -2147483648 || i64 > 2147483647 {
		return fmt.Errorf("ValueStore[%v] number outof int32 range data[%v]", vs.Type(), data)
	}

	vs.value = append(vs.value, int32(i64))
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtInt32sMulti] = func(vt ValueType) ValueStore {
		return &valueStoreInt32sMulti{
			valueStore: &valueStore{
				vt: vt,
			},
			value: make([]int32, 0, 1),
		}
	}
}
