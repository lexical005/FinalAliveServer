package uuid

import (
	"ffCommon/log/log"

	"fmt"
	"sync"
	"time"
)

// UUID 全球唯一标识，64位无符号整数
type UUID uint64

const (
	// InvalidUUID 无效UUID
	InvalidUUID = 0

	// uuid内的时间戳部分, 是相对北京时间2017年1月1日0时0分0秒对应的时间戳的偏移, 单位秒
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

// Requester 请求者
func (u UUID) Requester() uint16 {
	return uint16((uint64(u) >> requesterBitOffset) & requesterBitMask)
}

// Timestamp 时间戳（毫秒）
func (u UUID) Timestamp() uint64 {
	return (uint64(u) >> timestampBitOffset) & timestampBitMask
}

// Generator UUID生成器，不支持 multi goroutine call
type Generator struct {
	requester uint64

	sn        uint64
	timestamp int64
}

// Gen 生成一个UUID
// 如果当前毫秒内的所有可用的 UUID 都已经生成完毕(16个), 将阻塞, 直到下一毫秒到来再返回
func (gen *Generator) Gen() UUID {
	// 以秒为单位的时间戳
	now := time.Now().Unix()
	if now < gen.timestamp {
		log.FatalLogger.Println("UUID.Gen clock shift happened, waiting until the clock moving to the next millisecond.")
		now = gen.timestamp
	}

	// 时间戳处理以及序号生成
	if gen.timestamp < now {
		gen.sn = 0
		gen.timestamp = now
	} else {
		gen.sn++
		if gen.sn > snBitMask {
			// 旋转锁等待下一毫秒的到来
			for now <= gen.timestamp {
				now = time.Now().Unix()
			}

			gen.sn = 0
			gen.timestamp = now
		}
	}

	// 生成 UUID
	// 0.............0	0.............0	0......0
	// 32bit timestamp	12bit requester	20bit sn
	return UUID((uint64(gen.timestamp-uuidTimestampOffset) << timestampBitOffset) + gen.requester + gen.sn)
}
func (gen *Generator) String() string {
	return fmt.Sprintf("requester[%v] sn[%v] timestamp[%v]", gen.requester, gen.sn, gen.timestamp)
}

// GeneratorSafe support multi goroutine call
type GeneratorSafe struct {
	*Generator

	muLock sync.Mutex
}

// Gen 生成 UUID
func (g *GeneratorSafe) Gen() UUID {
	g.muLock.Lock()
	defer g.muLock.Unlock()

	return g.Generator.Gen()
}

// NewGenerator 返回一个Generator，不支持 multi goroutine call
// requester: 使用者自定义的12位无符号数, 可根据项目需要, 全局规划这12位
func NewGenerator(requester uint64) (*Generator, error) {
	if requester > requesterBitMask {
		return nil, fmt.Errorf("uuid.NewGenerator: invalid requester, must between[0, %v]", requesterBitMask)
	}

	return &Generator{
		requester: requester << requesterBitOffset,
	}, nil
}

// NewGeneratorSafe 返回一个GeneratorSafe，支持 multi goroutine call
// requester: 使用者自定义的12位无符号数, 可根据项目需要, 全局规划这12位
func NewGeneratorSafe(requester uint64) (*GeneratorSafe, error) {
	if requester > requesterBitMask {
		return nil, fmt.Errorf("uuid.NewGeneratorSafe: invalid requester, must between[0, %v]", requesterBitMask)
	}

	return &GeneratorSafe{
		Generator: &Generator{
			requester: requester << requesterBitOffset,
		},
	}, nil
}
