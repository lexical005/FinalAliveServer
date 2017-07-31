package base

import (
	"ffCommon/uuid"
	"ffProto"
)

// ServerNetEventData 向外界通知的事件数据, 事件处理完毕后, 必须执行Back方法, 以回收所有相关资源
type ServerNetEventData interface {
	// NetEventData 将NetEventData定义组合进来
	NetEventData

	// SessionUUID 事件关联的session的UUID, 当NetEventType为NetEventEnd时无效
	SessionUUID() uuid.UUID

	// Server 事件关联的server
	Server() Server
}

// Server Server自身, 未实现多goroutine安全, 由使用者确保
type Server interface {
	// Start 启动 Server 监听 Client 链接
	Start(chNetEventData chan NetEventData, recvProtoExtraDataType ffProto.ExtraDataType) error

	// SendProto 发送Proto到指定对端, 异步
	SendProto(uuidSession uuid.UUID, p *ffProto.Proto) error

	// Close 关闭, 一次性, 外界以处理到NetEventEnd作为Session结束的最后一个事件, 异步
	// 内部确保只有首次调用时有效
	//	uuidSession: 要关闭的连接
	// 	delayMillisecond: 延迟多少毫秒关闭
	CloseSession(uuidSession uuid.UUID, delayMillisecond int64)

	// Close 关闭, 一次性, 外界以处理到NetEventEnd作为Server结束的最后一个事件, 异步
	// 内部确保只有首次调用时有效
	// 	delayMillisecond: 延迟多少毫秒关闭
	Close(delayMillisecond int64)

	String() string
}
