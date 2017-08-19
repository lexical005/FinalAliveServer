package cellvalue

import (
	"fmt"
	"regexp"
	"strings"
)

var regexpStringArray = regexp.MustCompile(`([\w]+)`)

type valueStoreEnumArray struct {
	*valueStore

	value []string
}

func (vs *valueStoreEnumArray) Store(data string, vt ValueType) error {
	if vs.value == nil {
		vs.value = make([]string, 0, 1)
	}

	t := vt.valueType()
	allEnumKeys := mapEnums[t.arrayValueType]

	if strings.HasPrefix(data, "[") && strings.HasSuffix(data, "]") {
		result := regexpInt32Array.FindAllString(data, -1)
		for _, s := range result {
			found := false
			for _, key := range allEnumKeys {
				if key == s {
					vs.value = append(vs.value, t.arrayValueType+"."+s)
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("valueStoreEnum.Store invalid data[%v] valueType[%v]", data, vt.toString())
			}
		}
		vs.valueStore.value = vs.value
		return nil
	}

	for _, key := range allEnumKeys {
		if key == data {
			vs.value = append(vs.value, t.arrayValueType+"."+data)
			vs.valueStore.value = vs.value
			return nil
		}
	}

	return fmt.Errorf("valueStoreEnum.Store invalid data[%v] valueType[%v]", data, vt.toString())
}

func newValueStoreEnumArray(vt *valueType) *valueStoreEnumArray {
	return &valueStoreEnumArray{
		valueStore: &valueStore{
			vt: vt,
		},
	}
}
