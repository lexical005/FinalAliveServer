package main

import (
	"time"
)

const (
	// worldTimeUpdateInterval 游戏世界更新间隔
	worldTimeUpdateInterval = 1 * time.Millisecond

	// uuid.NewGenerator的requester参数的格式：4位类型+12位服务器编号
	uuidRequesterBitType     = 4
	uuidRequesterBitServerID = 12
)
