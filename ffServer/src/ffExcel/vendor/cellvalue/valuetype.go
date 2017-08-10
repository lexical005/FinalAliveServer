package cellvalue

const (
	// 本列配置的值都将被忽略
	vtEmpty = valueType("")

	// 本列配置的值为32位有符号整型, 名称字段在头部只能出现一次
	vtInt32 = valueType("int32")

	// 本列配置的值为32位有符号整型数组, 名称字段在头部只能出现一次
	vtInt32sSingle = valueType("int32s_single")

	// 由多列32位有符号整型配置合并成32位有符号整型数组, 名称字段在头部允许出现多次, 追加到数组的顺序为从左向右
	vtInt32sMulti = valueType("int32s_multi")

	// 本列配置的值为64位有符号整型, 名称字段在头部只能出现一次
	vtInt64 = valueType("int64")

	// 本列配置的值为64位有符号整型数组, 名称字段在头部只能出现一次
	vtInt64sSingle = valueType("int64s_single")

	// 由多列64位有符号整型配置合并成64位有符号整型数组, 名称字段在头部允许出现多次, 追加到数组的顺序为从左向右
	vtInt64sMulti = valueType("int64s_multi")

	// 本列配置的值为字符串, 名称字段在头部只能出现一次
	vtString = valueType("str")

	// 本列配置的值为字符串数组, 名称字段在头部只能出现一次
	vtStringsSingle = valueType("strs_single")

	// 由多列字符串配置合并成字符串数组, 名称字段在头部允许出现多次, 追加到数组的顺序为从左向右
	vtStringsMulti = valueType("strs_multi")

	// 本列配置的值为自定义语法语句, 名称字段在头部只能出现一次
	vtGrammar = valueType("grammar")

	// 由多列自定义语法语句配置合并成自定义语法语句数组, 名称字段在头部允许出现多次, 追加到数组的顺序为从左向右
	// todo: 实现
	vtGrammarsMulti = valueType("grammars_multi")
)

// 允许用户配置的值类型
var vtExist = map[string]valueType{
	string(vtEmpty):         vtEmpty,
	string(vtInt32):         vtInt32,
	string(vtInt32sSingle):  vtInt32sSingle,
	string(vtInt32sMulti):   vtInt32sMulti,
	string(vtInt64):         vtInt64,
	string(vtInt64sSingle):  vtInt64sSingle,
	string(vtInt64sMulti):   vtInt64sMulti,
	string(vtString):        vtString,
	string(vtStringsSingle): vtStringsSingle,
	string(vtStringsMulti):  vtStringsMulti,
	string(vtGrammar):       vtGrammar,
	string(vtGrammarsMulti): vtGrammarsMulti,
}

// 用户配置的值类型中, 有效的值类型到程序内部类型的匹配
var mapValueTypeToRealType = map[valueType]string{
	vtInt32:         "int32",
	vtInt32sSingle:  "[]int32",
	vtInt32sMulti:   "[]int32",
	vtInt64:         "int64",
	vtInt64sSingle:  "[]int64",
	vtInt64sMulti:   "[]int64",
	vtString:        "string",
	vtStringsSingle: "[]string",
	vtStringsMulti:  "[]string",
	vtGrammar:       "ffGrammar.Grammar",
	vtGrammarsMulti: "[]ffGrammar.Grammar",
}

// 用户配置的值类型中, 有效的值类型到程序内部类型的匹配
var mapValueTypeToProtoType = map[valueType]string{
	vtInt32:         "int32",
	vtInt32sSingle:  "repeated int32",
	vtInt32sMulti:   "repeated int32",
	vtInt64:         "int64",
	vtInt64sSingle:  "repeated int64",
	vtInt64sMulti:   "repeated int64",
	vtString:        "string",
	vtStringsSingle: "repeated string",
	vtStringsMulti:  "repeated string",
	vtGrammar:       "Grammar",
	vtGrammarsMulti: "Grammar",
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
	return *vt == vtInt32sMulti || *vt == vtInt64sMulti || *vt == vtStringsMulti || *vt == vtGrammarsMulti
}
func (vt *valueType) IsGrammar() bool {
	return *vt == vtGrammar || *vt == vtGrammarsMulti
}
func (vt *valueType) IsString() bool {
	return *vt == vtString || *vt == vtStringsMulti || *vt == vtStringsSingle
}

func (vt *valueType) valueType() valueType {
	return *vt
}
