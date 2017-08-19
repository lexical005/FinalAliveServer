package cellvalue

import (
	"fmt"
)

type valueStoreEnum struct {
	*valueStore

	value string
}

func (vs *valueStoreEnum) Store(data string, vt ValueType) error {
	t := vt.valueType()

	allEnumKeys := mapEnums[t.enumType]
	for _, key := range allEnumKeys {
		if key == data {
			vs.value = t.enumType + "." + data
			vs.valueStore.value = vs.value
			return nil
		}
	}

	return fmt.Errorf("valueStoreEnum.Store invalid data[%v] valueType[%v]", data, vt.toString())
}

func newValueStoreEnum(vt *valueType) *valueStoreEnum {
	return &valueStoreEnum{
		valueStore: &valueStore{
			vt: vt,
		},
	}
}
