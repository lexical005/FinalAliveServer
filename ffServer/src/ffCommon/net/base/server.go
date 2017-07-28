package base

import (
	"ffProto"
)

// Server Server自身, 未实现多goroutine安全, 由使用者确保
type Server interface {
	// 关闭
	Close()

	// 启动 Server 监听 Client 链接
	Start(agentCreator AgentCreator, recvProtoExtraDataType ffProto.ExtraDataType) error

	String() string
}

// Agent 处理连接事件以及连接接收到的协议处理, 由底层base.Session触发
type Agent interface {
	// OnConnect 连接已建立, 自此调用开始, 可发送协议, 此方法在base.Session生命周期内只会被调用一次
	OnConnect()

	// OnDisConnect 连接已断开, 此调用返回后, base.Session将被回收, 不再引用base.Session对象, 此方法在base.Session生命周期内只会被调用一次
	// OnDisConnect之后, 不会再触发OnProto
	// manual: 是不是上层逻辑主动断开的
	OnDisConnect(manual bool)

	// OnProto 处理协议, 在上层在收到OnDisConnect调用后, 将不再会收到协议, 也就是说, 在Close调用之后, 依然有可能收到协议
	//	如果不修改proto状态, 则在此调用返回后, proto将被回收
	//	可修改proto状态为缓存等待分发(SetCacheWaitDispatch)或缓存等待发送(SetCacheWaitSend), 以避免回收, 具体proto使用, 参见ffProto模块readme.txt
	OnProto(protoID ffProto.MessageType, proto *ffProto.Proto)
}

// AgentCreator AgentCreator
type AgentCreator interface {
	// 创建 Agent 实例
	Create(s Session) Agent
}
