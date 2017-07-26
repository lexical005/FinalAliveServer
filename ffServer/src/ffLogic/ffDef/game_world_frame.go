package ffDef

import (
	"ffAutoGen/ffError"
	"ffCommon/uuid"
	"ffProto"
	"time"
)

// UUIDType UUID 使用分类，最多支持16种
type UUIDType byte

const (
	// UUIDTimer 定时器
	UUIDTimer UUIDType = 0

	// UUIDEquipment 装备实例
	UUIDEquipment UUIDType = 1

	// UUIDTypeCount UUIDType 数量
	UUIDTypeCount = 2
)

// ITimer 定时器接口
type ITimer interface {
	// FirstOffset 首次触发偏移，单位：毫秒
	FirstOffset() time.Duration

	// Interval 触发间隔，单位：毫秒
	Interval() time.Duration

	// LeftCount 剩余触发次数. -1 代表不限次数
	LeftCount() int

	String() string
}

// TimerFunc 定时器触发时回调函数
type TimerFunc func(ITimer)

// ProtoCallback 协议响应回调函数
type ProtoCallback func(p *ffProto.Proto, uuidAgent uuid.UUID)

// IGameWorldFrame 游戏世界框架支持
type IGameWorldFrame interface {
	// DefaultOnlineCount 初始默认在线人数限定。multi goroutine safe
	// 返回值: 初始默认在线人数限定
	DefaultOnlineCount() int

	// Kick 踢出。multi goroutine safe
	//  uuidAgent: 连接标识
	//	notifyKick: 是否发送协议进行通知
	//	kickReason: 踢出原因
	Kick(uuidAgent uuid.UUID, notifyKick bool, kickReason ffError.Error)

	// SendProto 发送协议。multi goroutine safe
	//  p: 待发送协议
	//  uuidAgent: 连接标识
	// 返回值: 连接是否存在
	SendProto(uuidAgent uuid.UUID, p *ffProto.Proto) bool

	// AddTimer 设置定时器, 单位: 毫秒。只允许在主循环驱动的事件里调用。非multi goroutine safe
	//  firstOffset: 首次触发偏移
	//  interval: 触发间隔
	//  count: 触发次数. -1 代表不限次数. >0: 代表指定限定次数
	//  timerFunc: 触发时回调函数
	// 返回值: 定时器对象 ITimer
	AddTimer(firstOffset time.Duration, interval time.Duration, count int, timerFunc TimerFunc) ITimer

	// StopTimer 停止定时器。只允许在主循环驱动的事件里调用。非multi goroutine safe
	//  timer: 定时器对象
	StopTimer(timer ITimer)

	// UUID UUID 生成唯一标识。只允许在主循环驱动的事件里调用。非multi goroutine safe
	//  uuidType: UUID 使用分类
	// 返回值: UUID 唯一标识
	UUID(uuidType UUIDType) uuid.UUID

	// DBQuery 返回一个新的 IDBQueryRequest，用于数据库操作。callback和args，在请求者主动取消前或者接收到查询结果前，必须有效！
	//  idMysqlDB   int             // 哪个数据库
	//  idMysqlStmt int             // 哪个语句
	//  callback    DBQueryCallback // 查询结果回调函数
	//  args        []interface{}   // 查询参数
	DBQuery(idMysqlDB int, idMysqlStmt int, callback DBQueryCallback, args ...interface{}) IDBQueryRequest
}
