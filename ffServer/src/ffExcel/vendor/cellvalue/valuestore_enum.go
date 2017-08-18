package cellvalue

type valueStoreEnum struct {
	*valueStore

	value string
}

func (vs *valueStoreEnum) Store(data string) error {
	vs.value = data
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtEnum] = func(vt ValueType) ValueStore {
		return &valueStoreEnum{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
