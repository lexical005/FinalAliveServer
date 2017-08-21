package main

// import (
// 	"ffAutoGen/ffError"
// 	"ffProto"
// 	"math/rand"
// )

// // 协议回调函数
// //	返回值表明接收到的Proto是否进入了发送逻辑(如果未正确设置返回值, 将导致泄露或者异常)
// var mapProtoCallback = map[ffProto.MessageType]func(agent *agentUser, proto *ffProto.Proto) bool{
// 	ffProto.MessageType_EnterGameWorld:               onProtoEnterGameWorld,
// 	ffProto.MessageType_PrepareLoginPlatformUniqueId: onProtoPrepareLoginPlatformUniqueID,
// 	ffProto.MessageType_LoginPlatformUniqueId:        onProtoLoginPlatformUniqueID,
// }

// func onProtoEnterGameWorld(agent *agentUser, proto *ffProto.Proto) bool {
// 	message, _ := proto.Message().(*ffProto.MsgEnterGameWorld)
// 	message.Result = ffError.ErrNone.Code()

// 	proto.ChangeLimitStateRecvToSend()
// 	agent.SendProto(proto)
// 	return true
// }

// func onProtoPrepareLoginPlatformUniqueID(agent *agentUser, proto *ffProto.Proto) bool {
// 	fixSalt := rand.Int31()
// 	for fixSalt == 0 {
// 		fixSalt = rand.Int31()
// 	}

// 	message, _ := proto.Message().(*ffProto.MsgPrepareLoginPlatformUniqueId)
// 	message.FixSalt = fixSalt

// 	proto.ChangeLimitStateRecvToSend()
// 	agent.SendProto(proto)
// 	return true
// }

// func onProtoLoginPlatformUniqueID(agent *agentUser, proto *ffProto.Proto) bool {
// 	message, _ := proto.Message().(*ffProto.MsgLoginPlatformUniqueId)
// 	message.UUIDLogin = agent.UUID().Value()

// 	proto.ChangeLimitStateRecvToSend()
// 	agent.SendProto(proto)
// 	return true
// }
