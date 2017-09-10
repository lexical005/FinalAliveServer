package main

import (
	"ffCommon/log/log"
	"ffCommon/net/netmanager"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
)

// 连接对象
type agentUser struct {
	name string

	netsession netmanager.INetSession

	uuidPlatformLogin string // 用户平台唯一标识
	uuidAccount       uint64 // 用户唯一id

	// todo: 移除测试代码
	battleUser *battleUser
}

func (agent *agentUser) String() string {
	return ""
}

// OnConnect 连接建立
func (agent *agentUser) OnConnect() {
	log.RunLogger.Printf("%v.OnConnect", agent.name)
}

// OnDisConnect 连接断开, 此事件处理完毕后, session将不可用
func (agent *agentUser) OnDisConnect() {
	log.RunLogger.Printf("%v.OnDisConnect", agent.name)

	// 离开MatchServer
	if agent.uuidAccount != uuid.InvalidUUID.Value() {
		proto := ffProto.ApplyProtoForSend(ffProto.MessageType_LeaveMatchServer)
		ffProto.SendProtoExtraDataUUID(instMatchServerClient, agent.UUID(), proto, false)
	}

	// 在战场, 则退出
	if agent.battleUser != nil {
		agent.battleUser.RunAway()
		agent.battleUser = nil
	}
}

// OnProto 收到Proto
func (agent *agentUser) OnProto(proto *ffProto.Proto) bool {
	log.RunLogger.Printf("%v.OnProto proto[%v]", agent.name, proto)

	protoID := proto.ProtoID()

	if callback, ok := mapAgentUserProtoCallback[protoID]; ok {
		return callback(agent, proto)
	}

	if result, ok := onBattleProto(agent, proto); ok {
		return result
	}

	log.FatalLogger.Printf("%v.OnProto not support protoID[%v]", agent.name, protoID)

	return false
}

// SendProtoExtraDataNormal 发送Proto
//	返回值仅表明请求发送的协议, 是否被添加到待发送管道内, 不代表一定能发送到对端
func (agent *agentUser) SendProtoExtraDataNormal(proto *ffProto.Proto) bool {
	return agent.netsession.SendProtoExtraDataNormal(proto)
}

// UUID 本次连接唯一标识
func (agent *agentUser) UUID() uuid.UUID {
	return agent.netsession.UUID()
}

// Init 初始化
func (agent *agentUser) Init(netsession netmanager.INetSession) {
	agent.name = fmt.Sprintf("agentUser[%v]", netsession.UUID())
	agent.netsession = netsession

	// 用户数据
	agent.uuidPlatformLogin = ""
	agent.uuidAccount = uuid.InvalidUUID.Value()
}

// Back 回收
func (agent *agentUser) Back() {
	agent.netsession = nil
}

// Close 主动关闭
func (agent *agentUser) Close() {
	agent.netsession.Close()
}

func newAgentUser() *agentUser {
	return &agentUser{}
}
