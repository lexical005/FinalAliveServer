package main

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffLogic/ffDef"
	"ffProto"

	"time"
)

// 特殊事件
type specialEvent byte

const (
	// 特殊事件--与AgentServer连接断开, 踢所有人下线
	specialEventKickAll specialEvent = 0

	// 特殊事件--关闭
	specialEventClose specialEvent = 1
)

type gameWorldFrame struct {
	// 待处理的接收到的协议
	chRecvProtos chan *ffProto.Proto

	// 特殊事件
	chSpecialEvent chan specialEvent

	// 模块
	timeManager *gameWorldTimeManager
	uuidGen     *gameWorldUUIDGen
}

// DefaultOnlineCount 初始默认在线人数限定
// 返回值: 初始默认在线人数限定。multi goroutine safe
func (gwf *gameWorldFrame) DefaultOnlineCount() int {
	return appConfig.Session.OnlineCount
}

// Kick 踢出。multi goroutine safe
//  uuidAgent: 连接标识
//	notify: 是否发送协议进行通知
//	kickReason: 踢出原因
func (gwf *gameWorldFrame) Kick(uuidAgent uuid.UUID, notifyKick bool, kickReason ffError.Error) {
	log.RunLogger.Printf("gameWorldFrame.Kick: uuidAgent[%x] notifyKick[%v] kickReason[%v]",
		uuidAgent, notifyKick, kickReason)

	// 踢出通知
	if notifyKick {
		p := ffProto.ApplyProtoForSend(ffProto.MessageType_Kick)
		message := p.Message().(*ffProto.MsgKick)
		message.Result = kickReason.Code()
		gwf.SendProto(uuidAgent, p)
	}
}

// SendProto 发送协议。multi goroutine safe
//  uuidAgent: 连接标识
//  p: 待发送协议
// 返回值: 连接是否存在
func (gwf *gameWorldFrame) SendProto(uuidAgent uuid.UUID, p *ffProto.Proto) bool {
	return agentServerMgr.sendProto(p, uint64(uuidAgent))
}

// AddTimer 添加定时器, 单位: 毫秒。只允许在主循环驱动的事件里调用。非multi goroutine safe
//  firstOffset: 首次触发偏移
//  interval: 触发间隔
//  count: 触发次数. -1 代表不限次数. >0: 代表指定限定次数
//  timerFunc: 触发时回调函数
// 返回值: 定时器对象 ITimer
func (gwf *gameWorldFrame) AddTimer(firstOffset time.Duration, interval time.Duration, count int, timerFunc ffDef.TimerFunc) ffDef.ITimer {
	return gwf.timeManager.AddTimer(firstOffset, interval, count, timerFunc)
}

// StopTimer 停止定时器。只允许在主循环驱动的事件里调用。非multi goroutine safe
//  timer: 定时器对象
// multi goroutine safe
func (gwf *gameWorldFrame) StopTimer(timer ffDef.ITimer) {
	gwf.timeManager.StopTimer(timer)
}

// UUID UUID 生成唯一标识。只允许在主循环驱动的事件里调用。非multi goroutine safe
//  uuidType: UUID 使用分类
// 返回值: UUID 唯一标识
func (gwf *gameWorldFrame) UUID(uuidType ffDef.UUIDType) uuid.UUID {
	return uuid.InvalidUUID
}

// DBQuery 返回一个新的 IDBQueryRequest，用于数据库操作。callback和args，在请求者主动取消前或者接收到查询结果前，必须有效！
//  idMysqlDB   int             // 哪个数据库
//  idMysqlStmt int             // 哪个语句
//  callback    DBQueryCallback // 查询结果回调函数
//  args        []interface{}   // 查询参数
func (gwf *gameWorldFrame) DBQuery(idMysqlDB int, idMysqlStmt int, callback ffDef.DBQueryCallback, args ...interface{}) ffDef.IDBQueryRequest {
	return nil
}

// init 初始化
func (gwf *gameWorldFrame) init() {
	// 缓存的等待处理的协议数目
	cacheProtoCount := appConfig.Session.OnlineCount * 20 / 100
	if cacheProtoCount < 10 {
		cacheProtoCount = 10
	}
	gwf.chRecvProtos = make(chan *ffProto.Proto, cacheProtoCount)

	// 特殊事件
	gwf.chSpecialEvent = make(chan specialEvent, 1)

	// 管理器
	gwf.timeManager = &gameWorldTimeManager{}
	gwf.timeManager.init()

	gwf.uuidGen = &gameWorldUUIDGen{}
	gwf.uuidGen.init()
}

// onRecvProto 接收到协议
func (gwf *gameWorldFrame) onRecvProto(p *ffProto.Proto) {
	// 缓存以待分发
	p.SetCacheWaitDispatch()
	gwf.chRecvProtos <- p
}

// dispatchProto 处理接收到协议
func (gwf *gameWorldFrame) dispatchProto(p *ffProto.Proto) {
	// 分发完毕后，尝试回收协议
	defer ffProto.BackProtoAfterDispatch(p)

	// 由具体逻辑处理协议
	world.DispatchProto(uuid.UUID(p.ExtraData()), p)
}

// onAgentServerDisConnect 与AgentServer断开连接
func (gwf *gameWorldFrame) onAgentServerDisConnect() {
	gwf.chSpecialEvent <- specialEventKickAll
}

// dispatchSpecialEvent 处理特殊事件
// 返回值: 是否结束主循环
func (gwf *gameWorldFrame) dispatchSpecialEvent(specialEvent specialEvent) bool {
	if specialEvent == specialEventKickAll {
		world.KickAll(ffError.ErrKickConnection)
		return false
	}

	if specialEvent == specialEventClose {
		return true
	}

	return false
}

// mainLoop 开启游戏世界主循环
func (gwf *gameWorldFrame) mainLoop(params ...interface{}) {
	// 游戏世界启动
	world.Start()

	// 事件/时间驱动
deadLoop:
	for {
		select {
		case <-time.After(worldTimeUpdateInterval): // 时间驱动
			// 主循环间隔

		case p := <-gwf.chRecvProtos: // 接收到的协议驱动
			gwf.dispatchProto(p)

		case specialEvent := <-gwf.chSpecialEvent: // 特殊事件
			if gwf.dispatchSpecialEvent(specialEvent) {
				break deadLoop
			}

			// 数据库查询驱动
		}
	}
}
