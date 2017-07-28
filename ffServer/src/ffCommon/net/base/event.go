package base

import "ffProto"

// NetEventType Net向外界通知的事件类型
type NetEventType byte

const (
	// NetEventInvalid Net不会向外界抛出此事件, 其只作为内部使用
	NetEventInvalid NetEventType = 0

	// NetEventOn Net连接建立完成事件
	//	此事件之后, Net将可用, 此后方可执行发送协议请求
	NetEventOn NetEventType = 1

	// NetEventOff Net连接已断开事件
	//	此事件之后, Net将不可发送Proto
	NetEventOff NetEventType = 2

	// NetEventProto Net向外界抛出收到Proto事件
	//	在外界主动关闭Net请求后, 依然可能收到Proto, 直到外界处理完毕NetEventOff事件
	NetEventProto NetEventType = 3

	// NetEventEnd Net最后一个事件
	//	此事件之后, Net将不可用, 外界不应再对Net进行引用
	NetEventEnd NetEventType = 4
)

// NetEventData 向外界通知的事件数据, 事件处理完毕后, 必须执行Back方法, 以回收所有相关资源
type NetEventData interface {
	// NetEventType 获取事件类型
	NetEventType() NetEventType

	// ManualClose 当NetEventType为NetEventOff时有效, 返回是不是主动关闭引发的Net断开
	ManualClose() bool

	// Proto 当NetEventType为NetEventProto时有效, 返回事件携带的协议
	Proto() *ffProto.Proto

	// Back 上层事件处理完毕后, 必须执行Back方法, 以回收所有相关资源
	Back()

	String() string
}
