package cellvalue

type valueStoreString struct {
	*valueStore

	value string
}

func (vs *valueStoreString) Store(data string) error {
	vs.value = data
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtString] = func(vt ValueType) ValueStore {
		return &valueStoreString{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
