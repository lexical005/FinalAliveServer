package ffProto

import (
	"ffCommon/pool"
	"fmt"

	"github.com/golang/protobuf/proto"
)

var msgPool []*pool.Pool

func init() {
	maxProtoID := len(listProtoID)
	for _, v := range listProtoID {
		if maxProtoID < int(v) {
			maxProtoID = int(v)
		}
	}

	msgPool = make([]*pool.Pool, maxProtoID+1, maxProtoID+1)
	for i, protoID := range listProtoID {
		msgPool[protoID] = pool.New(fmt.Sprintf("ffProto.message_pool.msgPool_%d", i), true, mapMessageCreator[protoID], 10, 50)
	}
}

func applyMessage(mt MessageType) proto.Message {
	m, _ := msgPool[int32(mt)].Apply().(proto.Message)
	// todo: 增加追踪-添加到被使用列表
	return m
}

func backMessage(mt MessageType, m proto.Message) {
	// todo: 增加追踪-从被使用列表移除，供缓存泄露分析
	m.Reset()
	msgPool[int32(mt)].Back(m)
}
