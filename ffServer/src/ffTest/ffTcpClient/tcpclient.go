package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpclient"
	"ffCommon/util"
	"ffProto"
	"sync/atomic"
)

type tcpClient struct {
	c base.Client

	chNewSession   chan base.Session
	chClientClosed chan struct{}

	sess              base.Session
	chSendProto       chan *ffProto.Proto
	chNetEventData    chan base.NetEventData
	sendExtraDataType ffProto.ExtraDataType
}

func (client *tcpClient) mainLoop(params ...interface{}) {
	atomic.AddInt32(&waitQuitCount, 1)
	for {
		select {
		case sess := <-client.chNewSession:
			client.sess = sess
			client.sess.Start(client.chSendProto, client.chNetEventData, ffProto.ExtraDataTypeNormal)
		case data := <-client.chNetEventData:
			client.onNetEventData(data)
		case <-chApplicationQuit:
			return
		}
	}
}
func (client *tcpClient) mainLoopEnd() {
	atomic.AddInt32(&waitQuitCount, -1)
}

func (client *tcpClient) onNetEventData(data base.NetEventData) {
	switch data.NetEventType() {
	case base.NetEventOn:
		{
			var targetServerID int32 = 1
			var uuid uint64 = 0x1234567890

			proto := ffProto.ApplyProtoForSend(ffProto.MessageType_EnterGameWorld)
			message := proto.Message().(*ffProto.MsgEnterGameWorld)
			message.ServerID = targetServerID
			message.UUIDLogin = uuid
			client.sendProto(proto)
		}
	case base.NetEventOff:
	case base.NetEventProto:
	}
}

func (client *tcpClient) sendProto(proto *ffProto.Proto) {
	if client.sendExtraDataType == ffProto.ExtraDataTypeNormal {
		proto.SetExtraDataNormal()
	} else if client.sendExtraDataType == ffProto.ExtraDataTypeUUID {
		proto.SetExtraDataUUID(client.sess.UUID().Value())
	}

	client.chSendProto <- proto
}

func (client *tcpClient) start(addr string) {
	c, err := tcpclient.NewClient(addr)
	if err != nil {
		log.FatalLogger.Printf("tcpClient.start error[%v]", err)
		return
	}

	client.c = c

	client.chNewSession = make(chan base.Session, 1)
	client.chClientClosed = make(chan struct{}, 1)

	client.chSendProto = make(chan *ffProto.Proto, 2)
	client.chNetEventData = make(chan base.NetEventData, 2)

	client.sendExtraDataType = ffProto.ExtraDataTypeNormal

	c.Start(client.chNewSession, client.chClientClosed)

	go util.SafeGo(client.mainLoop, client.mainLoopEnd)
}
