package base

import (
	"ffCommon/uuid"
	"ffProto"

	"net"
)

// Session 封装网络连接, 内部开启接收Proto协程和发送Proto协程
//	如果在Start之前就要关闭连接, 则执行Close, 否则, 由外界关闭chSendProto来引发退出
type Session interface {
	// SetConn 设置底层连接, 第一优先调用, 同步
	SetConn(conn net.Conn)

	// Start 启动 Session 主循环, 一次性, 异步
	//	conn: 底层连接
	//	chSendProto: 外界发送协议的管道, Session仅有读取权. 当该管道关闭时, 即认为外界主动关闭Session
	//	chNetEventData: 外界接收Session反馈的事件数据管道, Session仅有写入权
	//	recvProtoExtraDataType: 接收的协议的附加数据类型
	Start(chSendProto chan *ffProto.Proto, chNetEventData chan NetEventData, recvProtoExtraDataType ffProto.ExtraDataType)

	// Close 在执行Start之前, 就直接关闭连接, 用于外界已决定关闭服务时新建立的连接需要立即关闭, 一次性, 同步
	Close()

	// UUID 唯一标识
	UUID() uuid.UUID

	String() string
}
