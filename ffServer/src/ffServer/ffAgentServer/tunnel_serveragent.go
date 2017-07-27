package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/uuid"
	"ffProto"

	"sync"
)

// 通道-服务端侧
type tunnelServerAgent struct {
	muSessionOn sync.RWMutex
	sess        base.Session

	serverID int
}

func (tsa *tunnelServerAgent) uuid() uuid.UUID {
	if tsa.sess != nil {
		return tsa.sess.UUID()
	}
	return uuid.InvalidUUID
}

// 返回值仅表明协议是否添加到底层连接的发送缓冲队列里了。当没有连接时，将返回false
func (tsa *tunnelServerAgent) sendProto(p *ffProto.Proto, extraDataUUID uint64) bool {
	tsa.muSessionOn.RLock()
	defer tsa.muSessionOn.RUnlock()

	p.SetExtraData(ffProto.ExtraDataTypeUUID, extraDataUUID)
	if tsa.sess != nil {
		log.RunLogger.Printf("tunnelServerAgent.sendProto: uuid[%x] proto[%v]", extraDataUUID, p)
		tsa.sess.SendProto(p)
		return true
	}

	ffProto.BackProtoAfterSend(p)
	return false
}

func (tsa *tunnelServerAgent) onConnect() {
	// 管理
	serverAgentMgr.onConnect(tsa)

	// 通知上层逻辑
}

func (tsa *tunnelServerAgent) onDisConnect() {
	tsa.muSessionOn.Lock()
	defer tsa.muSessionOn.Unlock()

	// 管理
	serverAgentMgr.onDisconnect(tsa.uuid())

	// 清空
	tsa.sess = nil
}

func (tsa *tunnelServerAgent) OnEvent(protoID ffProto.MessageType, data interface{}) {
	log.RunLogger.Printf("tunnelServerAgent.OnEvent: uuid[%x] protoID[%s] data[%v]\n",
		tsa.uuid(), ffProto.MessageType_name[int32(protoID)], data)

	if protoID == ffProto.MessageType_SessionConnect {

		// 连接建立
		tsa.onConnect()

	} else if protoID == ffProto.MessageType_SessionDisConnect {

		// 连接断开
		tsa.onDisConnect()

		wg, _ := data.(*sync.WaitGroup)
		wg.Done()

	} else if protoID == ffProto.MessageType_ServerKeepAlive {

		// 维持连接不断开协议(什么都不需要做)

	} else {

		proto, _ := data.(*ffProto.Proto)
		if protoID == ffProto.MessageType_Kick {

			// 踢出用户协议
			proto.Unmarshal() // 不验证服务端发来的协议
			message := proto.Message().(*ffProto.MsgKick)
			clientAgentMgr.kick(uuid.UUID(proto.ExtraData()), message.Result)

		} else if protoID == ffProto.MessageType_ServerRegister {

			// 服务器注册协议
			proto.Unmarshal() // 不验证服务端发来的协议
			message := proto.Message().(*ffProto.MsgServerRegister)
			serverID := int(message.GetServerID())
			if serverAgentMgr.serverType == message.GetServerType() && serverAgentMgr.getServerAgent(serverID) == nil {
				//有效的服务器注册
				tsa.serverID = serverID

				log.RunLogger.Printf("tunnelServerAgent.onConnect: uuid[%x] serverID[%d] serverType[%s]",
					tsa.uuid(), serverID, serverAgentMgr.serverType)
			} else {
				log.FatalLogger.Printf("tunnelServerAgent.OnEvent: uuid[%x] invalid MsgServerRegister[%s:%d] not match serverType[%s] or duplicate serverID\n", tsa.uuid(), message.GetServerType(), serverID, serverAgentMgr.serverType)

				// 1000毫秒后断开连接
				tsa.sess.Close(100)
			}

		} else {

			// 普通协议转发
			if !clientAgentMgr.sendProto(uuid.UUID(proto.ExtraData()), proto) {
				p := ffProto.ApplyProtoForSend(ffProto.MessageType_AgentDisConnect)
				tsa.sendProto(p, proto.ExtraData())
			}

		}

	}
}
