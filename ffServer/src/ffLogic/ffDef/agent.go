package ffDef

import (
	"ffAutoGen/ffError"
	"ffCommon/uuid"
	"ffProto"
)

// IAccountAgent IAccount 对应的在线代理
type IAccountAgent interface {
	// UUID agent connection UUID
	UUID() uuid.UUID

	// SendProto send proto
	SendProto(p *ffProto.Proto)

	// Kick kick agent
	Kick(notifyKick bool, kickReason ffError.Error)
}
