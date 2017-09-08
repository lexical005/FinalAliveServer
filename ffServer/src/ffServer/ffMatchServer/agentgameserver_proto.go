package main

import (
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
)

// 协议回调函数
//	返回值表明接收到的Proto是否进入了发送逻辑(如果未正确设置返回值, 将导致泄露或者异常)
var mapProtoCallback = map[ffProto.MessageType]func(agent *agentGameServer, proto *ffProto.Proto) bool{
	ffProto.MessageType_ServerRegister: onProtoServerRegister,

	ffProto.MessageType_EnterMatchServer: onProtoPlayerEnterMatchServer,
	ffProto.MessageType_LeaveMatchServer: onProtoPlayerLeaveMatchServer,

	ffProto.MessageType_StartMatch: onProtoStartMatch,
	ffProto.MessageType_StopMatch:  onProtoStopMatch,
}

// 服务器注册
func onProtoServerRegister(server *agentGameServer, proto *ffProto.Proto) bool {
	message, _ := proto.Message().(*ffProto.MsgServerRegister)
	if message.ServerType != "AgentGameServer" {
		log.FatalLogger.Printf("agentGameServer recv not support server register[%v]", message)
		server.Close()
	} else {
		server.serverID = message.ServerID
	}

	return false
}

// 开始匹配
func onProtoStartMatch(server *agentGameServer, proto *ffProto.Proto) bool {
	return instMatchMgr.OnPlayerMatchProto(proto)
}

// 停止匹配
func onProtoStopMatch(server *agentGameServer, proto *ffProto.Proto) bool {
	return instMatchMgr.OnPlayerMatchProto(proto)
}

// 用户加入匹配服务器
func onProtoPlayerEnterMatchServer(server *agentGameServer, proto *ffProto.Proto) bool {
	uuidPlayerKey := uuid.NewUUID(proto.ExtraData())
	message, _ := proto.Message().(*ffProto.MsgEnterMatchServer)
	instMatchPlayerMgr.AddPlayer(server, uuidPlayerKey, uuid.NewUUID(message.UUIDAccount), uuid.NewUUID(message.UUIDTeam))
	message.UUIDAccount = uuid.InvalidUUID.Value()
	message.UUIDTeam = uuid.InvalidUUID.Value()

	return ffProto.SendProtoExtraDataUUID(server, uuidPlayerKey, proto, true)
}

// 用户离开匹配服务器
func onProtoPlayerLeaveMatchServer(server *agentGameServer, proto *ffProto.Proto) bool {
	return instMatchMgr.OnPlayerMatchProto(proto)
}
