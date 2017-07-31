package base

import (
	"ffCommon/uuid"
	"ffProto"

	"net"
)

// SessionNetEventData 向外界通知的事件数据, 事件处理完毕后, 必须执行Back方法, 以回收所有相关资源
type SessionNetEventData interface {
	// NetEventData 将NetEventData定义组合进来
	NetEventData

	// Session 事件关联的session
	Session() Session
}

// Session 多goroutine安全
type Session interface {
	// Start 启动 Session 主循环, 一次性, 异步
	//	conn: 底层连接
	//	chNetEventData: 事件数据管道, 仅有写入权
	//	recvProtoExtraDataType: 接收的协议的附加数据类型
	Start(conn net.Conn, chNetEventData chan NetEventData, recvProtoExtraDataType ffProto.ExtraDataType)

	// UUID 唯一标识
	UUID() uuid.UUID

	// SendProto 发送Proto到对端, 外界只应该在收到连接建立完成事件之后再调用此接口, 多goroutine安全, 异步
	SendProto(p *ffProto.Proto)

	// Close 关闭, 一次性, 外界以处理到NetEventEnd作为Session结束的最后一个事件, 异步
	// 内部确保只有首次调用时有效
	// 	delayMillisecond: 延迟多少毫秒关闭
	Close(delayMillisecond int64)

	String() string
}
