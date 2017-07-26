package cellvalue

import (
	"fmt"
	"regexp"
	"strconv"
)

var regexpIntsSingle = regexp.MustCompile(`[\d]+`)

type valueStoreIntsSingle struct {
	*valueStore

	value []int
}

func (vs *valueStoreIntsSingle) Store(data string) error {
	result := regexpIntsSingle.FindAllString(data, -1)

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
	mapValueStoreCreator[vtIntsSingle] = func(vt ValueType) ValueStore {
		return &valueStoreIntsSingle{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
