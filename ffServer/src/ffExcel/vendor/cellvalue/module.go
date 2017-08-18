package cellvalue

import (
	"fmt"
)

// ValueType 定义标准的值类型定义
type ValueType interface {
	// IsMulti 是否允许多个列组合成一个数组, 合并顺序为从左到右, 如果列的值为数组, 则进行数组合并
	IsMulti() bool

	// Ignore 是否忽略本列配置
	IsIgnore() bool

	// IsGrammar 是不是grammar配置
	IsGrammar() bool

	// IsString 是不是字符串配置列
	IsString() bool

	// IsMap 是不是字典类型
	IsMap() bool

	// Type 返回程序内部使用的类型的字符串描述
	Type() string

	// ProtoType 返回该字段在Proto定义里的类型
	ProtoType() string

	toString() string
	valueType() valueType
}

// NewValueType 根据值类型描述返回ValueType
func NewValueType(v string) (ValueType, error) {
	if _, ok := vtExist[v]; ok {
		vt := valueType(v)
		return &vt, nil
	}

	success, _, _ := checkMapKeyValueType(v)
	if success {
		vt := valueType(v)
		return &vt, nil
	}

	return nil, fmt.Errorf("invalid value type[%v]", v)
}

// ValueStore 定义标准的值类型存储
type ValueStore interface {
	// Store 将字符串形式的值存储起来
	Store(data string) error

	// Type 返回程序内部使用的类型的字符串描述
	Type() string

	// Value 返回实际值
	Value() interface{}

	// ValueToml 返回导出toml时的字符串
	ValueToml() string

	String() string
}

// NewValueStore 根据ValueType和字符串形式的值, 返回存储了值的ValueType
func NewValueStore(vt ValueType) (ValueStore, error) {
	creator, ok := mapValueStoreCreator[vt.valueType()]
	if !ok {
		return nil, fmt.Errorf("ValueType[%v] not exist creator", vt.valueType())
	}

	r := creator(vt)
	return r, nil
}
