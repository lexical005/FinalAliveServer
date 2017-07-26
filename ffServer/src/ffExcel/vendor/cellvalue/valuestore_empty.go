package cellvalue

type valueStoreEmpty struct {
	*valueStore

	value string
}

func (vs *valueStoreEmpty) Store(data string) error {
	vs.value = data
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtEmpty] = func(vt ValueType) ValueStore {
		return &valueStoreEmpty{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
