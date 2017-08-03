package base

import "ffCommon/uuid"

// Client 主动连接服务端, 尝试建立连接, 进行通讯
type Client interface {
	// Start 开始连接Server, 只执行一次, 异步
	//	chNewSession: 外界接收新连接被创建事件的管道, Client仅有写入权
	//	chClientClosed: 用于向外界通知关闭完成的管道, Client仅有写入权
	Start(chNewSession chan Session, chClientClosed chan struct{})

	// Stop 停止连接
	Stop()

	// Back 回收Client资源, 只应在连接已完成关闭情况下执行, 一次性, 同步
	Back()

	// UUID 唯一标识
	UUID() uuid.UUID

	String() string
}
