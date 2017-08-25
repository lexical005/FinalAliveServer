package main

import (
	"ffCommon/log/log"
	"ffCommon/net/netmanager"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
)

// 连接对象
type agentGameServer struct {
	name string

	netsession netmanager.INetSession
}

func (agent *agentGameServer) String() string {
	return ""
}

// OnConnect 连接建立
func (agent *agentGameServer) OnConnect() {
	log.RunLogger.Printf("%v.OnConnect", agent.name)
}

// OnDisConnect 连接断开, 此事件处理完毕后, session将不可用
func (agent *agentGameServer) OnDisConnect() {
	log.RunLogger.Printf("%v.OnDisConnect", agent.name)
}

// OnProto 收到Proto
func (agent *agentGameServer) OnProto(proto *ffProto.Proto) bool {
	log.RunLogger.Printf("%v.OnProto proto[%v]", agent.name, proto)

	protoID := proto.ProtoID()

	if callback, ok := mapProtoCallback[protoID]; ok {
		return callback(agent, proto)
	}

	log.FatalLogger.Printf("%v.OnProto not support protoID[%v]", agent.name, protoID)

	return false
}

// SendProtoExtraDataUUID 发送Proto
//	返回值仅表明请求发送的协议, 是否被添加到待发送管道内, 不代表一定能发送到对端
func (agent *agentGameServer) SendProtoExtraDataUUID(uuidExtraData uuid.UUID, proto *ffProto.Proto) bool {
	return agent.netsession.SendProtoExtraDataUUID(uuidExtraData.Value(), proto)
}

// UUID
func (agent *agentGameServer) UUID() uuid.UUID {
	return agent.netsession.UUID()
}

// Init 初始化
func (agent *agentGameServer) Init(netsession netmanager.INetSession) {
	agent.name = fmt.Sprintf("agentGameServer[%v]", netsession.UUID())
	agent.netsession = netsession
}

// Back 回收
func (agent *agentGameServer) Back() {
	agent.netsession = nil
}

// Close 主动关闭
func (agent *agentGameServer) Close() {
	agent.netsession.Close()
}

func newAgentGameServer() *agentGameServer {
	return &agentGameServer{}
}
