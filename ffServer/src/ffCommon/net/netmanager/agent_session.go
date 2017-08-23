package netmanager

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
)

// 连接对象
type agentSession struct {
	name              string
	uuid              uuid.UUID
	handler           INetSessionHandler
	responseKeepAlive bool // 是否响应并返回KeepAlive协议

	sendExtraDataType ffProto.ExtraDataType  // 发送的Proto的附加数据类型
	chSendProto       chan *ffProto.Proto    // 待发送协议管道
	chNetEventData    chan base.NetEventData // session网络事件管道

	chClose       chan struct{}      // 用于接收外界通知关闭
	chAgentClosed chan *agentSession // 用于向外界汇报关闭了, 仅有使用权

	status util.Worker // 可使用性状态管理, 内含一次性关闭
}

func (agent *agentSession) String() string {
	return fmt.Sprintf("name[%v] status[%v]", agent.name, agent.status.String())
}

// UUID 唯一标识
func (agent *agentSession) UUID() uuid.UUID {
	return agent.uuid
}

func (agent *agentSession) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("%v.mainLoop", agent.name)

	// 主循环
	{
		for {
			select {
			case data := <-agent.chNetEventData: // 处理事件

				if agent.onNetEventData(data) {
					// 在外界通知关闭之前, 连接已关闭
					return
				}

			case <-agent.chClose: // 外界通知关闭

				log.RunLogger.Printf("%v.mainLoop start close", agent.name)

				agent.Close()

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
func (agent *agentSession) mainLoopEnd(isPanic bool) {
	log.RunLogger.Printf("%v.mainLoopEnd isPanic[%v]", agent.name, isPanic)
}

func (agent *agentSession) onNetEventData(data base.NetEventData) bool {
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
func (agent *agentSession) onConnect(data base.NetEventData) {
	log.RunLogger.Printf("%v.onConnect data[%v]", agent.name, data)

	agent.status.Ready()

	agent.handler.OnConnect()
}

// onDisConnect 连接断开, 此事件处理完毕后, session将不可用
func (agent *agentSession) onDisConnect(data base.NetEventData) {
	log.RunLogger.Printf("%v.onDisConnect data[%v]", agent.name, data)

	agent.status.Close()

	agent.handler.OnDisConnect()

	agent.chAgentClosed <- agent
}

// onProto 收到Proto
func (agent *agentSession) onProto(data base.NetEventData) {
	log.RunLogger.Printf("%v.onProto data[%v]", agent.name, data)

	changedToSendState := false

	proto := data.Proto()
	protoID := proto.ProtoID()

	// 如果协议在处理完毕后, 未进入发送逻辑, 则回收
	defer func() {
		if !changedToSendState {
			proto.BackAfterDispatch()
		}
	}()

	if protoID == ffProto.MessageType_KeepAlive {

		if agent.responseKeepAlive {
			// 维持活跃协议, 直接返回
			changedToSendState = true
			proto.ChangeLimitStateRecvToSend()

			if agent.sendExtraDataType == ffProto.ExtraDataTypeNormal {
				agent.SendProtoExtraDataNormal(proto)
			} else if agent.sendExtraDataType == ffProto.ExtraDataTypeUUID {
				agent.SendProtoExtraDataUUID(proto.ExtraData(), proto)
			}
		}

	} else {

		// 其他协议, 则反序列化
		if err := proto.Unmarshal(); err != nil {
			log.RunLogger.Printf("%v.onProto proto[%v] Unmarshal error[%v]", agent.name, proto, err)
			agent.Close()
			return
		}

		changedToSendState = agent.handler.OnProto(proto)
	}
}

// init 初始化
func (agent *agentSession) init(sess base.Session, net inet, chAgentClosed chan *agentSession, responseKeepAlive bool) {
	agent.name = fmt.Sprintf("agentSession[%v]", sess.UUID())
	agent.uuid = sess.UUID()
	agent.responseKeepAlive = responseKeepAlive

	agent.sendExtraDataType, agent.chAgentClosed = net.SendExtraDataType(), chAgentClosed

	agent.chSendProto = make(chan *ffProto.Proto, net.SessionSendProtoCache())
	agent.chNetEventData = make(chan base.NetEventData, net.SessionNetEventDataCache())

	log.RunLogger.Printf("%v.init sendExtraDataType[%v]", agent.name, agent.sendExtraDataType)

	agent.chClose = make(chan struct{}, 1)

	agent.status.Reset()
}

// Start 启动, 收发协议
func (agent *agentSession) Start(sess base.Session, net inet, handler INetSessionHandler) {
	agent.handler = handler

	sess.Start(agent.chSendProto, agent.chNetEventData, net.RecvExtraDataType())

	go util.SafeGo(agent.mainLoop, agent.mainLoopEnd)
}

// Close
func (agent *agentSession) Close() {
	agent.status.Close()
}

// SendProtoExtraDataNormal 发送Proto, 附加数据类型ExtraDataTypeNormal
//	返回值仅表明请求发送的协议, 是否被添加到待发送管道内, 不代表一定能发送到对端. 当协议未被添加到待发送管道内时, 将被执行回收
func (agent *agentSession) SendProtoExtraDataNormal(proto *ffProto.Proto) bool {
	if agent.sendExtraDataType != ffProto.ExtraDataTypeNormal {
		log.FatalLogger.Printf("%v.SendProtoExtraDataNormal not match agent sendExtraDataType[%v] vs [%v]",
			agent.name, agent.sendExtraDataType, ffProto.ExtraDataTypeNormal)

		proto.BackAfterSend()
		return false
	}

	work := agent.status.EnterWork()

	defer func() {
		agent.status.LeaveWork(work)

		// 直接回收
		if !work {
			proto.BackAfterSend()
		}
	}()

	if work {
		proto.SetExtraDataNormal()

		agent.chSendProto <- proto
	}

	return work
}

// SendProtoExtraDataUUID 发送Proto, 附加数据类型ExtraDataTypeUUID
//	返回值仅表明请求发送的协议, 是否被添加到待发送管道内, 不代表一定能发送到对端. 当协议未被添加到待发送管道内时, 将被执行回收
func (agent *agentSession) SendProtoExtraDataUUID(uuidSender uint64, proto *ffProto.Proto) bool {
	if agent.sendExtraDataType != ffProto.ExtraDataTypeUUID {
		log.FatalLogger.Printf("%v.SendProtoExtraDataUUID not match agent sendExtraDataType[%v] vs [%v]",
			agent.name, agent.sendExtraDataType, ffProto.ExtraDataTypeUUID)

		proto.BackAfterSend()

		return false
	}

	work := agent.status.EnterWork()

	defer func() {
		agent.status.LeaveWork(work)

		// 直接回收
		if !work {
			proto.BackAfterSend()
		}
	}()

	if work {
		proto.SetExtraDataUUID(uuidSender)

		agent.chSendProto <- proto
	}

	return work
}

// Back 可安全回收了
func (agent *agentSession) Back() {
	agent.chAgentClosed = nil
	agent.handler = nil

	// 等待使用完毕
	agent.status.WaitWorkEnd(10)

	// 清理待发送Proto
	{
		// 关闭发送协议管道
		close(agent.chSendProto)
		for proto := range agent.chSendProto {
			if proto != nil {
				proto.BackAfterSend()
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

// keepAlive 发送KeepAlive协议, 保持连接有效
func (agent *agentSession) keepAlive() {
	proto := ffProto.ApplyProtoForSend(ffProto.MessageType_KeepAlive)
	message := proto.Message().(*ffProto.MsgKeepAlive)
	message.Number = 0
	agent.SendProtoExtraDataUUID(agent.UUID().Value(), proto)
}

func newAgentSession() *agentSession {
	return &agentSession{}
}
