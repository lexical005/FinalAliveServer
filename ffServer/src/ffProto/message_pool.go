package ffProto

import (
	"ffCommon/pool"
	"fmt"

	"github.com/golang/protobuf/proto"
)

var messagePool []*pool.Pool

func init() {
	maxProtoID := int(listMessageID[len(listMessageID)-1])
	messagePool = make([]*pool.Pool, maxProtoID+1, maxProtoID+1)
	for _, protoID := range listMessageID {
		messagePool[protoID] = pool.New(
			fmt.Sprintf("ffProto.message_pool.messagePool_%v", protoID),
			false,
			mapMessageCreator[protoID],
			10,
			50)
	}
}

func applyMessage(mt MessageType) proto.Message {
	m, _ := messagePool[int32(mt)].Apply().(proto.Message)
	return m
}

func backMessage(mt MessageType, m proto.Message) {
	m.Reset()
	messagePool[int32(mt)].Back(m)
}
