package netmanager

import (
	"ffCommon/net/base"
	"ffCommon/uuid"
	"ffProto"
)

type inet interface {
	// onAgentClosed 连接断开事件
	onAgentClosed()

	// Stop 通知底层停止
	Stop()

	// BackNet 回收底层
	BackNet()

	// NewSessionChan 返回新连接管道
	NewSessionChan() chan base.Session

	// WaitNetExit 等待底层退出, 直到底层抛出退出事件后, 才返回
	WaitNetExit()

	// UUID 底层唯一标识
	UUID() uuid.UUID

	// Clear 清理
	Clear()

	// SendExtraDataType 发送的协议的附加数据类型
	SendExtraDataType() ffProto.ExtraDataType

	// RecvExtraDataType 发送的协议的附加数据类型
	RecvExtraDataType() ffProto.ExtraDataType

	// SessionNetEventDataCache 网络事件管道的缓存大小. 影响处理网络事件的速度.
	SessionNetEventDataCache() int

	// SessionSendProtoCache 待发送协议管道的缓存大小. 影响发送协议的速度
	SessionSendProtoCache() int
}
