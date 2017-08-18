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

func (vs *valueStoreMap) Store(data string) error {
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

func checkMapKeyValueType(desc string) (bool, string, string) {
	result := regexpMap.FindAllString(desc, -1)
	if len(result) == 2 {
		return true, result[0], result[1]
	}
	return false, "", ""
}

func init() {
	mapValueStoreCreator[vtMap] = func(vt ValueType) ValueStore {
		_, keyType, valueType := checkMapKeyValueType(vt.toString())

		// 键是不是枚举
		keyTypeEnum := true
		if keyType == vtInt32.toString() || keyType == vtInt64.toString() || keyType == vtString.toString() {
			keyTypeEnum = false
		}

		// 如果键是枚举, 则判定键值, 是配置在键里的还是在值里的
		keyEnumFixed := keyTypeEnum && strings.Contains(keyType, ".")

		return &valueStoreMap{
			keyType:   keyType,
			valueType: valueType,

			keyTypeEnum:  keyTypeEnum,
			keyEnumFixed: keyEnumFixed,

			valueStore: &valueStore{
				vt: vt,
			},
		}
	}
}
