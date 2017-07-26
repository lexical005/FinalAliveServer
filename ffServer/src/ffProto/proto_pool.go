package ffProto

import (
	"ffCommon/pool"

	"github.com/golang/protobuf/proto"

	"fmt"
)

var arraySizeLimit = []int{16, 32, 64, 128, 256, 512, 1024, 2048}
var arrayCountRate = []int{50, 50, 50, 50, 25, 20, 10, 5}
var protoPool []*pool.Pool

func init() {
	l1, l2 := len(arraySizeLimit), len(arrayCountRate)
	if l1 != l2 {
		panic(fmt.Sprintf("ffProto.proto_pool len(arraySizeLimit)[%d] != len(arrayCountRate)[%d]", l1, l2))
	}

	protoPool = make([]*pool.Pool, len(arrayCountRate), len(arrayCountRate)*3/2)
	for i, v := range arrayCountRate {
		protoPool[i] = newProtoPool(arraySizeLimit[i], v)
	}
}

func newProtoPool(bufSize, countRate int) *pool.Pool {
	creator := func() interface{} {
		buf := make([]byte, 0, bufSize)
		p := &Proto{
			pb:        proto.NewBuffer(nil),
			extraData: make([]byte, extraDataMaxLength, extraDataMaxLength),
		}
		p.setBuf(buf)
		return p
	}

	return pool.New(fmt.Sprintf("ffProto.proto_pool.protoPool_%v", bufSize), true, creator, initProtoCount*countRate/100, 50)
}

// bufLengthLimit: >=0 限定协议体缓冲区大小; <0 协议体缓冲区越大越好
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
	// 	arrayCountRate = append(arrayCountRate, arrayCountRate[l-1])

	// 	protoPool = append(protoPool, newProtoPool(bufLengthLimit, arrayCountRate[l-1]))
	// }

	index := findIndex(bufLengthLimit)
	p, _ = protoPool[index].Apply().(*Proto)
	if p.Cap() < bufLengthLimit {
		p.setBuf(make([]byte, 0, bufLengthLimit))
	}
	p.resetForRecv(header, bufLengthLimit)
	return p
}

// ApplyProtoForSend apply a Proto with specified Message
func ApplyProtoForSend(protoID MessageType) (p *Proto) {
	index := findIndex(-1)
	p, _ = protoPool[index].Apply().(*Proto)
	p.resetForSend(protoID)
	return p
}

// BackProtoAfterSend back Proto to pool after send
func BackProtoAfterSend(p *Proto) {
	p.onBackPoolAfterSend()

	index := findIndex(p.Cap())
	protoPool[index].Back(p)
}

// BackProtoAfterRecv back Proto to pool after recv
func BackProtoAfterRecv(p *Proto) {
	p.onBackPoolAfterRecv()

	index := findIndex(p.Cap())
	protoPool[index].Back(p)
}

// BackProtoAfterDispatch back Proto to pool after cache
func BackProtoAfterDispatch(p *Proto) {
	p.onBackPoolAfterDispatch()

	index := findIndex(p.Cap())
	protoPool[index].Back(p)
}

// ForceBackProtoInWaitSend force back Proto to pool in state useStateCacheWaitSend
func ForceBackProtoInWaitSend(p *Proto) {
	p.forceBackProtoInWaitSend()

	index := findIndex(p.Cap())
	protoPool[index].Back(p)
}
