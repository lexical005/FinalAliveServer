package main

import (
	"ffProto"
	"sync"
	"time"

	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpclient"
	"ffCommon/util"

	"github.com/golang/protobuf/proto"
)

type tcpclient1 struct {
	client base.Client

	isSessionOn bool
	muSessionOn sync.RWMutex
}

func (c *tcpclient1) SendProto(p *ffProto.Proto) {
	c.muSessionOn.RLock()
	defer c.muSessionOn.RUnlock()

	if c.isSessionOn {
		c.client.SendProto(p)
	}
}

func (c *tcpclient1) onDisConnect() {
	c.muSessionOn.Lock()
	defer c.muSessionOn.Unlock()

	c.isSessionOn = false
}

func (c *tcpclient1) OnEvent(protoID ffProto.MessageType, data interface{}) {
	log.RunLogger.Printf("tcpclient1.OnEvent protoID[%s]\n", ffProto.MessageType_name[int32(protoID)])

	if protoID == ffProto.MessageType_MT_Connect {
		c.isSessionOn = true

		p := ffProto.ApplyProtoForSend(ffProto.MessageType_MT_MsgChatData)
		msg, _ := p.MessageForSend().(*ffProto.MsgChatData)
		msg.MsgData = proto.String("Proto")
		msg.FromName = proto.String("FromName")
		msg.ChannelType = proto.Uint32(32)

		c.SendProto(p)
	} else if protoID == ffProto.MessageType_MT_DisConnect {
		c.onDisConnect()

		wg, _ := data.(*sync.WaitGroup)
		wg.Done()
	} else {
		p, _ := data.(*ffProto.Proto)
		m, err := p.Unmarshal()
		if err != nil {
			log.RunLogger.Println(err)
			return
		}

		log.RunLogger.Println(m)

		if p.ProtoID() != ffProto.MessageType_MT_MsgChatData {
			log.RunLogger.Printf("recv invalid ProtID: ProtoID[%d]\n", p.ProtoID())
			return
		}

		m1 := m.(*ffProto.MsgChatData)

		p2 := ffProto.ApplyProtoForSend(ffProto.MessageType_MT_MsgChatData)
		msg, _ := p2.MessageForSend().(*ffProto.MsgChatData)
		msg.MsgData = proto.String(m1.GetMsgData())
		msg.FromName = proto.String(m1.GetFromName())
		msg.ChannelType = proto.Uint32(m1.GetChannelType())

		go util.SafeGo(func(params ...interface{}) {
			<-time.After(5 * time.Second)

			c.SendProto(p2)
		})
	}
}

func (c *tcpclient1) start(addr string, autoReconnect bool) (err error) {
	c.client, err = tcpclient.New(addr, autoReconnect)
	c.client.Start(c)
	return
}

var tc1 tcpclient1
