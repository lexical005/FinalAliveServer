package base

import (
	"ffProto"
)

// Client 主动连接服务端, 尝试建立连接, 进行通讯
type Client interface {
	// Start 开始连接Server, 只执行一次, 异步
	//	chNetEventData: 事件数据管道, 仅有写入权
	//	recvProtoExtraDataType: 接收的协议的附加数据类型
	Start(chNetEventData chan NetEventData, recvProtoExtraDataType ffProto.ExtraDataType) error

	// SendProto 发送Proto到对端, 只应该在收到连接建立事件之后再调用, 异步
	SendProto(p *ffProto.Proto)

	// Close 关闭, 一次性, 外界以处理到NetEventEnd作为Client结束的最后一个事件, 异步
	// 内部确保只有首次调用时有效
	// 	delayMillisecond: 延迟多少毫秒关闭
	Close(delayMillisecond int64)

	String() string
}
