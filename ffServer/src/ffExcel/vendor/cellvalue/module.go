package cellvalue

import (
	"fmt"
)

// ValueType 定义标准的值类型定义
type ValueType interface {
	// IsArray 是否允许多个列组合成一个数组, 合并顺序为从左到右, 如果列的值为数组, 则进行数组合并
	IsArray() bool

	// Ignore 是否忽略本列配置
	IsIgnore() bool

	// IsGrammar 是不是grammar配置
	IsGrammar() bool

	// IsString 是不是字符串配置列
	IsString() bool

	// IsMap 是不是字典类型
	IsMap() bool

	// MapKeyGoType map 的 key 类型
	MapKeyGoType() string

	// MapValueGoType map 的 value 类型
	MapValueGoType() string

	// MapKeyProtoType map 的 key 类型
	MapKeyProtoType() string

	// MapValueProtoType map 的 value 类型
	MapValueProtoType() string

	// IsEnum 是不是枚举
	IsEnum() bool

	// HasEnum 是否有枚举
	HasEnum() bool

	// GoType 返回Go使用的类型的字符串描述
	GoType() string

	// ProtoType 返回该字段在Proto定义里的类型
	ProtoType() string

	toString() string
	valueType() *valueType
}

// NewValueType 根据值类型描述返回ValueType
func NewValueType(v string) (ValueType, error) {
	if vt := newValueType(v); vt != nil {
		return vt, nil
	}

	return nil, fmt.Errorf("invalid value type[%v]", v)
}

// ValueStore 定义标准的值类型存储
type ValueStore interface {
	// Store 将字符串形式的值存储起来
	Store(data string, vt ValueType) error

	// GoType 返回Go使用的类型的字符串描述
	GoType() string

	// Value 返回实际值
	Value() interface{}

	// ValueToml 返回导出toml时的字符串
	ValueToml() string

	// ValueTomlMapKeys 返回Map导出toml时的key字符串
	ValueTomlMapKeys() string

	// ValueTomlMapValues 返回Map导出toml时的value字符串
	ValueTomlMapValues() string

	String() string
}

// NewValueStore 根据ValueType和字符串形式的值, 返回存储了值的ValueType
func NewValueStore(vt ValueType) (ValueStore, error) {
	r := newValueStore(vt)
	if r == nil {
		return nil, fmt.Errorf("ValueType[%v] not exist creator", vt.valueType())
	}

	return r, nil
}

// InitEnum 外界设置enum的toml文件
func InitEnum(tomlFile string) {
	initEnum(tomlFile)
}
