package ffDef

import (
	"ffAutoGen/ffError"
	"ffCommon/uuid"
	"ffProto"
)

// IGameWorld 游戏世界接口
type IGameWorld interface {
	// Start 启动
	Start()

	// Stop 停止
	Stop()

	// DispatchProto 处理协议
	//  uuidAgent: 连接唯一标识
	//  p: 待处理协议
	DispatchProto(uuidAgent uuid.UUID, p *ffProto.Proto)

	// KickAll 所有人全部立即下线
	//	kickReason: 踢出原因
	KickAll(kickReason ffError.Error)
}
