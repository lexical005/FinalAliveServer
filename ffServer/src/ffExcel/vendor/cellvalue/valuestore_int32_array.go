package cellvalue

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var regexpInt32Array = regexp.MustCompile(`([-\d]+)`)

type valueStoreInt32Array struct {
	*valueStore

	value []int32
}

func (vs *valueStoreInt32Array) Store(data string, vt ValueType) error {
	if vs.value == nil {
		vs.value = make([]int32, 0, 1)
	}

	if strings.HasPrefix(data, "[") && strings.HasSuffix(data, "]") {
		result := regexpInt32Array.FindAllString(data, -1)
		for _, s := range result {
			i64, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				return fmt.Errorf("valueStoreInt32Array.ValueStore[%v] Invalid int array data[%v]", vs.GoType(), data)
			}
			vs.value = append(vs.value, int32(i64))
		}
	} else {
		i64, err := strconv.ParseInt(data, 10, 0)
		if err != nil {
			return fmt.Errorf("valueStoreInt32Array.ValueStore[%v] Invalid int array data[%v]", vs.GoType(), data)
		}
		vs.value = append(vs.value, int32(i64))
	}

	vs.valueStore.value = vs.value
	return nil
}

func init() {
	basicValueStoreCreator[vtInt32Array] = func(vt ValueType) ValueStore {
		return &valueStoreInt32Array{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
