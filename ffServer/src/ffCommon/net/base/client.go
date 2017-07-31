package base

import (
	"ffProto"
)

// ClientNetEventData 向外界通知的事件数据, 事件处理完毕后, 必须执行Back方法, 以回收所有相关资源
type ClientNetEventData interface {
	// NetEventData 将NetEventData定义组合进来
	NetEventData

	// Client Client
	Client() Client
}

// Client Client自身, 未实现多goroutine安全, 由使用者确保
type Client interface {
	// Start 开始连接Server, 只执行一次, 异步
	//	chNetEventData: 事件数据管道, 仅有写入权
	//	recvProtoExtraDataType: 接收的协议的附加数据类型
	Start(chNetEventData chan NetEventData, recvProtoExtraDataType ffProto.ExtraDataType) error

	// ReConnect 重连, 只能在外界处理到连接断开事件时, 如果需要恢复连接(比如不能在执行了Close之后再尝试ReConnect), 才可调用此接口, 异步
	ReConnect()

	// SendProto 发送Proto到对端, 只应该在收到连接建立事件之后再调用, 异步
	SendProto(p *ffProto.Proto)

	// Close 关闭, 一次性, 外界以处理到NetEventEnd作为Client结束的最后一个事件, 异步
	// 内部确保只有首次调用时有效
	// 	delayMillisecond: 延迟多少毫秒关闭
	Close(delayMillisecond int64)

	String() string
}
