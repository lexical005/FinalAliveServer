package cellvalue

type valueStoreEnumArray struct {
	*valueStore

	value string
}

func (vs *valueStoreEnumArray) Store(data string, vt ValueType) error {
	vs.value = data
	vs.valueStore.value = vs.value
	return nil
}

func newValueStoreEnumArray(vt *valueType) *valueStoreEnumArray {
	return &valueStoreEnumArray{
		valueStore: &valueStore{
			vt: vt,
		},
	}
}
