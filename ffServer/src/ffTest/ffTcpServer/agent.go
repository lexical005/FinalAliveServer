package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/util"
	"ffProto"

	"sync"
	"time"

	"github.com/golang/protobuf/proto"
)

type agent struct {
	muSessionOn sync.RWMutex
	sess        base.Session
}

func (a *agent) SendProto(p *ffProto.Proto) {
	a.muSessionOn.RLock()
	defer a.muSessionOn.RUnlock()

	if a.sess != nil {
		a.sess.SendProto(p)
	}
}

func (a *agent) onDisConnect() {
	a.muSessionOn.Lock()
	defer a.muSessionOn.Unlock()

	a.sess = nil
}

func (a *agent) OnEvent(protoID ffProto.MessageType, data interface{}) {
	log.RunLogger.Printf("agent.OnEvent protoID[%s]\n", ffProto.MessageType_name[int32(protoID)])

	if protoID == ffProto.MessageType_MT_Connect {
		// p := ffProto.ApplyProtoForSend(ffProto.MessageType_MT_MsgChatData)
		// msg, _ := p.MessageForSend().(*ffProto.MsgChatData)
		// msg.MsgData = proto.String("Proto")
		// msg.FromName = proto.String("FromName")
		// msg.ChannelType = proto.Uint32(32)

		// a.SendProto(p)
	} else if protoID == ffProto.MessageType_MT_DisConnect {
		a.onDisConnect()

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

			a.SendProto(p2)
		})
	}
}
