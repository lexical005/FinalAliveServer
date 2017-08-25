package ffProto

import (
	"ffCommon/uuid"
)

// ISendProtoExtraDataNormal 发送附加数据类型为ExtraDataNormal的协议
type ISendProtoExtraDataNormal interface {
	SendProtoExtraDataNormal(proto *Proto) bool
}

// ISendProtoExtraDataUUID 发送附加数据类型为ExtraDataUUID的协议
type ISendProtoExtraDataUUID interface {
	SendProtoExtraDataUUID(uuid uuid.UUID, proto *Proto) bool
}

// SendProtoExtraDataNormal 发送附加数据类型为ExtraDataNormal的协议
//	sender: 实现了ISendProtoExtraDataNormal接口的对象
//	proto: 待发送的协议
//	isProtoRecved: 被发送的协议, 是不是接收到的协议再被用来发送的
func SendProtoExtraDataNormal(sender ISendProtoExtraDataNormal, proto *Proto, isProtoRecved bool) bool {
	if isProtoRecved {
		proto.ChangeLimitStateRecvToSend()
	}
	sender.SendProtoExtraDataNormal(proto)
	return isProtoRecved
}

// SendProtoExtraDataUUID 发送附加数据类型为ExtraDataNormal的协议
//	sender: 实现了ISendProtoExtraDataNormal接口的对象
//	uuid: 附加数据
//	proto: 待发送的协议
//	isProtoRecved: 被发送的协议, 是不是接收到的协议再被用来发送的
func SendProtoExtraDataUUID(sender ISendProtoExtraDataUUID, uuid uuid.UUID, proto *Proto, isProtoRecved bool) bool {
	if isProtoRecved {
		proto.ChangeLimitStateRecvToSend()
	}
	sender.SendProtoExtraDataUUID(uuid, proto)
	return isProtoRecved
}
