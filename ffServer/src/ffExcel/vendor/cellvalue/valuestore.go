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

func (vs *valueStore) Type() string {
	return vs.vt.Type()
}
func (vs *valueStore) Value() interface{} {
	return vs.value
}
func (vs *valueStore) ValueToml() string {
	rv := reflect.ValueOf(vs.value)
	if rv.Kind() == reflect.Array || rv.Kind() == reflect.Slice {
		result := "["
		for i := 0; i < rv.Len(); i++ {
			if i > 0 {
				result += ", "
			}
			v := rv.Index(i)
			k := v.Kind()
			if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 {
				result += fmt.Sprintf("%v", v.Int())
			} else if k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 {
				return fmt.Sprintf("%v", v.Int())
			} else if v.Kind() == reflect.String {
				s := v.String()
				s = strings.Replace(s, "\"", "\\\"", -1)
				result += fmt.Sprintf("\"%v\"", s)
			}
		}
		return result + "]"
	}

	k := rv.Kind()
	if k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64 {
		return fmt.Sprintf("%v", rv.Int())
	} else if k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64 {
		return fmt.Sprintf("%v", rv.Uint())
	} else if rv.Kind() == reflect.String {
		s := rv.String()
		s = strings.Replace(s, "\"", "\\\"", -1)
		return fmt.Sprintf("\"%v\"", s)
	}

	panic(fmt.Sprintf("ValueToml failed: ValueType[%v] value[%v:%v]", vs.vt, rv.Kind().String(), vs.value))
}

func (vs *valueStore) String() string {
	return fmt.Sprintf("[%v:%v]", vs.Type(), vs.value)
}

var mapValueStoreCreator = make(map[valueType]func(vt ValueType) ValueStore)
