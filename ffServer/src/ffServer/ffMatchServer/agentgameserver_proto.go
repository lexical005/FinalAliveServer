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

func onProtoServerRegister(agent *agentGameServer, proto *ffProto.Proto) bool {
	message, _ := proto.Message().(*ffProto.MsgServerRegister)
	if message.ServerType != "AgentGameServer" {
		log.FatalLogger.Printf("agentGameServer recv not support server register[%v]", message)
		agent.Close()
		return false
	}

	return false
}
