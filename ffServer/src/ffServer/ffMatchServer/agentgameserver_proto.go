package main

import (
	"ffCommon/log/log"
	"ffProto"
)

// 协议回调函数
//	返回值表明接收到的Proto是否进入了发送逻辑(如果未正确设置返回值, 将导致泄露或者异常)
var mapProtoCallback = map[ffProto.MessageType]func(agent *agentGameServer, proto *ffProto.Proto) bool{
	ffProto.MessageType_ServerRegister: onProtoServerRegister,
}

// 服务器注册
func onProtoServerRegister(server *agentGameServer, proto *ffProto.Proto) bool {
	message, _ := proto.Message().(*ffProto.MsgServerRegister)
	if message.ServerType != "AgentGameServer" {
		log.FatalLogger.Printf("agentGameServer recv not support server register[%v]", message)
		server.Close()
		return false
	}

	return false
}

// 开始匹配
func onProtoStartMatch(server *agentGameServer, proto *ffProto.Proto) bool {
	proto.ChangeLimitStateRecvToSend()
	instMatchMgr.OnPlayerMatchProto(proto)
	return true
}

// 停止匹配
func onProtoStopMatch(server *agentGameServer, proto *ffProto.Proto) bool {
	proto.ChangeLimitStateRecvToSend()
	instMatchMgr.OnPlayerMatchProto(proto)
	return true
}

// 用户加入匹配服务器
func onProtoPlayerEnterMatchServer(server *agentGameServer, proto *ffProto.Proto) bool {
	proto.ChangeLimitStateRecvToSend()
	uuidPlayerKey := proto.ExtraData()

	message, _ := proto.Message().(*ffProto.MsgEnterMatchServer)
	instMatchPalyerMgr.AddPlayer(server, uuidPlayerKey, message.UUIDAccount, message.UUIDTeam)

	server.SendProto(uuidPlayerKey, proto)

	return true
}

// 用户离开匹配服务器
func onProtoPlayerLeaveMatchServer(server *agentGameServer, proto *ffProto.Proto) bool {
	proto.ChangeLimitStateRecvToSend()
	uuidPlayerKey := proto.ExtraData()

	instMatchPalyerMgr.RemovePlayer(uuidPlayerKey)

	server.SendProto(uuidPlayerKey, proto)

	return true
}
