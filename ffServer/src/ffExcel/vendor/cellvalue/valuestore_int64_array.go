package cellvalue

import (
	"fmt"
	"regexp"
	"strconv"
)

var regexpInt64Array = regexp.MustCompile(`[\d]+`)

type valueStoreInt64Array struct {
	*valueStore

	value []int
}

func (vs *valueStoreInt64Array) Store(data string) error {
	result := regexpInt64Array.FindAllString(data, -1)

	value := make([]int, len(result), len(result))
	for i, s := range result {
		i64, err := strconv.ParseInt(s, 10, 0)
		if err != nil {
			return fmt.Errorf("ValueStore[%v] Invalid int array data[%v]", vs.Type(), data)
		}
		value[i] = int(i64)
	}

	vs.value = value
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtInt64Array] = func(vt ValueType) ValueStore {
		return &valueStoreInt64Array{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
