package base

import (
	"ffCommon/uuid"
	"ffProto"
)

// Session 封装网络连接, 内部开启接收Proto协程和发送Proto协程
//	如果在Start之前就要关闭连接, 则执行Close, 否则, 由外界关闭chSendProto来引发退出
type Session interface {
	// Start 启动 Session 主循环, 一次性, 异步
	//	chSendProto: 外界发送协议的管道, Session仅有读取权. 当从该管道读取到nil时, 即认为外界要求关闭Session
	//	chNetEventData: 外界接收Session反馈的事件数据管道, Session仅有写入权
	//	recvProtoExtraDataType: 接收的协议的附加数据类型
	Start(chSendProto chan *ffProto.Proto, chNetEventData chan NetEventData, recvProtoExtraDataType ffProto.ExtraDataType)

	// Close 外界决定关闭Client或Server后, 通过此方法, 关闭chNewSession内尚未处理的新建立的连接. 必须发生在Start之前, 一次性, 同步.
	// 如果在Start之前就要关闭连接, 则执行Close, 否则, 由外界关闭chSendProto来引发退出
	Close()

	// UUID 唯一标识
	UUID() uuid.UUID

	String() string
}
