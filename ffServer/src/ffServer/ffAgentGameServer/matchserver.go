package main

import (
	"ffCommon/log/log"
	"ffCommon/net/netmanager"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
)

// 连接对象
type matchServer struct {
	name string

	netsession netmanager.INetSession
}

func (agent *matchServer) String() string {
	return ""
}

// OnConnect 连接建立
func (agent *matchServer) OnConnect() {
	log.RunLogger.Printf("%v.OnConnect", agent.name)

	proto := ffProto.ApplyProtoForSend(ffProto.MessageType_ServerRegister)
	message := proto.Message().(*ffProto.MsgServerRegister)
	message.ServerID = appConfig.Server.ServerID
	message.ServerType = appConfig.Server.ServerType

	agent.netsession.SendProto(proto)
}

// OnDisConnect 连接断开, 此事件处理完毕后, session将不可用
func (agent *matchServer) OnDisConnect() {
	log.RunLogger.Printf("%v.OnDisConnect", agent.name)
}

// OnProto 收到Proto
func (agent *matchServer) OnProto(proto *ffProto.Proto) bool {
	log.RunLogger.Printf("%v.OnProto proto[%v]", agent.name, proto)

	protoID := proto.ProtoID()

	// if callback, ok := mapProtoCallback[protoID]; ok {
	// 	return callback(agent, proto)
	// }

	log.FatalLogger.Printf("%v.OnProto not support protoID[%v]", agent.name, protoID)

	return false
}

// SendProto 发送Proto
//	返回值仅表明请求发送的协议, 是否被添加到待发送管道内, 不代表一定能发送到对端
func (agent *matchServer) SendProto(proto *ffProto.Proto) bool {
	return agent.netsession.SendProto(proto)
}

// UUID
func (agent *matchServer) UUID() uuid.UUID {
	return agent.netsession.UUID()
}

// Init 初始化
func (agent *matchServer) Init(netsession netmanager.INetSession) {
	agent.name = fmt.Sprintf("matchServer[%v]", netsession.UUID())
	agent.netsession = netsession
}

// Back 回收
func (agent *matchServer) Back() {
	agent.netsession = nil
}

// Close 主动关闭
func (agent *matchServer) Close() {
	agent.netsession.Close()
}

func newMatchServer() *matchServer {
	return &matchServer{}
}
