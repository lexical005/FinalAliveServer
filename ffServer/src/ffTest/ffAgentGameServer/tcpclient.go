package main

import (
	"ffProto"
	"sync"
	// "time"

	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpclient"
	// "ffCommon/util"
	// "github.com/golang/protobuf/proto"
)

type clientAgent struct {
	client base.Client

	isSessionOn bool
	muSessionOn sync.RWMutex
}

func (c *clientAgent) sendProto(p *ffProto.Proto) {
	c.muSessionOn.RLock()
	defer c.muSessionOn.RUnlock()

	if c.isSessionOn {
		c.client.SendProto(p)
	}
}

func (c *clientAgent) onConnect() {
	c.isSessionOn = true

	var serverType string = "GameServer"
	var serverID int32 = 1

	p := ffProto.ApplyProtoForSend(ffProto.MessageType_ServerRegister)
	message := p.Message().(*ffProto.MsgServerRegister)
	message.ServerType = &serverType
	message.ServerID = &serverID
	c.sendProto(p)
}

func (c *clientAgent) onDisConnect() {
	c.muSessionOn.Lock()
	defer c.muSessionOn.Unlock()

	c.isSessionOn = false
}

func (c *clientAgent) OnEvent(protoID ffProto.MessageType, data interface{}) {
	log.RunLogger.Printf("clientAgent.OnEvent: protoID[%s]\n", ffProto.MessageType_name[int32(protoID)])

	if protoID == ffProto.MessageType_SessionConnect {
		c.onConnect()
	} else if protoID == ffProto.MessageType_SessionDisConnect {
		c.onDisConnect()

		wg, _ := data.(*sync.WaitGroup)
		wg.Done()
	} else {
		proto := data.(*ffProto.Proto)
		proto.Unmarshal()
		log.RunLogger.Println(proto)
		// p, _ := data.(*ffProto.Proto)
		// m, err := p.Unmarshal()
		// if err != nil {
		//  log.RunLogger.Println(err)
		//  return
		// }

		// log.RunLogger.Println(m)

		// if p.ProtoID() != ffProto.MessageType_ChatData {
		//  log.RunLogger.Printf("recv invalid ProtID: ProtoID[%d]\n", p.ProtoID())
		//  return
		// }

		// m1 := m.(*ffProto.MsgChatData)

		// p2 := ffProto.ApplyProtoForSend(ffProto.MessageType_ChatData)
		// msg, _ := p2.MessageForSend().(*ffProto.MsgChatData)
		// msg.MsgData = proto.String(m1.GetMsgData())
		// msg.FromName = proto.String(m1.GetFromName())
		// msg.ChannelType = proto.Uint32(m1.GetChannelType())

		// go util.SafeGo(func(params ...interface{}) {
		//  <-time.After(5 * time.Second)

		//  c.SendProto(p2)
		// })
	}
}

func (c *clientAgent) start(addr string, autoReconnect bool) (err error) {
	c.client, err = tcpclient.New(addr, autoReconnect)
	c.client.Start(c)
	return
}

var agent clientAgent
