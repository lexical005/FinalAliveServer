package cellvalue

import (
	"fmt"
	"strings"
)

type mapInfo struct {
	mapKeyIsNumber bool
	mapKeys        []string

	mapValueIsNumber bool
	mapValues        []string
}

type valueStoreMap struct {
	*valueStore

	mapInfo *mapInfo
}

func (vs *valueStoreMap) Store(data string, vt ValueType) error {
	t := vt.valueType()

	if vs.mapInfo == nil {

		mapKeyIsNumber := false
		if isNumber, ok := basicGoType[t.mapKeyType]; ok && isNumber {
			mapKeyIsNumber = true
		}

		mapValueIsNumber := false
		if isNumber, ok := basicGoType[t.mapValueType]; ok && isNumber {
			mapValueIsNumber = true
		}

		vs.mapInfo = &mapInfo{
			mapKeyIsNumber: mapKeyIsNumber,
			mapKeys:        make([]string, 0, 1),

			mapValueIsNumber: mapValueIsNumber,
			mapValues:        make([]string, 0, 1),
		}
	}

	datas := strings.Split(data, ",")
	for i, s := range datas {
		datas[i] = strings.TrimSpace(s)
	}

	if !t.hasEnum {
		if len(t.fixedLineMapKey) > 0 {
			vs.mapInfo.mapKeys = append(vs.mapInfo.mapKeys, t.fixedLineMapKey)
			vs.mapInfo.mapValues = append(vs.mapInfo.mapValues, datas[0])
		} else {
			vs.mapInfo.mapKeys = append(vs.mapInfo.mapKeys, datas[0])
			vs.mapInfo.mapValues = append(vs.mapInfo.mapValues, datas[1])
		}
	} else {
		if t.mapKeyIsEnum {
			if len(t.fixedLineMapKey) > 0 {
				vs.mapInfo.mapKeys = append(vs.mapInfo.mapKeys, t.mapKeyType+"."+t.fixedLineMapKey)
			} else {
				found := false
				enumKeys, _ := mapEnums[t.mapKeyType]
				for _, key := range enumKeys {
					if key == datas[0] {
						found = true
						vs.mapInfo.mapKeys = append(vs.mapInfo.mapKeys, t.mapKeyType+"."+datas[0])
						break
					}
				}

				if !found {
					return fmt.Errorf("valueStoreMap.Store invalid data[%v] for valueType[%v]", data, vs.GoType())
				}
			}
		} else {
			if len(t.fixedLineMapKey) > 0 {
				vs.mapInfo.mapKeys = append(vs.mapInfo.mapKeys, t.fixedLineMapKey)
			} else {
				vs.mapInfo.mapKeys = append(vs.mapInfo.mapKeys, datas[0])
			}
		}

		if t.mapValueIsEnum {
			if len(t.fixedLineMapKey) > 0 {
				vs.mapInfo.mapValues = append(vs.mapInfo.mapValues, t.mapValueType+"."+datas[0])
			} else {
				vs.mapInfo.mapValues = append(vs.mapInfo.mapValues, t.mapValueType+"."+datas[1])
			}
		} else {
			if len(t.fixedLineMapKey) > 0 {
				vs.mapInfo.mapValues = append(vs.mapInfo.mapValues, datas[0])
			} else {
				vs.mapInfo.mapValues = append(vs.mapInfo.mapValues, datas[1])
			}
		}
	}

	vs.valueStore.value = vs.mapInfo
	return nil
}

func newValueStoreMap(vt *valueType) *valueStoreMap {
	return &valueStoreMap{
		valueStore: &valueStore{
			vt: vt,
		},
	}
}
