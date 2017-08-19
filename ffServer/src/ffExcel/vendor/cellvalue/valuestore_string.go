package cellvalue

type valueStoreString struct {
	*valueStore

	value string
}

func (vs *valueStoreString) Store(data string, vt ValueType) error {
	vs.value = data
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	basicValueStoreCreator[vtString] = func(vt ValueType) ValueStore {
		return &valueStoreString{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
