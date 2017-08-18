package ffGameConfig

import (
	"ffCommon/log/log"

	"fmt"
)

// Error Error
type Error interface {
	Code() int32
	Error() string
	String() string
}

type errReason struct {
	code int32
	desc string
}

func (ec *errReason) Code() int32 {
	return ec.code
}

func (ec *errReason) Error() string {
	return fmt.Sprintf("ffGameConfig[%d-%s]", ec.code, ec.desc)
}

func (ec *errReason) String() string {
	return fmt.Sprintf("ffGameConfig[%d-%s]", ec.code, ec.desc)
}

// ErrNone 无错误
var ErrNone Error = &errReason{code: 0, desc: "无错误"}

// ErrUnknown 未知错误
var ErrUnknown Error = &errReason{code: 1, desc: "未知错误"}

// ErrInvalidParam 无效参数
var ErrInvalidParam Error = &errReason{code: 2, desc: "无效参数"}

// ErrAccountLevelLess 账号等级不足
var ErrAccountLevelLess Error = &errReason{code: 3, desc: "账号等级不足"}

// ErrAccountVitalityLess 账号体力不足
var ErrAccountVitalityLess Error = &errReason{code: 4, desc: "账号体力不足"}

// ErrHeroLevelLess 英雄等级不足
var ErrHeroLevelLess Error = &errReason{code: 5, desc: "英雄等级不足"}

// ErrTemplateItemLess 物品数量不足
var ErrTemplateItemLess Error = &errReason{code: 6, desc: "物品数量不足"}

// ErrTemplateItemTooMuch 持有物品数量太多
var ErrTemplateItemTooMuch Error = &errReason{code: 7, desc: "持有物品数量太多"}

// ErrGameServerOffline 目标服务器不在线
var ErrGameServerOffline Error = &errReason{code: 8, desc: "目标服务器不在线"}

// ErrKickConnection 连接异常
var ErrKickConnection Error = &errReason{code: 9, desc: "连接异常"}

// ErrKickProtoInvalid 连接异常
var ErrKickProtoInvalid Error = &errReason{code: 10, desc: "连接异常"}

// ErrKickExotic 帐号异地登录, 您被迫离线
var ErrKickExotic Error = &errReason{code: 11, desc: "帐号异地登录, 您被迫离线"}

var errByCode = []Error{
	ErrNone,
	ErrUnknown,
	ErrInvalidParam,
	ErrAccountLevelLess,
	ErrAccountVitalityLess,
	ErrHeroLevelLess,
	ErrTemplateItemLess,
	ErrTemplateItemTooMuch,
	ErrGameServerOffline,
	ErrKickConnection,
	ErrKickProtoInvalid,
	ErrKickExotic,
}

// ErrByCode 根据错误码获取Error
func ErrByCode(errCode int32) Error {
	if errCode >= 0 && int(errCode) < len(errByCode) {
		return errByCode[errCode]
	}

	log.FatalLogger.Printf("ffGameConfig.ErrByCode: invalid errCode[%d]", errCode)

	return ErrUnknown
}
