package session

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/pool"
	"ffCommon/uuid"
	"fmt"
)

// 没有网络包进入的最大间隔(秒)
var readDeadtime int

// Session缓存
var sp *sessionPool

// Init 初始session模块
// _readDeadtime: 没有网络包进入的最大间隔(秒)，超过此时限，则认为对端关闭了
// _onlineCount: 预计多少Session同时连接
func Init(_readDeadtime int, _onlineCount int) (err error) {
	if _readDeadtime < 30 {
		return fmt.Errorf("session.Init: invalid _readDeadtime[%v]", _readDeadtime)
	}

	if _onlineCount < 1 {
		return fmt.Errorf("session.Init: invalid _onlineCount[%v]", _onlineCount)
	}

	uuidGenerator, err := uuid.NewGenerator(0)
	if err != nil {
		return err
	}

	readDeadtime = _readDeadtime

	creator := func() interface{} {
		return newSession()
	}

	sp = &sessionPool{
		pool: pool.New("session.sp.pool", false, creator, _onlineCount, 50),

		uuidGenerator: uuidGenerator,
	}

	return nil
}

// Apply idle session, and session will be back automatically on lower conn closed
func Apply() (s base.Session) {
	return sp.apply()
}

// PrintModule 输出Session模块信息
func PrintModule() {
	runLogger := log.RunLogger
	runLogger.Println("PrintModule Start session:")

	runLogger.Println("session.sp.pool:")
	runLogger.Println(sp)

	runLogger.Printf("PrintModule End session\n\n")
}
