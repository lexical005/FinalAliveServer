package main

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/uuid"
	"ffProto"

	"sync"
)

// 通道-客户端侧
type tunnelClientAgent struct {
	muSessionOn sync.RWMutex // 确保对sess的操作顺序进行
	sess        base.Session

	serverAgent *tunnelServerAgent // 协议转发目标

	kicked bool // 被踢出
}

func (tca *tunnelClientAgent) uuid() uuid.UUID {
	if tca.sess != nil {
		return tca.sess.UUID()
	}
	return uuid.InvalidUUID
}

// 返回值仅表明协议是否添加到底层连接的发送缓冲队列里了。当没有连接时，将返回false
func (tca *tunnelClientAgent) sendProto(p *ffProto.Proto) bool {
	tca.muSessionOn.RLock()
	defer tca.muSessionOn.RUnlock()

	if tca.sess != nil && !tca.kicked {
		p.SetExtraData(ffProto.ExtraDataTypeNormal, 0)
		log.RunLogger.Printf("tunnelClientAgent.sendProto: uuid[%x] proto[%v]", tca.uuid(), p)
		tca.sess.SendProto(p)
		return true
	}

	ffProto.BackProtoAfterSend(p)
	return false
}

func (tca *tunnelClientAgent) onConnect() {
	tca.serverAgent = nil
	tca.kicked = false

	// 管理
	clientAgentMgr.onConnect(tca)

	// 通知上层逻辑
}

func (tca *tunnelClientAgent) onDisConnect() {
	tca.muSessionOn.Lock()
	defer tca.muSessionOn.Unlock()

	// 管理
	clientAgentMgr.onDisconnect(tca.uuid())

	// 通知服务端侧
	if tca.serverAgent != nil {
		sendProto := ffProto.ApplyProtoForSend(ffProto.MessageType_MT_MsgAgentDisConnect)
		tca.serverAgent.sendProto(sendProto, uint64(tca.uuid()))
	}

	// 清空
	tca.sess = nil
}

func (tca *tunnelClientAgent) OnEvent(protoID ffProto.MessageType, data interface{}) {
	log.RunLogger.Printf("tunnelClientAgent.OnEvent: uuid[%x] protoID[%s] data[%v]\n",
		tca.uuid(), ffProto.MessageType_name[int32(protoID)], data)

	if protoID == ffProto.MessageType_MT_Connect {

		// 连接建立
		tca.onConnect()

	} else if protoID == ffProto.MessageType_MT_DisConnect {

		// 连接断开
		tca.onDisConnect()

		wg, _ := data.(*sync.WaitGroup)
		wg.Done()

	} else if protoID == ffProto.MessageType_MT_MsgCSKeepAlive {

		// 维持连接不断开协议(什么都不需要做)

	} else {

		// 已被踢出
		if tca.kicked {
			return
		}

		// todo: 帐号登录协议, 转到帐号服务器处理

		// 进入游戏世界协议
		if protoID == ffProto.MessageType_MT_MsgEnterGameWorld {
			proto, _ := data.(*ffProto.Proto)
			if err := proto.Unmarshal(); err == nil {
				message := proto.Message().(*ffProto.MsgEnterGameWorld)
				tca.serverAgent = serverAgentMgr.getServerAgent(int(message.GetServerID()))

				// 目标服务器不存在
				if tca.serverAgent == nil {
					message.Result = ffError.ErrGameServerOffline.Code()
					tca.sendProto(proto)
					return
				}
			} else {
				tca.kick(true, ffError.ErrKickProtoInvalid)
				return
			}
		}

		if tca.serverAgent != nil {
			proto, _ := data.(*ffProto.Proto)
			if !tca.serverAgent.sendProto(proto, uint64(tca.uuid())) {
				// 对端的服务器连接已断开, 断开与此客户端的连接
				tca.kick(true, ffError.ErrKickConnection)
				return
			}
		} else {
			// 没有对端的服务器连接, 断开与此客户端的连接
			tca.kick(true, ffError.ErrKickConnection)
			return
		}

	}
}

// 踢出服务器, 断开连接
func (tca *tunnelClientAgent) kick(notifyKick bool, kickReason ffError.Error) {
	log.RunLogger.Printf("tunnelClientAgent.kick: uuid[%x] notifyKick[%v] kickReason[%v]", tca.uuid(), notifyKick, kickReason)

	// 踢出通知
	if !tca.kicked {
		tca.kicked = true

		if notifyKick {
			sendProto := ffProto.ApplyProtoForSend(ffProto.MessageType_MT_MsgKick)
			message := sendProto.Message().(*ffProto.MsgKick)
			message.Result = kickReason.Code()
			tca.sendProto(sendProto)
		}
	}

	// 2秒后关闭连接
	tca.sess.Close(2000)
}

func newClientAgent() interface{} {
	return &tunnelClientAgent{}
}
