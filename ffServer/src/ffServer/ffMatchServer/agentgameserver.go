package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"sync/atomic"
	"time"
)

// 协议回调函数
//	返回值表明接收到的Proto是否进入了发送逻辑(如果未正确设置返回值, 将导致泄露或者异常)
var mapProtoCallback = map[ffProto.MessageType]func(agent *agentGameServer, proto *ffProto.Proto) bool{
	ffProto.MessageType_KeepAlive: onProtoKeepAlive,
}

// AgentGameServer连接到本服务器后, 在本服务器的对象
type agentGameServer struct {
	uuid uuid.UUID

	sendExtraDataType ffProto.ExtraDataType  // 发送的Proto的附加数据类型
	sendStatus        int32                  // 运行状态  0不可使用状态 1可发送 2正在发送 -2发送完毕即进入不可使用状态 -1关闭完成
	chSendProto       chan *ffProto.Proto    // 待发送协议管道
	chNetEventData    chan base.NetEventData // session网络事件管道

	chClose       chan struct{}         // 用于接收外界通知关闭
	chAgentClosed chan *agentGameServer // 用于向外界汇报关闭了, 仅有使用权

	onceClose util.Once // 一次连接期间, 关闭一次
}

func (agent *agentGameServer) String() string {
	return fmt.Sprintf("uuid[%v]", agent.uuid)
}

func (agent *agentGameServer) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("agentGameServer.mainLoop: %v", agent)

	// 主循环
	{
		for {
			select {
			case data := <-agent.chNetEventData: // 处理事件

				if agent.onNetEventData(data) {
					// 在外界通知关闭之前, 连接已报错关闭
					return
				}

			case <-agent.chClose: // 外界通知关闭

				log.RunLogger.Printf("agentGameServer.mainLoop start close: %v", agent)

				// 发送中状态 ==> 等待发送结束关闭状态
				if !atomic.CompareAndSwapInt32(&agent.sendStatus, 2, -2) {
					// 不可使用状态
					atomic.StoreInt32(&agent.sendStatus, 0)
				}

				agent.chSendProto <- nil

				for {
					select {
					case data := <-agent.chNetEventData: // 继续处理, 直到处理到连接关闭
						if agent.onNetEventData(data) {
							return
						}
					}
				}
			}
		}
	}
}
func (agent *agentGameServer) mainLoopEnd() {
	log.RunLogger.Printf("agentGameServer.mainLoopEnd: %v", agent)
}

func (agent *agentGameServer) onNetEventData(data base.NetEventData) bool {
	defer data.Back()

	switch data.NetEventType() {
	case base.NetEventOn:
		agent.onConnect(data)
	case base.NetEventOff:
		agent.onDisConnect(data)
		return true
	case base.NetEventProto:
		agent.onProto(data)
	}
	return false
}

// onConnect 连接建立
func (agent *agentGameServer) onConnect(data base.NetEventData) {
	log.RunLogger.Printf("agentGameServer.onConnect data[%v]: %v", data, agent)

	// 可发送状态
	atomic.StoreInt32(&agent.sendStatus, 1)
}

// onDisConnect 连接断开
func (agent *agentGameServer) onDisConnect(data base.NetEventData) {
	log.RunLogger.Printf("agentGameServer.onDisConnect data[%v]: %v", data, agent)

	// 发送中状态 ==> 等待发送结束关闭状态
	if !atomic.CompareAndSwapInt32(&agent.sendStatus, 2, -2) {
		// 不可使用状态
		atomic.StoreInt32(&agent.sendStatus, 0)
	}

	agent.chAgentClosed <- agent
}

// onProto 收到Proto
func (agent *agentGameServer) onProto(data base.NetEventData) {
	log.RunLogger.Printf("agentGameServer.onProto data[%v]: %v", data, agent)

	changedToSendState := false

	proto := data.Proto()
	protoID := proto.ProtoID()

	// 如果协议在处理完毕后, 未进入发送逻辑, 则回收
	defer func() {
		if !changedToSendState {
			proto.BackAfterDispatch()
		}
	}()

	// todo: 区分协议号, 有些协议直接转发的
	// 反序列化
	if err := proto.Unmarshal(); err != nil {
		log.RunLogger.Printf("agentGameServer.onProto proto[%v] Unmarshal error[%v]: %v", proto, err, agent)
		agent.Close()
		return
	}

	log.RunLogger.Printf("agentGameServer.onProto proto[%v]: %v", proto, agent)

	if callback, ok := mapProtoCallback[protoID]; ok {
		changedToSendState = callback(agent, proto)
	} else {
		log.FatalLogger.Printf("agentGameServer.onProto unknown protoID[%v]: %v", protoID, agent)
	}
}

// Start 初始化, 然后开始收发协议并处理
func (agent *agentGameServer) Start(sess base.Session, agentManager *agentGameServerManager) {
	agent.uuid = sess.UUID()
	agent.sendExtraDataType, agent.chAgentClosed = agentManager.sendExtraDataType, agentManager.chAgentClosed

	agent.chSendProto = make(chan *ffProto.Proto, agentManager.config.SessionSendProtoCache)
	agent.chNetEventData = make(chan base.NetEventData, agentManager.config.SessionNetEventDataCache)

	agent.chClose = make(chan struct{}, 1)

	agent.onceClose.Reset()

	sess.Start(agent.chSendProto, agent.chNetEventData, agentManager.recvExtraDataType)

	go util.SafeGo(agent.mainLoop, agent.mainLoopEnd)
}

// Close
func (agent *agentGameServer) Close() {
	agent.onceClose.Do(func() {
		agent.chClose <- struct{}{}
	})
}

// SendProto 发送Proto
//	返回值仅表明请求发送的协议, 是否被添加到待发送管道内, 不代表一定能发送到对端
func (agent *agentGameServer) SendProto(proto *ffProto.Proto) bool {
	// 可发送状态 ==> 发送中状态
	if !atomic.CompareAndSwapInt32(&agent.sendStatus, 1, 2) {
		return false
	}

	if agent.sendExtraDataType == ffProto.ExtraDataTypeNormal {
		proto.SetExtraDataNormal()
	} else if agent.sendExtraDataType == ffProto.ExtraDataTypeUUID {
		proto.SetExtraDataUUID(agent.uuid.Value())
	}

	agent.chSendProto <- proto

	// 发送中状态 ==> 可发送状态
	if !atomic.CompareAndSwapInt32(&agent.sendStatus, 2, 1) {
		// 关闭状态
		atomic.StoreInt32(&agent.sendStatus, 0)
	}

	return true
}

// Back 可安全回收了
func (agent *agentGameServer) Back() {
	agent.chAgentClosed = nil

	// 清理待发送Proto
	{
		// 确保关闭完成
		for index := 0; index < 10; index++ {
			// 不可使用状态 ==> 关闭完成状态
			if atomic.CompareAndSwapInt32(&agent.sendStatus, 0, -1) {
				break
			}

			// 等待1秒
			<-time.After(time.Second)
		}

		// 关闭发送协议管道
		close(agent.chSendProto)
		for proto := range agent.chSendProto {
			if proto != nil {
				proto.BackAfterSend()
			} else {
				break
			}
		}
		agent.chSendProto = nil
	}

	// 关闭
	close(agent.chNetEventData)
	agent.chNetEventData = nil

	// 关闭
	close(agent.chClose)
	agent.chClose = nil
}

func onProtoKeepAlive(agent *agentGameServer, proto *ffProto.Proto) bool {
	proto.ChangeLimitStateRecvToSend()
	agent.SendProto(proto)
	return true
}

func newAgentGameServer() *agentGameServer {
	return &agentGameServer{}
}
