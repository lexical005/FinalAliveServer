package base

import (
	"ffCommon/uuid"
	"ffProto"

	"net"
)

// ServerEvent ServerEvent
type ServerEvent interface {
	OnServerEnd()

	OnSessionStart()
	OnSessionData()
	OnSessionEnd()
}

// Session interface
type Session interface {
	// 启动 Session 主循环
	Start(conn net.Conn, agent Agent, recvProtoExtraDataType ffProto.ExtraDataType)

	// UUID 唯一标识
	UUID() uuid.UUID

	// 发送 Proto 到对端, 异步，当Agent处理完毕MessageType_MT_DisConnect后，禁止再调用此方法
	SendProto(p *ffProto.Proto)

	// 关闭，异步
	// delayMillisecond: 延迟多少毫秒关闭
	Close(delayMillisecond int64)

	// 外界调用此接口后，将返回一个管道，供外界等候Session关闭
	WaitCloseChan() <-chan bool
}

// Agent 底层连接向上层通知事件和协议的接口列表
// Agent 的方法, 是由运行在独立 goroutine 里的 Session 调用的
// 故:
// 单个 Agent 操作自身内部数据时, 可以当作单线程处理
// 多个 Agent 操作同一份共享资源时, 需要同步锁
type Agent interface {
	// 接收到对端发来的协议。实现时，如果缓存ffProto.Proto, 则必须调用Proto.SetCached, 否则, 其生命周期只在函数执行期间内有效
	// protoID: MessageType_MT_DisConnect 连接断开，Session以上层处理完毕该事件来作为Session结束的依据, data为sync.WaitGroup，处理完毕后，必须执行data.Done，
	//										以便Session逻辑继续执行。Agent执行完毕此事件后，即代表连接完全断开，禁止再调用Session的SendProto方法
	//			MessageType_MT_Connect 连接建立，上层以该事件来作为Session有效的起点，data无效
	//			>=0 协议编号，data为*ffProto.Proto
	OnEvent(protoID ffProto.MessageType, data interface{})
}

// AgentCreator AgentCreator
type AgentCreator interface {
	// 创建 Agent 实例
	Create(s Session) Agent
}

// Client connect Server
type Client interface {
	// 关闭
	Close()

	// 启动 Client 连接 Server
	Start(agent Agent, recvProtoExtraDataType ffProto.ExtraDataType) error

	// 发送 Proto 到对端, 由使用者确保当前连接有效
	SendProto(p *ffProto.Proto)

	String() string
}

// Server listen Client
type Server interface {
	// 关闭
	Close()

	// 启动 Client 连接 Server
	Start(agentCreator AgentCreator, recvProtoExtraDataType ffProto.ExtraDataType) error

	String() string
}
