package ffProto

import (
	"ffCommon/pool"

	"github.com/golang/protobuf/proto"
)

// protoInitCount 初始缓存多少个Proto以供外界使用
var protoInitCount = 64

var protoPool *pool.Pool

func init() {
	protoPool = newProtoPool(protoInitCount)
}

func newProtoPool(initCount int) *pool.Pool {
	creator := func() interface{} {
		p := &Proto{
			pb:        proto.NewBuffer(nil),
			extraData: make([]byte, extraDataMaxLength, extraDataMaxLength),
		}
		p.setBuf(nil)
		return p
	}

	return pool.New("ffProto.proto_pool.protoPool", false, creator, initCount, 50)
}

// ApplyProtoForRecv apply a Proto from pool for recv
func ApplyProtoForRecv(header *ProtoHeader) (p *Proto) {
	bufLengthLimit := protoHeaderLength + header.contentLength + header.recvExtraDataType.BufferLength()
	buf := applyBuffer(bufLengthLimit)

	p, _ = protoPool.Apply().(*Proto)
	p.setBuf(buf)
	p.resetForRecv(header, bufLengthLimit)
	return p
}

// ApplyProtoForSend apply a Proto for send specified Message
func ApplyProtoForSend(protoID MessageType) (p *Proto) {
	// todo: 解析协议结构，尽可能得到确切的大小
	index := findBufferIndex(-1)
	p, _ = protoPool.Apply().(*Proto)
	buf, _ := bufferPool[index].Apply().([]byte)
	p.setBuf(buf)
	p.resetForSend(protoID)
	return p
}

func backProto(p *Proto) {
	protoPool.Back(p)
}
