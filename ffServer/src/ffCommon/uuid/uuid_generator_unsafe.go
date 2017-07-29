package uuid

import (
	"ffCommon/log/log"
	"fmt"
	"time"
)

// uuidGeneratorUnsafe UUID生成器，不支持 multi goroutine call
type uuidGeneratorUnsafe struct {
	requester uint64

	sn        uint64
	timestamp int64
}

// Gen 生成一个UUID
// 如果当前毫秒内的所有可用的 UUID 都已经生成完毕(16个), 将阻塞, 直到下一毫秒到来再返回
func (gen *uuidGeneratorUnsafe) Gen() UUID {
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

func (gen *uuidGeneratorUnsafe) String() string {
	return fmt.Sprintf("requester[%v] sn[%v] timestamp[%v]", gen.requester, gen.sn, gen.timestamp)
}
