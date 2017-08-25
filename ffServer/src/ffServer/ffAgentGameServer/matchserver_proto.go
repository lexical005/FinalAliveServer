package main

import (
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
)

// 协议回调函数
//	返回值表明接收到的Proto是否进入了发送逻辑(如果未正确设置返回值, 将导致泄露或者异常)
var mapMatchServerProtoCallback = map[ffProto.MessageType]func(server *matchServer, proto *ffProto.Proto) bool{
	ffProto.MessageType_EnterMatchServer: onProtoEnterMatchServer,
	ffProto.MessageType_LeaveMatchServer: onProtoLeaveMatchServer,

	ffProto.MessageType_StartMatch:  onProtoStartMatchResult,
	ffProto.MessageType_StopMatch:   onProtoStopMatchResult,
	ffProto.MessageType_MatchResult: onProtoMatchResult,
}

func onProtoEnterMatchServer(server *matchServer, proto *ffProto.Proto) (result bool) {
	log.RunLogger.Printf("onProtoEnterMatchServer agent[%v] proto[%v]", uuid.NewUUID(proto.ExtraData()), proto)
	return false
}

func onProtoLeaveMatchServer(server *matchServer, proto *ffProto.Proto) (result bool) {
	log.RunLogger.Printf("onProtoLeaveMatchServer agent[%v] proto[%v]", uuid.NewUUID(proto.ExtraData()), proto)
	return false
}

func onProtoStartMatchResult(server *matchServer, proto *ffProto.Proto) (result bool) {
	return instAgentUserServer.OnMatchServerProto(proto)
}

func onProtoStopMatchResult(server *matchServer, proto *ffProto.Proto) (result bool) {
	return instAgentUserServer.OnMatchServerProto(proto)
}

func onProtoMatchResult(server *matchServer, proto *ffProto.Proto) (result bool) {
	return instAgentUserServer.OnMatchServerProto(proto)
}
