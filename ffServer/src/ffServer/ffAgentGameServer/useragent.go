package main

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
)

// 用户连接到本服务器后的agent
type userAgent struct {
	uuidSession uuid.UUID

	sendExtraDataType ffProto.ExtraDataType  // 发送的Proto的附加数据类型
	chSendProto       chan *ffProto.Proto    // 待发送协议管道
	chNetEventData    chan base.NetEventData // session网络事件管道

	chClose       chan struct{}   // 用于接收外界通知关闭
	chAgentClosed chan *userAgent // 用于向外界汇报关闭了, 仅有使用权

	onceClose util.Once // 一次性关闭
}

func (agent *userAgent) String() string {
	return fmt.Sprintf("uuidUserAgent[%v]", agent.uuidSession)
}

func (agent *userAgent) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("userAgent.mainLoop: %v", agent)

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

				log.RunLogger.Printf("userAgent.mainLoop start close: %v", agent)

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
func (agent *userAgent) mainLoopEnd() {
	log.RunLogger.Printf("userAgent.mainLoopEnd: %v", agent)
}

func (agent *userAgent) onNetEventData(data base.NetEventData) bool {
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
func (agent *userAgent) onConnect(data base.NetEventData) {
	log.RunLogger.Printf("userAgent.onConnect data[%v]: %v", data, agent)
}

// onDisConnect 连接断开
func (agent *userAgent) onDisConnect(data base.NetEventData) {
	log.RunLogger.Printf("userAgent.onDisConnect data[%v]: %v", data, agent)

	agent.chAgentClosed <- agent
}

// onProto 收到Proto
func (agent *userAgent) onProto(data base.NetEventData) {
	log.RunLogger.Printf("userAgent.onProto data[%v]: %v", data, agent)

	proto := data.Proto()
	protoID := proto.ProtoID()

	// todo: 区分协议号, 有些协议直接转发的
	// 反序列化
	if err := proto.Unmarshal(); err != nil {
		log.RunLogger.Printf("userAgent.onProto proto[%v] Unmarshal error[%v]: %v", proto, err, agent)
		agent.Close()
		return
	}

	log.RunLogger.Printf("userAgent.onProto proto[%v]: %v", proto, agent)

	switch protoID {
	case ffProto.MessageType_EnterGameWorld:
		agent.onProtoEnterGameWorld(proto)
	}
}

func (agent *userAgent) onProtoEnterGameWorld(proto *ffProto.Proto) {
	message, _ := proto.Message().(*ffProto.MsgEnterGameWorld)
	message.Result = ffError.ErrNone.Code()

	agent.SendProto(proto)
}

// Start 初始化, 然后开始收发协议并处理
func (agent *userAgent) Start(sess base.Session, agentServer *userAgentServer) {
	agent.uuidSession = sess.UUID()
	agent.sendExtraDataType, agent.chAgentClosed = agentServer.sendExtraDataType, agentServer.chAgentClosed

	agent.chSendProto = make(chan *ffProto.Proto, agentServer.config.SessionSendProtoCache)
	agent.chNetEventData = make(chan base.NetEventData, agentServer.config.SessionNetEventDataCache)

	agent.chClose = make(chan struct{}, 1)

	agent.onceClose.Reset()

	sess.Start(agent.chSendProto, agent.chNetEventData, agentServer.recvExtraDataType)

	go util.SafeGo(agent.mainLoop, agent.mainLoopEnd)
}

// Close
func (agent *userAgent) Close() {
	agent.onceClose.Do(func() {
		agent.chClose <- struct{}{}
	})
}

// SendProto 发送Proto
func (agent *userAgent) SendProto(proto *ffProto.Proto) {
	if agent.sendExtraDataType == ffProto.ExtraDataTypeNormal {
		proto.SetExtraDataNormal()
	} else if agent.sendExtraDataType == ffProto.ExtraDataTypeUUID {
		proto.SetExtraDataUUID(agent.uuidSession.Value())
	}

	agent.chSendProto <- proto
}

// Back 可安全回收了
func (agent *userAgent) Back() {
	agent.chAgentClosed = nil

	// 清理待发送Proto
	close(agent.chSendProto)
	for proto := range agent.chSendProto {
		if proto != nil {
			proto.BackAfterSend()
		}
	}
	agent.chSendProto = nil

	// 关闭
	close(agent.chNetEventData)
	agent.chNetEventData = nil

	// 关闭
	close(agent.chClose)
	agent.chClose = nil
}

func newUserAgent() *userAgent {
	return &userAgent{}
}
