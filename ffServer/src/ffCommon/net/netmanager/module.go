package netmanager

import (
	"ffCommon/net/base"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"time"
)

// INetSession 连接对象对外提供的方法
type INetSession interface {
	// UUID 唯一标识
	UUID() uuid.UUID

	// SendProtoExtraDataNormal 发送协议
	SendProtoExtraDataNormal(proto *ffProto.Proto) bool

	// SendProtoExtraDataUUID 发送协议
	SendProtoExtraDataUUID(uuidSender uint64, proto *ffProto.Proto) bool

	// Close 主动关闭
	Close()
}

// INetSessionHandler 连接对象关联的外界逻辑处理对象
//	INetSessionHandler 内部, 是单线的
//	INetSessionHandler 之间, 没有任何耦合
type INetSessionHandler interface {
	// UUID 唯一标识
	UUID() uuid.UUID

	// Init 初始化
	Init(netsession INetSession)

	// OnConnect 连接建立完成事件
	OnConnect()

	// OnDisConnect 连接关闭事件
	OnDisConnect()

	// OnProto 接收到协议, 返回值表明接收到的Proto是否进入了发送逻辑, 必须正确设置, 否则将导致泄露或者异常
	OnProto(proto *ffProto.Proto) bool
}

// INetSessionHandlerManager INetSessionHandler 管理器, 不会被异步调用
type INetSessionHandlerManager interface {
	// Create 创建INetSessionHandler
	Create(netSession INetSession) INetSessionHandler

	// Back 回收INetSessionHandler
	Back(handler INetSessionHandler)

	// End 退出完成
	End()
}

// NewServer 根据配置返回一个Server管理器
func NewServer(
	handlerManager INetSessionHandlerManager,
	config *base.ServeConfig,
	countApplicationQuit *int32,
	chApplicationQuit chan struct{}) (mgr *Manager, err error) {

	net, err := newNetServer(config)
	if err != nil {
		return
	}

	name := fmt.Sprintf("ServerManager[%v]", net.UUID())

	mgr = &Manager{
		name: name,

		sendKeepAliveInterval: 0,

		handlerManager: handlerManager,

		net: net,

		chAgentClosed: make(chan *agentSession, 2),
		mapAgent:      make(map[uuid.UUID]*agentSession, config.InitOnlineCount),

		countApplicationQuit: countApplicationQuit,
		chApplicationQuit:    chApplicationQuit,
	}

	go util.SafeGo(mgr.mainLoop, mgr.mainLoopEnd)

	return
}

// NewClient 根据配置返回一个Client管理器
func NewClient(
	handlerManager INetSessionHandlerManager,
	configConnect *base.ConnectConfig,
	configSession *base.SessionConfig,
	countApplicationQuit *int32,
	chApplicationQuit chan struct{}) (mgr *Manager, err error) {

	net, err := newNetClient(configConnect)
	if err != nil {
		return
	}

	name := fmt.Sprintf("ClientManager[%v]", net.UUID())

	mgr = &Manager{
		name: name,

		sendKeepAliveInterval: time.Duration(configConnect.KeepAliveInterval) * time.Second,

		handlerManager: handlerManager,

		net: net,

		chAgentClosed: make(chan *agentSession, 2),
		mapAgent:      make(map[uuid.UUID]*agentSession, 1),

		countApplicationQuit: countApplicationQuit,
		chApplicationQuit:    chApplicationQuit,
	}

	go util.SafeGo(mgr.mainLoop, mgr.mainLoopEnd)

	return
}
