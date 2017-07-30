package main

import (
	"ffAutoGen/ffError"
	"ffProto"
	"sync"
	"time"

	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpclient"
	"ffCommon/util"
)

type clientAgent struct {
	client         base.Client            // client 实现了tcp连接的client端
	extraDataType  ffProto.ExtraDataType  // extraDataType 附加数据类型
	chNetEventData chan base.NetEventData // chNetEventData 接收client反馈的事件

	mutexWorking sync.RWMutex // mutexWorking 状态锁
	working      bool         // mutexWorking 状态

	chNtfClose          chan struct{} // chClose 通知关闭
	chWaitCloseComplete chan struct{} // chWaitCloseComplete 等待关闭完成

	number    int32
	startTime time.Time
}

// close 同步关闭
func (c *clientAgent) close(delayMillisecond int64) {
	close(c.chNtfClose)

	<-c.chWaitCloseComplete
}

// changeWorking 改变当前状态
func (c *clientAgent) changeWorking(working bool) {
	c.mutexWorking.Lock()
	defer c.mutexWorking.Lock()
	c.working = working
}

func (c *clientAgent) sendProto(p *ffProto.Proto) {
	c.mutexWorking.RLock()
	defer c.mutexWorking.RUnlock()

	p.SetExtraDataNormal()
	if c.working {
		c.client.SendProto(p)
		return
	}
	p.BackAfterSend()
}

func (c *clientAgent) onConnect(data base.NetEventData) {
	c.changeWorking(true)

	var targetServerID int32 = 1
	var uuid uint64 = 0x1234567890

	p := ffProto.ApplyProtoForSend(ffProto.MessageType_EnterGameWorld)
	message := p.Message().(*ffProto.MsgEnterGameWorld)
	message.ServerID = targetServerID
	message.UUIDLogin = uuid
	c.sendProto(p)
}

func (c *clientAgent) onDisConnect(data base.NetEventData) {
	c.changeWorking(false)

	if !data.ManualClose() {
		c.client.ReConnect()
	}
}

// onEnd 只会在主动调用client.Close, 且底层完成关闭处理后, 才会触发该事件. 处理完该事件后, client将被回收. 此处需要解除引用.
func (c *clientAgent) onEnd(data base.NetEventData) {
	c.client = nil
}

func (c *clientAgent) onProto(data base.NetEventData) {
	proto := data.Proto()
	protoID := proto.ProtoID()

	log.RunLogger.Printf("clientAgent.onProto: proto[%v]\n", proto)

	if protoID == ffProto.MessageType_EnterGameWorld {
		if err := proto.Unmarshal(); err != nil {
			log.FatalLogger.Println(err)
			return
		}

		c.onMessageEnterGameWorld(proto)

		return
	}
}

func (c *clientAgent) onMessageEnterGameWorld(proto *ffProto.Proto) {
	message, _ := proto.Message().(*ffProto.MsgEnterGameWorld)
	if message.Result != ffError.ErrNone.Code() {
		log.FatalLogger.Println(ffError.ErrByCode(message.Result))
		return
	}

	log.RunLogger.Printf("clientAgent.onMessageEnterGameWorld message[%v]\n", message)

	c.number = 0
	c.startTime = time.Now()

	p := ffProto.ApplyProtoForSend(ffProto.MessageType_KeepAlive)
	m, _ := p.Message().(*ffProto.MsgKeepAlive)
	m.Number = c.number
	c.sendProto(p)
}

func (c *clientAgent) onMessageKeepAlive(proto *ffProto.Proto) {
	message, _ := proto.Message().(*ffProto.MsgKeepAlive)

	log.RunLogger.Printf("clientAgent.onMessageKeepAlive message[%v]\n", message)

	c.number++
	if c.number != message.Number {
		log.FatalLogger.Printf("clientAgent.onMessageKeepAlive number not match[%v:%v]", c.number, message.Number)
	}

	c.number++

	message.Number = c.number
	c.sendProto(proto)
}

func (c *clientAgent) mainLoop(params ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			util.PrintPanicStack(err, "clientAgent.mainLoop", c)
		}

		c.chWaitCloseComplete <- struct{}{}
	}()

	c.client.Start(c.chNetEventData, ffProto.ExtraDataTypeNormal)

	for {
		select {
		case data := <-c.chNetEventData: // 网络事件数据
			if c.onNetDataEvent(data) {
				// 退出
				return
			}

		case <-c.chNtfClose: // 通知退出
			for {
				select {
				case data := <-c.chNetEventData: // 网络事件数据
					if c.onNetDataEvent(data) {
						// 退出
						return
					}
				}
			}
		}
	}
}

// onNetDataEvent 处理网络事件. 返回是否退出网络事件处理协程
func (c *clientAgent) onNetDataEvent(data base.NetEventData) bool {
	defer data.Back()

	switch data.NetEventType() {
	case base.NetEventOn:
		c.onConnect(data)
	case base.NetEventOff:
		c.onDisConnect(data)
	case base.NetEventProto:
		c.onProto(data)
	case base.NetEventEnd:
		c.onEnd(data)
		return true
	}
	return false
}

func newClientAgent(addr string, extraDataType ffProto.ExtraDataType) (agent *clientAgent, err error) {
	client, err := tcpclient.NewClient(addr)
	if err != nil {
		return nil, err
	}

	agent = &clientAgent{
		extraDataType:  extraDataType,
		client:         client,
		chNetEventData: make(chan base.NetEventData, tcpclient.DefaultClientNetEventDataChanCount),

		chNtfClose:          make(chan struct{}),
		chWaitCloseComplete: make(chan struct{}, 1),
	}
	return
}
