package cellvalue

type valueStoreStringsMulti struct {
	*valueStore

	value []string
}

func (vs *valueStoreStringsMulti) Store(data string) error {
	vs.value = append(vs.value, data)
	vs.valueStore.value = vs.value
	return nil
}

func init() {
	mapValueStoreCreator[vtStringsMulti] = func(vt ValueType) ValueStore {
		return &valueStoreStringsMulti{
			valueStore: &valueStore{
				vt: vt,
			},
			value: make([]string, 0, 1),
		}
	}
}
