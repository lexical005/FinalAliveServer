package uuid

import (
	"fmt"
)

const (
	// InvalidUUID 无效UUID
	InvalidUUID UUID = 0

	// uuidTimestampOffset uuid内的时间戳部分, 是相对北京时间2017年1月1日0时0分0秒对应的时间戳的偏移, 单位秒
	uuidTimestampOffset = 1483200000
)

//---------------------------------------------------------------------------------
// generate uuid, format:
//
// 0.............0	0.............0	0......0
// 32bit timestamp	12bit requester	20bit sn
// snowflake生成的ID整体上按照时间自增排序, 并且整个分布式系统内不会产生ID碰撞（由请求者字段来保证）
const (
	timestampBitCount = 32 // 42bit, 时间戳, 相对北京时间2017年1月1日0时0分0秒对应的时间戳的偏移, 可表达136年
	requesterBitCount = 12 // 12bit, 生成器全球唯一标识, 最多同时支持4096个UUID生成器
	snBitCount        = 20 // 20bit, 序号, 一秒内最多生成1048576个UUID编号(100万QPS)

	timestampBitMask = 0xffffffff // 32bit, 时间戳, 相对北京时间2017年1月1日0时0分0秒对应的时间戳的偏移, 可表达136年
	requesterBitMask = 0xFFF      // 12bit, 生成器全球唯一标识, 最多同时支持4096个UUID生成器
	snBitMask        = 0xfffff    // 20bit, 序号, 一秒内最多生成1048576个UUID编号(100万QPS)

	timestampBitOffset = requesterBitCount + snBitCount
	requesterBitOffset = snBitCount
	snBitOffset        = 0
)

// UUID 全球唯一标识，64位无符号整数
type UUID uint64

// Requester 请求者
func (u UUID) Requester() uint16 {
	return uint16((uint64(u) >> requesterBitOffset) & requesterBitMask)
}

// Timestamp 时间戳（毫秒）
func (u UUID) Timestamp() uint64 {
	return (uint64(u) >> timestampBitOffset) & timestampBitMask
}

// Value uint64值
func (u UUID) Value() uint64 {
	return uint64(u)
}

// String
func (u UUID) String() string {
	// 生成 UUID
	// 0.............0	0.............0	0......0
	// 32bit timestamp	12bit requester	20bit sn
	v := u.Value()
	timestampOffset := (v >> timestampBitOffset) & timestampBitMask
	requester := (v >> requesterBitOffset) & requesterBitMask
	sn := (v >> snBitOffset) & snBitMask
	return fmt.Sprintf("%d-%d-%d", timestampOffset, requester, sn)
}

// Generator uuid生成器
type Generator interface {
	// Gen 生成一个uuid
	Gen() UUID

	String() string
}

// NewUUID 根据value返回一个UUID
func NewUUID(value uint64) UUID {
	return UUID(value)
}
