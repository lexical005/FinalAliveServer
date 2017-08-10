package cellvalue

import (
	"fmt"
	"regexp"
	"strconv"
)

var regexpInt64sSingle = regexp.MustCompile(`[\d]+`)

type valueStoreInt64sSingle struct {
	*valueStore

	value []int
}

func (vs *valueStoreInt64sSingle) Store(data string) error {
	result := regexpInt64sSingle.FindAllString(data, -1)

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
	mapValueStoreCreator[vtInt64sSingle] = func(vt ValueType) ValueStore {
		return &valueStoreInt64sSingle{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
