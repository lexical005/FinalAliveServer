package cellvalue

import (
	"regexp"
	"strings"
)

// 本列配置的值都将被忽略
var vtEmpty = &valueType{
	ignore: true,

	excel:     "",
	goType:    "",
	protoType: "",
}

// 本列配置的值为32位有符号整型, 名称字段在头部只能出现一次
var vtInt32 = &valueType{
	isNumber: true,

	excel:     "int32",
	goType:    "int32",
	protoType: "int32",
}

// 本列配置的值为32位有符号整型数组
var vtInt32Array = &valueType{
	isNumber: true,

	isArray: true,

	excel:     "[]int32",
	goType:    "[]int32",
	protoType: "repeated int32",
}

// 本列配置的值为64位有符号整型, 名称字段在头部只能出现一次
var vtInt64 = &valueType{
	isNumber: true,

	excel:     "int64",
	goType:    "int64",
	protoType: "int64",
}

// 本列配置的值为64位有符号整型数组
var vtInt64Array = &valueType{
	isNumber: true,
	isArray:  true,

	excel:     "[]int64",
	goType:    "[]int64",
	protoType: "repeated int64",
}

// 本列配置的值为字符串, 名称字段在头部只能出现一次
var vtString = &valueType{
	isStr: true,

	excel:     "string",
	goType:    "string",
	protoType: "string",
}

// 本列配置的值为字符串数组
var vtStringArray = &valueType{
	isStr:   true,
	isArray: true,

	excel:     "[]string",
	goType:    "[]string",
	protoType: "repeated string",
}

// 本列配置的值为自定义语法语句, 名称字段在头部只能出现一次
var vtGrammar = &valueType{
	grammar: true,

	excel:     "grammar",
	goType:    "ffGrammar.Grammar",
	protoType: "Grammar",
}

type valueType struct {
	ignore  bool
	grammar bool

	hasEnum bool

	isNumber bool

	isStr bool

	isEnum   bool
	enumType string

	isArray          bool
	arrayValueType   string
	arrayValueIsEnum bool

	isMap           bool
	mapKeyType      string
	mapKeyIsEnum    bool
	mapValueType    string
	mapValueIsEnum  bool
	fixedLineMapKey string

	goKeyType   string
	goValueType string

	protoKeyType   string
	protoValueType string

	excel     string
	goType    string
	protoType string
}

func (vt *valueType) GoType() string {
	return vt.goType
}
func (vt *valueType) ProtoType() string {
	return vt.protoType
}
func (vt *valueType) IsIgnore() bool {
	return vt.ignore
}
func (vt *valueType) IsArray() bool {
	return vt.isArray
}
func (vt *valueType) IsGrammar() bool {
	return vt.grammar
}
func (vt *valueType) IsString() bool {
	return vt.isStr
}
func (vt *valueType) IsMap() bool {
	return vt.isMap
}
func (vt *valueType) MapKeyGoType() string {
	return vt.goKeyType
}
func (vt *valueType) MapValueGoType() string {
	return vt.goValueType
}
func (vt *valueType) MapKeyProtoType() string {
	return vt.protoKeyType
}
func (vt *valueType) MapValueProtoType() string {
	return vt.protoValueType
}
func (vt *valueType) IsEnum() bool {
	return vt.isEnum
}
func (vt *valueType) HasEnum() bool {
	return vt.hasEnum
}
func (vt valueType) toString() string {
	return vt.excel
}

func (vt *valueType) valueType() *valueType {
	return vt
}

// 允许用户配置的值类型
var basicValueType = map[string]*valueType{
	vtEmpty.excel:       vtEmpty,
	vtInt32.excel:       vtInt32,
	vtInt32Array.excel:  vtInt32Array,
	vtInt64.excel:       vtInt64,
	vtInt64Array.excel:  vtInt64Array,
	vtString.excel:      vtString,
	vtStringArray.excel: vtStringArray,
	vtGrammar.excel:     vtGrammar,
}

// 基础类型是不是数值
var basicGoType = map[string]bool{
	"int32":  true,
	"int64":  true,
	"string": false,
}

var regexpMap = regexp.MustCompile(`map\[([\w\.]+)\]([\w]+)`)

func newValueType(v string) *valueType {
	// 基本类型
	if v, ok := basicValueType[v]; ok {
		return v
	}

	enumType := v
	if strings.HasPrefix(v, "[]") {
		enumType = v[2:]
		if _, ok := mapEnums[enumType]; !ok {
			return nil
		}

		return &valueType{
			hasEnum: true,

			isArray:          true,
			arrayValueType:   enumType,
			arrayValueIsEnum: true,

			excel:     v,
			goType:    "[]ffEnum." + enumType,
			protoType: "repeated int32",
		}

	} else if strings.HasPrefix(v, "map[") {
		i := strings.Index(v, "]")
		mapKeyType := v[len("map["):i]
		mapValueType := v[i+1:]

		tmpGoKey, tmpGoValue := "", ""
		tmpProtoKey, tmpProtoValue := mapKeyType, mapValueType

		fixedLineMapKey := ""

		temp := strings.Split(mapKeyType, ".")
		if len(temp) == 2 {
			mapKeyType, fixedLineMapKey = temp[0], temp[1]
		}

		enumKeys, mapKeyIsEnum := mapEnums[mapKeyType]
		if !mapKeyIsEnum {
			if _, ok := basicGoType[mapKeyType]; !ok {
				return nil
			}
			tmpGoKey = mapKeyType
		} else {
			tmpGoKey = "ffEnum." + mapKeyType
			tmpProtoKey = "int32"

			if len(fixedLineMapKey) > 0 {
				found := false
				for _, key := range enumKeys {
					if fixedLineMapKey == key {
						found = true
						break
					}
				}
				if !found {
					return nil
				}
			}
		}
		_, mapValueIsEnum := mapEnums[mapValueType]
		if !mapValueIsEnum {
			if _, ok := basicGoType[mapValueType]; !ok {
				return nil
			}
			tmpGoValue = mapValueType
		} else {
			tmpGoValue = "ffEnum." + mapValueType
			tmpProtoValue = "int32"
		}

		return &valueType{
			hasEnum: mapKeyIsEnum || mapValueIsEnum,

			isMap:           true,
			mapKeyType:      mapKeyType,
			mapKeyIsEnum:    mapKeyIsEnum,
			mapValueType:    mapValueType,
			mapValueIsEnum:  mapValueIsEnum,
			fixedLineMapKey: fixedLineMapKey,

			goKeyType:   "[]" + tmpGoKey,
			goValueType: "[]" + tmpGoValue,

			protoKeyType:   "repeated " + tmpProtoKey,
			protoValueType: "repeated " + tmpProtoValue,

			excel:     v,
			goType:    "map[" + tmpGoKey + "]" + tmpGoValue,
			protoType: "map<" + tmpProtoKey + "," + tmpProtoValue + ">",
		}
	} else {
		if _, ok := mapEnums[enumType]; !ok {
			return nil
		}

		return &valueType{
			hasEnum: true,

			isEnum:   true,
			enumType: enumType,

			excel:     v,
			goType:    "ffEnum." + enumType,
			protoType: "int32",
		}
	}
}
