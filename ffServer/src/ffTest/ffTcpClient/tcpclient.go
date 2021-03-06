package main

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpclient"
	"ffCommon/util"
	"ffProto"
	"fmt"
	"sync/atomic"
	"time"
)

// 协议回调函数
//	返回值表明接收到的Proto是否进入了发送逻辑
var mapProtoCallback = map[ffProto.MessageType]func(agent *tcpClient, proto *ffProto.Proto) bool{
	ffProto.MessageType_EnterGameWorld: onProtoEnterGameWorld,
	ffProto.MessageType_KeepAlive:      onProtoKeepAlive,
}

type tcpClient struct {
	c base.Client

	chNewSession   chan base.Session
	chClientClosed chan struct{}

	sess              base.Session
	chSendProto       chan *ffProto.Proto
	chNetEventData    chan base.NetEventData
	sendExtraDataType ffProto.ExtraDataType

	number             int32
	keepAliveStartTime time.Time
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
func (client *tcpClient) mainLoopEnd(isPanic bool) {
	atomic.AddInt32(&waitQuitCount, -1)
}

func (client *tcpClient) onNetEventData(data base.NetEventData) {
	defer data.Back()

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
		client.onProto(data)
	}
}
func (client *tcpClient) onProto(data base.NetEventData) {
	log.RunLogger.Printf("tcpClient.onProto data[%v]: %v", data, client)

	changedToSendState := false
	proto := data.Proto()
	protoID := proto.ProtoID()

	// 如果协议在处理完毕后, 未进入发送逻辑, 则回收
	defer func() {
		if !changedToSendState {
			proto.BackAfterDispatch()
		}
	}()

	if err := proto.Unmarshal(); err != nil {
		log.FatalLogger.Printf("tcpClient.onProto proto[%v] Unmarshal error[%v]: %v", proto, err, client)
		client.close()
		return
	}

	log.RunLogger.Printf("tcpClient.onProto proto[%v]", proto)

	if callback, ok := mapProtoCallback[protoID]; ok {
		changedToSendState = callback(client, proto)
	} else {
		log.FatalLogger.Printf("tcpClient.onProto unknown protoID[%v]: %v", protoID, client)
	}
}

func onProtoEnterGameWorld(client *tcpClient, proto *ffProto.Proto) bool {
	msgEnterGameWorld, _ := proto.Message().(*ffProto.MsgEnterGameWorld)
	if msgEnterGameWorld.Result != ffError.ErrNone.Code() {
		log.RunLogger.Printf("tcpClient.onProtoEnterGameWorld Result[%v]", ffError.ErrByCode(msgEnterGameWorld.Result))
		client.close()
		return false
	}

	client.number = 1
	client.keepAliveStartTime = time.Now()

	proto = ffProto.ApplyProtoForSend(ffProto.MessageType_KeepAlive)
	message, _ := proto.Message().(*ffProto.MsgKeepAlive)
	message.Number = client.number

	client.sendProto(proto)
	return false
}

func onProtoKeepAlive(client *tcpClient, proto *ffProto.Proto) bool {
	message, _ := proto.Message().(*ffProto.MsgKeepAlive)
	if client.number != message.Number {
		log.RunLogger.Printf("tcpClient.onProtoKeepAlive number not match[%v-%v]", message.Number, client.number)
		client.close()
		return false
	} else if message.Number%10 == 0 {
		nanosecond := time.Now().Sub(client.keepAliveStartTime)
		fmt.Printf("average go-back net lag is %v %v %v\n", nanosecond, message.Number, uint64(nanosecond.Nanoseconds()/time.Microsecond.Nanoseconds())/uint64(message.Number))
	}

	client.number++
	message.Number = client.number

	proto.ChangeLimitStateRecvToSend()
	client.sendProto(proto)
	return true
}

func (client *tcpClient) sendProto(proto *ffProto.Proto) {
	if client.sendExtraDataType == ffProto.ExtraDataTypeNormal {
		proto.SetExtraDataNormal()
	} else if client.sendExtraDataType == ffProto.ExtraDataTypeUUID {
		proto.SetExtraDataUUID(client.sess.UUID().Value())
	}

	client.chSendProto <- proto
}

func (client *tcpClient) close() {
	fmt.Printf("close\n")
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
