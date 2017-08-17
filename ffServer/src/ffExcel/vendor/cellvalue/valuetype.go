package cellvalue

const (
	// 本列配置的值都将被忽略
	vtEmpty = valueType("")

	// 本列配置的值为32位有符号整型, 名称字段在头部只能出现一次
	vtInt32 = valueType("int32")

	// 本列配置的值为32位有符号整型数组
	vtInt32Array = valueType("[]int32")

	// 本列配置的值为64位有符号整型, 名称字段在头部只能出现一次
	vtInt64 = valueType("int64")

	// 本列配置的值为64位有符号整型数组
	vtInt64Array = valueType("[]int64")

	// 本列配置的值为字符串, 名称字段在头部只能出现一次
	vtString = valueType("string")

	// 本列配置的值为字符串数组
	vtStringArray = valueType("[]string")

	// 本列配置的值为自定义语法语句, 名称字段在头部只能出现一次
	vtGrammar = valueType("grammar")
)

// 允许用户配置的值类型
var vtExist = map[string]valueType{
	string(vtEmpty):       vtEmpty,
	string(vtInt32):       vtInt32,
	string(vtInt32Array):  vtInt32Array,
	string(vtInt64):       vtInt64,
	string(vtInt64Array):  vtInt64Array,
	string(vtString):      vtString,
	string(vtStringArray): vtStringArray,
	string(vtGrammar):     vtGrammar,
}

// 用户配置的值类型中, 有效的值类型到程序内部类型的匹配
var mapValueTypeToRealType = map[valueType]string{
	vtInt32:       "int32",
	vtInt32Array:  "[]int32",
	vtInt64:       "int64",
	vtInt64Array:  "[]int64",
	vtString:      "string",
	vtStringArray: "[]string",
	vtGrammar:     "ffGrammar.Grammar",
}

// 用户配置的值类型中, 有效的值类型到程序内部类型的匹配
var mapValueTypeToProtoType = map[valueType]string{
	vtInt32:       "int32",
	vtInt32Array:  "repeated int32",
	vtInt64:       "int64",
	vtInt64Array:  "repeated int64",
	vtString:      "string",
	vtStringArray: "repeated string",
	vtGrammar:     "Grammar",
}

type valueType string

func (vt *valueType) Type() string {
	return mapValueTypeToRealType[*vt]
}
func (vt *valueType) ProtoType() string {
	return mapValueTypeToProtoType[*vt]
}
func (vt *valueType) IsIgnore() bool {
	return *vt == vtEmpty
}
func (vt *valueType) IsMulti() bool {
	return *vt == vtInt32Array || *vt == vtInt64Array || *vt == vtStringArray
}
func (vt *valueType) IsGrammar() bool {
	return *vt == vtGrammar
}
func (vt *valueType) IsString() bool {
	return *vt == vtString || *vt == vtStringArray
}

func (vt *valueType) valueType() valueType {
	return *vt
}
