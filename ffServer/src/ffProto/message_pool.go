package ffProto

import (
	"ffCommon/pool"
	"fmt"

	"github.com/golang/protobuf/proto"
)

var msgPool []*pool.Pool

func init() {
	maxID := len(listProtoID)
	for _, v := range listProtoID {
		if maxID < int(v) {
			maxID = int(v)
		}
	}

	msgPool = make([]*pool.Pool, maxID+1, maxID+1)
	for i, protoID := range listProtoID {
		msgPool[protoID] = pool.New(fmt.Sprintf("ffProto.message_pool.msgPool_%d", i), true, mapMessageCreator[protoID], 10, 50)
	}
}

func applyMessage(mt MessageType) proto.Message {
	m, _ := msgPool[int32(mt)].Apply().(proto.Message)
	return m
}

func backMessage(mt MessageType, m proto.Message) {
	m.Reset()
	msgPool[int32(mt)].Back(m)
}
