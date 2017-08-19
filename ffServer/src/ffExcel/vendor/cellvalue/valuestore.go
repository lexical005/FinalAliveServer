package cellvalue

import (
	"fmt"
	"reflect"
	"strings"
)

type valueStore struct {
	vt ValueType

	value interface{}
}

func (vs *valueStore) GoType() string {
	return vs.vt.GoType()
}
func (vs *valueStore) Value() interface{} {
	return vs.value
}
func (vs *valueStore) ValueToml() string {
	t := vs.vt.valueType()
	rv := reflect.ValueOf(vs.value)

	// 基本类型
	if _, ok := basicValueType[t.excel]; ok {
		// 数组
		if t.IsArray() {
			result := "["
			for i := 0; i < rv.Len(); i++ {
				if i > 0 {
					result += ", "
				}
				v := rv.Index(i)
				if t.isNumber {
					result += fmt.Sprintf("%v", v.Int())
				} else {
					s := v.String()
					s = strings.Replace(s, "\"", "\\\"", -1)
					result += fmt.Sprintf("\"%v\"", s)
				}
			}
			return result + "]"
		}

		// 数值, 字符串, Grammar
		if t.isNumber {
			return fmt.Sprintf("%v", rv.Int())
		} else if rv.Kind() == reflect.String {
			s := rv.String()
			s = strings.Replace(s, "\"", "\\\"", -1)
			return fmt.Sprintf("\"%v\"", s)
		}
	}

	panic(fmt.Sprintf("ValueToml failed: ValueType[%v] value[%v:%v]", vs.vt, rv.Kind().String(), vs.value))
}

func (vs *valueStore) String() string {
	return fmt.Sprintf("[%v:%v]", vs.GoType(), vs.value)
}

var basicValueStoreCreator = make(map[*valueType]func(vt ValueType) ValueStore)

// NewValueStore 根据ValueType和字符串形式的值, 返回存储了值的ValueType
func newValueStore(vt ValueType) ValueStore {
	t := vt.valueType()
	creator, ok := basicValueStoreCreator[t]
	if ok {
		return creator(vt)
	}

	if t.IsMap() {
		return newValueStoreMap(t)
	} else if t.IsArray() {
		return newValueStoreEnumArray(t)
	} else if t.IsEnum() {
		return newValueStoreEnum(t)
	}

	return nil
}
