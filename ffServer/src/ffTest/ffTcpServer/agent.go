package main

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
	"sync"
)

type agentClient struct {
	uuidSession uuid.UUID // 底层连接

	mutexWorking sync.RWMutex // working锁
	working      bool         // session是否可收发数据

	chSendProto chan *ffProto.Proto
}

// SendProto 发送Proto
func (agent *agentClient) SendProto(proto *ffProto.Proto) {
	agent.mutexWorking.RLock()
	defer agent.mutexWorking.RUnlock()

	if agent.working {
		agent.chSendProto <- proto
	}
}

// CloseSession 关闭连接
func (agent *agentClient) CloseSession() {
	agent.mutexWorking.RLock()
	defer agent.mutexWorking.RUnlock()

	if agent.working {
		close(agent.chSendProto)
	}
}

// SendProtoChannel 发送Proto的管道
func (agent *agentClient) SendProtoChannel() chan *ffProto.Proto {
	return agent.chSendProto
}

// OnNetConnect 连接建立事件
func (agent *agentClient) OnNetConnect() {

}

// OnNetProto 接收到Proto事件
func (agent *agentClient) OnNetProto(proto *ffProto.Proto) {
	message, _ := proto.Message().(*ffProto.MsgEnterGameWorld)
	if message.Result != ffError.ErrNone.Code() {
		log.FatalLogger.Println(ffError.ErrByCode(message.Result))
		return
	}

	log.RunLogger.Printf("agentClient.onSessionProto message[%v]\n", message)

	protoID := proto.ProtoID()
	if protoID == ffProto.MessageType_EnterGameWorld {
		agent.onProtoEnterGameWorld(proto)
	}
}

// OnNetDisConnect 连接断开事件
func (agent *agentClient) OnNetDisConnect(manualClose bool) {
	agent.mutexWorking.Lock()
	defer agent.mutexWorking.Unlock()

	agent.working = false
}

func (agent *agentClient) onProtoEnterGameWorld(proto *ffProto.Proto) {
	message, _ := proto.Message().(*ffProto.MsgEnterGameWorld)
	message.Result = ffError.ErrNone.Code()
	agent.SendProto(proto)
}

func newAgentClient(uuidSession uuid.UUID) *agentClient {
	return &agentClient{
		uuidSession: uuidSession,
	}
}
