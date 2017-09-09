package main

import (
	"ffAutoGen/ffError"
	"ffProto"
	"math/rand"
)

// 协议回调函数
//	返回值表明接收到的Proto是否进入了发送逻辑(如果未正确设置返回值, 将导致泄露或者异常)
var mapAgentUserProtoCallback = map[ffProto.MessageType]func(agent *agentUser, proto *ffProto.Proto) bool{
	ffProto.MessageType_EnterGameWorld:               onProtoEnterGameWorld,
	ffProto.MessageType_PrepareLoginPlatformUniqueId: onProtoPrepareLoginPlatformUniqueID,
	ffProto.MessageType_LoginPlatformUniqueId:        onProtoLoginPlatformUniqueID,

	ffProto.MessageType_StartMatch: onProtoStartMatch,
	ffProto.MessageType_StopMatch:  onProtoStopMatch,

	ffProto.MessageType_PeaceCheat: onPeaceProtoCheat,
}

func onProtoEnterGameWorld(agent *agentUser, proto *ffProto.Proto) (result bool) {
	message, _ := proto.Message().(*ffProto.MsgEnterGameWorld)
	message.Result = ffError.ErrNone.Code()

	result = ffProto.SendProtoExtraDataNormal(agent, proto, true)

	// 进入MatchServer
	p := ffProto.ApplyProtoForSend(ffProto.MessageType_EnterMatchServer)
	m := p.Message().(*ffProto.MsgEnterMatchServer)
	m.UUIDAccount = agent.uuidAccount
	m.UUIDTeam = 0
	ffProto.SendProtoExtraDataUUID(instMatchServerClient, agent.UUID(), p, false)

	return
}

func onProtoPrepareLoginPlatformUniqueID(agent *agentUser, proto *ffProto.Proto) (result bool) {
	fixSalt := rand.Int31()
	for fixSalt == 0 {
		fixSalt = rand.Int31()
	}

	message, _ := proto.Message().(*ffProto.MsgPrepareLoginPlatformUniqueId)
	message.FixSalt = fixSalt

	agent.uuidPlatformLogin = message.UUIDPlatformLogin

	return ffProto.SendProtoExtraDataNormal(agent, proto, true)
}

func onProtoLoginPlatformUniqueID(agent *agentUser, proto *ffProto.Proto) (result bool) {
	loginData := &httpClientCustomLoginData{
		// 请求者
		uuidRequester: agent.UUID(),

		// 请求数据
		UUIDPlatform: agent.uuidPlatformLogin,
	}

	// 无法请求登录验证
	if !instHTTPLoginClient.PostCustomLogin(loginData) {
		agent.Close()
	}

	return
}

func onProtoStartMatch(agent *agentUser, proto *ffProto.Proto) (result bool) {
	if agent.battleUser != nil {
		return false
	}

	return ffProto.SendProtoExtraDataUUID(instMatchServerClient, agent.UUID(), proto, true)
}

func onProtoStopMatch(agent *agentUser, proto *ffProto.Proto) (result bool) {
	if agent.battleUser != nil {
		return false
	}

	return ffProto.SendProtoExtraDataUUID(instMatchServerClient, agent.UUID(), proto, true)
}

// 和平作弊指令
func onPeaceProtoCheat(agent *agentUser, proto *ffProto.Proto) (result bool) {
	return
}

func onCustomLoginResult(agent *agentUser, result *httpClientCustomLoginData) {
	proto := ffProto.ApplyProtoForSend(ffProto.MessageType_LoginPlatformUniqueId)
	message := proto.Message().(*ffProto.MsgLoginPlatformUniqueId)

	if result.err != nil {
		message.Result = ffError.ErrUnknown.Code()
	} else {
		agent.uuidAccount = result.UUIDAccount
		message.UUIDLogin = agent.UUID().Value()
	}

	ffProto.SendProtoExtraDataNormal(agent, proto, false)
}
