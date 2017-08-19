package cellvalue

type valueStoreEmpty struct {
	*valueStore

	value string
}

func (vs *valueStoreEmpty) Store(data string, vt ValueType) error {
	vs.value = data
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	basicValueStoreCreator[vtEmpty] = func(vt ValueType) ValueStore {
		return &valueStoreEmpty{
			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
