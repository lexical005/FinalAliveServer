package cellvalue

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var regexpInt64Array = regexp.MustCompile(`[\d]+`)

type valueStoreInt64Array struct {
	*valueStore

	value []int64
}

func (vs *valueStoreInt64Array) Store(data string) error {
	if vs.value == nil {
		vs.value = make([]int64, 0, 1)
	}

	if strings.HasPrefix(data, "[") && strings.HasSuffix(data, "]") {
		result := regexpInt64Array.FindAllString(data, -1)
		for _, s := range result {
			i64, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				return fmt.Errorf("valueStoreInt64Array.ValueStore[%v] Invalid int array data[%v]", vs.GoType(), data)
			}
			vs.value = append(vs.value, int64(i64))
		}
	} else {
		i64, err := strconv.ParseInt(data, 10, 0)
		if err != nil {
			return fmt.Errorf("valueStoreInt64Array.ValueStore[%v] Invalid int array data[%v]", vs.GoType(), data)
		}
		vs.value = append(vs.value, int64(i64))
	}

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
