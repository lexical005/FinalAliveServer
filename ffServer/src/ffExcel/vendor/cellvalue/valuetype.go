package cellvalue

const (
	// 本列配置的值都将被忽略
	vtEmpty = valueType("")

	// 本列配置的值为整型, 名称字段在头部只能出现一次
	vtInt = valueType("int")

	// 本列配置的值为整型数组, 名称字段在头部只能出现一次
	vtIntsSingle = valueType("ints_single")

	// 由多列整型配置合并成整型数组, 名称字段在头部允许出现多次, 追加到数组的顺序为从左向右
	vtIntsMulti = valueType("ints_multi")

	// 本列配置的值为字符串, 名称字段在头部只能出现一次
	vtString = valueType("str")

	// 本列配置的值为字符串数组, 名称字段在头部只能出现一次
	vtStringsSingle = valueType("strs_single")

	// 由多列字符串配置合并成字符串数组, 名称字段在头部允许出现多次, 追加到数组的顺序为从左向右
	vtStringsMulti = valueType("strs_multi")

	// 本列配置的值为自定义语法语句, 名称字段在头部只能出现一次
	vtGrammar = valueType("grammar")

	// 由多列自定义语法语句配置合并成自定义语法语句数组, 名称字段在头部允许出现多次, 追加到数组的顺序为从左向右
	vtGrammarsMulti = valueType("grammars_multi")
)

// 允许用户配置的值类型
var vtExist = map[string]valueType{
	string(vtEmpty):         vtEmpty,
	string(vtInt):           vtInt,
	string(vtIntsSingle):    vtIntsSingle,
	string(vtIntsMulti):     vtIntsMulti,
	string(vtString):        vtString,
	string(vtStringsSingle): vtStringsSingle,
	string(vtStringsMulti):  vtStringsMulti,
	string(vtGrammar):       vtGrammar,
	string(vtGrammarsMulti): vtGrammarsMulti,
}

// 用户配置的值类型中, 有效的值类型到程序内部类型的匹配
var mapValueTypeToRealType = map[valueType]string{
	vtInt:           "int",
	vtIntsSingle:    "[]int",
	vtIntsMulti:     "[]int",
	vtString:        "string",
	vtStringsSingle: "[]string",
	vtStringsMulti:  "[]string",
	vtGrammar:       "ffGrammar.Grammar",
	vtGrammarsMulti: "[]ffGrammar.Grammar",
}

type valueType string

func (vt *valueType) Type() string {
	return mapValueTypeToRealType[*vt]
}
func (vt *valueType) IsIgnore() bool {
	return *vt == vtEmpty
}
func (vt *valueType) IsMulti() bool {
	return *vt == vtIntsMulti || *vt == vtStringsMulti || *vt == vtGrammarsMulti
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
