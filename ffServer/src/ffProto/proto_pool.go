package ffProto

import (
	"ffCommon/pool"

	"github.com/golang/protobuf/proto"

	"fmt"
)

var arraySizeLimit = []int{16, 32, 64, 128, 256, 512, 1024, 2048, 4096}
var arrayInitCount = []int{100, 100, 100, 100, 50, 50, 20, 10, 5}
var protoPool []*pool.Pool

func init() {
	l1, l2 := len(arraySizeLimit), len(arrayInitCount)
	if l1 != l2 {
		panic(fmt.Sprintf("ffProto.proto_pool len(arraySizeLimit)[%d] != len(arrayInitCount)[%d]", l1, l2))
	}

	for index := 1; index < l1; index++ {
		if arraySizeLimit[index-1] >= arraySizeLimit[index] {
			panic(fmt.Sprintf("ffProto.proto_pool invalid arraySizeLimit at index %d:%d: size must increase", index-1, index))
		} else if arraySizeLimit[index] > protoMaxLength {
			panic(fmt.Sprintf("ffProto.proto_pool invalid arraySizeLimit: size must letter than protoMaxLength[%d]", protoMaxLength))
		}

		if arrayInitCount[index-1] < arrayInitCount[index] {
			panic(fmt.Sprintf("ffProto.proto_pool invalid arrayInitCount at index %d:%d: init count must decrease", index-1, index))
		}
	}

	protoPool = make([]*pool.Pool, l1, l1)
	for i, v := range arrayInitCount {
		protoPool[i] = newProtoPool(arraySizeLimit[i], v)
	}
}

func newProtoPool(bufSize, initCount int) *pool.Pool {
	creator := func() interface{} {
		buf := make([]byte, 0, bufSize)
		p := &Proto{
			pb:        proto.NewBuffer(nil),
			extraData: make([]byte, extraDataMaxLength, extraDataMaxLength),
		}
		p.setBuf(buf)
		return p
	}

	return pool.New(fmt.Sprintf("ffProto.proto_pool.protoPool_%v", bufSize), true, creator, initCount, 50)
}

// bufLengthLimit: >=0 限定协议缓冲区大小; <0 协议缓冲区越大越好
func findIndex(bufLengthLimit int) int {
	index := len(arraySizeLimit) - 1
	if bufLengthLimit >= 0 {
		for i, v := range arraySizeLimit {
			if bufLengthLimit <= v {
				index = i
				break
			}
		}
	}
	return index
}

// ApplyProtoForRecv apply a Proto from pool for recv
func ApplyProtoForRecv(header *ProtoHeader) (p *Proto) {
	bufLengthLimit := protoHeaderLength + header.contentLength + getExtraDataLength(ExtraDataType(header.recvExtraDataType))

	// // 暂不考虑动态新增更大的缓存
	// // 需求的大小, 超出当前已有的缓存
	// l := len(arraySizeLimit)
	// if bufLengthLimit > arraySizeLimit[l-1] {
	// 	arraySizeLimit = append(arraySizeLimit, bufLengthLimit)
	// 	arrayInitCount = append(arrayInitCount, arrayInitCount[l-1])

	// 	protoPool = append(protoPool, newProtoPool(bufLengthLimit, arrayInitCount[l-1]))
	// }

	index := findIndex(bufLengthLimit)
	p, _ = protoPool[index].Apply().(*Proto)
	p.resetForRecv(header, bufLengthLimit)
	return p
}

// ApplyProtoForSend apply a Proto for specified Message
func ApplyProtoForSend(protoID MessageType) (p *Proto) {
	// todo: 解析协议结构，尽可能得到确切的大小
	index := findIndex(-1)
	p, _ = protoPool[index].Apply().(*Proto)
	p.resetForSend(protoID)
	return p
}

// BackProtoAfterSend back Proto to pool after send
func BackProtoAfterSend(p *Proto) {
	if p.onBackPoolAfterSend() {
		backProto(p)
	}
}

// BackProtoAfterRecv back Proto to pool after recv
func BackProtoAfterRecv(p *Proto) {
	if p.onBackPoolAfterRecv() {
		backProto(p)
	}
}

// BackProtoAfterDispatch back Proto to pool after dispatch
func BackProtoAfterDispatch(p *Proto) {
	if p.onBackPoolAfterDispatch() {
		backProto(p)
	}
}

// BackProtoInWaitSend force back Proto to pool in state useStateCacheWaitSend
func BackProtoInWaitSend(p *Proto) {
	if p.onBackProtoInWaitSend() {
		backProto(p)
	}
}

func backProto(p *Proto) {
	index := findIndex(p.Cap())
	protoPool[index].Back(p)
}
