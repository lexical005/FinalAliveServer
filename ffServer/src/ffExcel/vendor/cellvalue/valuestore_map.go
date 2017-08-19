package cellvalue

import (
	"strings"
)

type mapValueData struct {
	key   string
	value string
}

type valueStoreMap struct {
	*valueStore

	keyType   string
	valueType string

	keyTypeEnum  bool
	keyEnumFixed bool

	data *mapValueData
}

func (vs *valueStoreMap) Store(data string, vt ValueType) error {
	if vs.data == nil {
		vs.data = &mapValueData{}
	}

	if vs.keyEnumFixed {
		vs.data.key = vs.keyType
		vs.data.value = data
	} else {
		tmp := strings.Split(data, ",")
		vs.data.key = tmp[0]
		vs.data.value = tmp[1]
	}

	vs.valueStore.value = vs.data
	return nil
}

func newValueStoreMap(vt *valueType) *valueStoreMap {
	return &valueStoreMap{
		valueStore: &valueStore{
			vt: vt,
		},
	}
}
