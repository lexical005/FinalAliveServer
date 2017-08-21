package netmanager

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

// 连接对象
type agentSession struct {
	name        string
	uuidSession uuid.UUID

	status int32 // 运行状态  0初始状态 1可使用 2使用中 -1关闭中(阻碍进入使用状态) -2完成了所有关闭工作(阻碍进入使用状态) -3关闭完成

	sendExtraDataType ffProto.ExtraDataType  // 发送的Proto的附加数据类型
	chSendProto       chan *ffProto.Proto    // 待发送协议管道
	chNetEventData    chan base.NetEventData // session网络事件管道

	chClose       chan struct{}      // 用于接收外界通知关闭
	chAgentClosed chan *agentSession // 用于向外界汇报关闭了, 仅有使用权

	onceClose util.Once // 一次连接期间, 关闭一次
}

func (agent *agentSession) String() string {
	return agent.name
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

				// 2使用中 ==> -1关闭中(阻碍进入使用状态)
				if !atomic.CompareAndSwapInt32(&agent.status, 2, -1) {
					// -2完成了所有关闭工作(阻碍进入使用状态)
					atomic.StoreInt32(&agent.status, -2)
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
func (agent *agentSession) mainLoopEnd() {
	log.RunLogger.Printf("%v.mainLoopEnd", agent.name)
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

	// 0初始状态 ==> 1可使用
	atomic.CompareAndSwapInt32(&agent.status, 0, 1)
}

// onDisConnect 连接断开, 此事件处理完毕后, session将不可用
func (agent *agentSession) onDisConnect(data base.NetEventData) {
	log.RunLogger.Printf("%v.onDisConnect data[%v]", agent.name, data)

	// 2使用中 ==> -1关闭中(阻碍进入使用状态)
	if !atomic.CompareAndSwapInt32(&agent.status, 2, -1) {
		// -2完成了所有关闭工作(阻碍进入使用状态)
		atomic.StoreInt32(&agent.status, -2)
	}

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

	// todo: 区分协议号, 有些协议直接转发的
	// 反序列化
	if err := proto.Unmarshal(); err != nil {
		log.RunLogger.Printf("%v.onProto proto[%v] Unmarshal error[%v]", agent.name, proto, err)
		agent.Close()
		return
	}

	log.RunLogger.Printf("%v.onProto proto[%v]", agent.name, proto)

	if callback, ok := mapProtoCallback[protoID]; ok {
		changedToSendState = callback(agent, proto)
	} else {
		log.FatalLogger.Printf("%v.onProto unknown protoID[%v]", agent.name, protoID)
	}
}

// Start 初始化, 然后开始收发协议并处理
func (agent *agentSession) Start(sess base.Session, net inet, chAgentClosed chan *agentSession) {
	agent.name = fmt.Sprintf("agentSession[%v]", sess.UUID())
	agent.uuidSession = sess.UUID()

	agent.sendExtraDataType, agent.chAgentClosed = net.SendExtraDataType(), chAgentClosed

	agent.chSendProto = make(chan *ffProto.Proto, net.SessionSendProtoCache())
	agent.chNetEventData = make(chan base.NetEventData, net.SessionNetEventDataCache())

	agent.chClose = make(chan struct{}, 1)

	agent.onceClose.Reset()

	agent.status = 0

	sess.Start(agent.chSendProto, agent.chNetEventData, net.RecvExtraDataType())

	go util.SafeGo(agent.mainLoop, agent.mainLoopEnd)
}

// Close
func (agent *agentSession) Close() {
	agent.onceClose.Do(func() {
		agent.chClose <- struct{}{}
	})
}

// SendProto 发送Proto
//	返回值仅表明请求发送的协议, 是否被添加到待发送管道内, 不代表一定能发送到对端
func (agent *agentSession) SendProto(proto *ffProto.Proto) bool {
	// 1可使用 ==> 2使用中
	if !atomic.CompareAndSwapInt32(&agent.status, 1, 2) {
		return false
	}

	if agent.sendExtraDataType == ffProto.ExtraDataTypeNormal {
		proto.SetExtraDataNormal()
	} else if agent.sendExtraDataType == ffProto.ExtraDataTypeUUID {
		proto.SetExtraDataUUID(agent.uuidSession.Value())
	}

	agent.chSendProto <- proto

	// 2使用中 ==> 1可使用
	if !atomic.CompareAndSwapInt32(&agent.status, 2, 1) {
		// -2完成了所有关闭工作(阻碍进入使用状态)
		atomic.StoreInt32(&agent.status, -2)
	}

	return true
}

// Back 可安全回收了
func (agent *agentSession) Back() {
	agent.chAgentClosed = nil

	// 清理待发送Proto
	{
		// 等待发送方法执行完毕
		waitCount, maxWaitCount := 0, 10
		for {
			// -2完成了所有关闭工作(阻碍进入使用状态) ==> -3关闭完成
			if atomic.CompareAndSwapInt32(&agent.status, -2, -3) {
				break
			}

			// 等待1秒
			<-time.After(time.Second)

			waitCount++
			if waitCount > maxWaitCount {
				log.FatalLogger.Printf("Back wait SendProto too long time[%v] to exit", waitCount)
				break
			}
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

func newAgentSession() *agentSession {
	return &agentSession{}
}
