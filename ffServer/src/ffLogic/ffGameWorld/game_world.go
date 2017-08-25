package ffGameWorld

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
)

type gameWorld struct {
}

// init 初始化
func (gw *gameWorld) init() error {
	// poolOfAccount 150%在线用户数
	err := poolOfAccount.init(worldFrame.DefaultOnlineCount() * 150 / 100)
	if err != nil {
		return nil
	}

	// managerOfAccount
	managerOfAccount.init()

	// 读取配置文件
	return nil
}

// Start 启动
func (gw *gameWorld) Start() {

}

// Stop 停止
func (gw *gameWorld) Stop() {

}

// DispatchProto 处理协议
//  uuidAgent: 连接唯一标识
//  p: 待处理协议
// 禁止缓存协议
func (gw *gameWorld) DispatchProto(uuidAgent uuid.UUID, p *ffProto.Proto) {
	account, okOnlineAccount := managerOfAccount.mapOnlineAccounts[uuidAgent]

	// 解析协议
	if err := p.Unmarshal(); err != nil {
		log.RunLogger.Printf("gameWorld.DispatchProto: uuidAgent[%x] proto[%v] Unmarshal error[%v]\n", uuidAgent, p, err)

		if okOnlineAccount {
			account.kick(ffError.ErrKickProtoInvalid)
		} else {
			worldFrame.Kick(uuidAgent, true, ffError.ErrKickProtoInvalid)
		}
		return
	}

	log.RunLogger.Printf("gameWorld.DispatchProto: uuidAgent[%x] proto[%v]\n", uuidAgent, p)

	if p.ProtoID() == ffProto.MessageType_EnterGameWorld {

		// 开始
		if okOnlineAccount {
			// 重复发送, 忽略
			log.RunLogger.Printf("gameWorld.DispatchProto: uuidAgent[%x]-uuidAccount[%x] ignore duplicate online\n", uuidAgent, account.uuidAccount)
			return
		}

		message := p.Message().(*ffProto.MsgEnterGameWorld)
		uuidAcount := uuid.NewUUID(message.UUIDLogin)

		// 帐号异地登录
		if uuidAgentOld, ok := managerOfAccount.mapAccountAgent[uuidAcount]; ok {
			if account, ok := managerOfAccount.mapOnlineAccounts[uuidAgentOld]; ok {
				log.RunLogger.Printf("gameWorld.DispatchProto: uuidAgent[%x]-uuidAccount[%x] exotic with uuidAgentOld[%x]\n", uuidAgent, uuidAcount, uuidAgentOld)

				account.kick(ffError.ErrKickExotic)
			} else {
				log.FatalLogger.Printf("gameWorld.DispatchProto: uuidAgent[%x]-uuidAccount[%x] exotic with uuidAgentOld[%x] but no online account\n", uuidAgent, uuidAcount, uuidAgentOld)

				delete(managerOfAccount.mapAccountAgent, uuidAcount)
			}
		}

		// 申请并初始化
		account = managerOfAccount.applyAccount(uuidAgent, uuidAcount)
		account.Init()

		// 发送进入游戏的反馈
		message.Result = ffError.ErrNone.Code()
		worldFrame.SendProto(uuidAgent, p)

		return

	} else if p.ProtoID() == ffProto.MessageType_AgentDisConnect {

		// 结束
		if okOnlineAccount {
			account.kicked(ffError.ErrKickConnection, true)
		} else {
			log.RunLogger.Printf("gameWorld.DispatchProto: uuidAgent[%x] not in online list\n", uuidAgent)
		}
		return

	} else {

		// 普通协议
		if !okOnlineAccount {
			log.RunLogger.Printf("gameWorld.DispatchProto: uuidAgent[%x] not in online list\n", uuidAgent)

			worldFrame.Kick(uuidAgent, true, ffError.ErrKickConnection)
			return
		}

		gw.onProto(uuidAgent, p)
		return

	}
}

// onProto 处理逻辑协议
func (gw *gameWorld) onProto(uuidAgent uuid.UUID, p *ffProto.Proto) {

}

// KickAll 所有人全部立即下线
//	kickReason: 踢出原因
func (gw *gameWorld) KickAll(kickReason ffError.Error) {
	managerOfAccount.onEventKickAll(kickReason)
}

// sendMsgKick 发送踢人协议
func (gw *gameWorld) sendMsgKick(uuidAgent uuid.UUID, kickReason ffError.Error) {
	p := ffProto.ApplyProtoForSend(ffProto.MessageType_Kick)
	message := p.Message().(*ffProto.MsgKick)
	message.Result = kickReason.Code()
	worldFrame.SendProto(uuidAgent, p)
}
