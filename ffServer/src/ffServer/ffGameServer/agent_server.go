package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpclient"
	"ffCommon/util"

	"ffProto"
	"sync"
	"sync/atomic"
	"time"
)

type agentServer struct {
	server      base.Client // 与AgentServer的连接
	serverState int32       // 与AgentServer的连接的状态. 0 初始状态; -1 未连接; 1 已连接
}

// 与目标服务器的连接建立
func (as *agentServer) onConnect() {
	// 服务器注册
	serverID := int32(appConfig.Server.ServerID)
	p := ffProto.ApplyProtoForSend(ffProto.MessageType_ServerRegister)
	message := p.Message().(*ffProto.MsgServerRegister)
	message.ServerType = appConfig.Server.ServerType
	message.ServerID = serverID
	p.SetExtraData(ffProto.ExtraDataTypeUUID, 0)
	as.server.SendProto(p)

	// 然后服务器再标识为连接在线
	atomic.StoreInt32(&as.serverState, 1)
}

func (as *agentServer) onDisConnect() {
	atomic.StoreInt32(&as.serverState, -1)

	worldFrame.onAgentServerDisConnect()
}

func (as *agentServer) OnEvent(protoID ffProto.MessageType, data interface{}) {
	if protoID == ffProto.MessageType_SessionConnect {

		log.RunLogger.Printf("agentServer.OnEvent: protoID[%s]\n",
			ffProto.MessageType_name[int32(protoID)])

		as.onConnect()

	} else if protoID == ffProto.MessageType_SessionDisConnect {

		log.RunLogger.Printf("agentServer.OnEvent: protoID[%s]\n",
			ffProto.MessageType_name[int32(protoID)])

		as.onDisConnect()

		wg, _ := data.(*sync.WaitGroup)
		wg.Done()


	} else {
		p, _ := data.(*ffProto.Proto)

		log.RunLogger.Printf("agentServer.OnEvent: protoID[%s] proto[%v]\n",
			ffProto.MessageType_name[int32(protoID)], p)

		worldFrame.onRecvProto(p)
	}
}

// 保持现有连接不断开
func (as *agentServer) keepAliveLoop(params ...interface{}) {
	interval := time.Duration(appConfig.Session.ReadDeadTime*40/100) * time.Second

	for {
		select {
		case <-time.After(interval):
			as.sendProto(ffProto.ApplyProtoForSend(ffProto.MessageType_ServerKeepAlive), 0)
		}
	}
}

// 由具体逻辑调用, 发送协议
func (as *agentServer) sendProto(p *ffProto.Proto, extraDataUUID uint64) bool {
	p.SetExtraData(ffProto.ExtraDataTypeUUID, extraDataUUID)

	if atomic.CompareAndSwapInt32(&as.serverState, 1, 1) {
		log.RunLogger.Printf("agentServer.sendProto: uuidAgent[%x] proto[%v] serverState[%d]",
			extraDataUUID, p, 1)

		as.server.SendProto(p)
		return true
	}

	log.RunLogger.Printf("agentServer.sendProto: uuidAgent[%x] proto[%v] serverState[%d]",
		extraDataUUID, p, -1)

	p.BackAfterSend()
	return false
}

// 根据配置创建Server
func (as *agentServer) init() error {
	server, err := tcpclient.New(appConfig.AgentServer.ListenAddr, appConfig.AgentServer.AutoReConnect)
	if err != nil {
		return err
	}

	as.server = server

	return nil
}

// Server启动监听
func (as *agentServer) start() error {
	// 新起一个协程, 用于定时向已建立的连接发送KeepAlive协议, 以避免连接被自然断开
	go util.SafeGo(as.keepAliveLoop, nil)

	return as.server.Start(as, ffProto.ExtraDataTypeUUID)
}
