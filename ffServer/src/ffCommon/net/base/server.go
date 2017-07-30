package base

import (
	"ffProto"
)

// Server Server自身, 未实现多goroutine安全, 由使用者确保
type Server interface {
	// 关闭
	Close()

	// 启动 Server 监听 Client 链接
	Start(chNetEventData chan NetEventData, recvProtoExtraDataType ffProto.ExtraDataType) error

	String() string
}
