package cellvalue

type valueStoreEnum struct {
	*valueStore

	value string
}

func (vs *valueStoreEnum) Store(data string, vt ValueType) error {
	vs.value = data
	vs.valueStore.value = vs.value
	return nil
}

func newValueStoreEnum(vt *valueType) *valueStoreEnum {
	return &valueStoreEnum{
		valueStore: &valueStore{
			vt: vt,
		},
	}
}
