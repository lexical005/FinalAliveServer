package main

import (
	"ffAutoGen/ffError"
	"ffProto"
	"sync"
	"time"

	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpclient"
)

type clientAgent struct {
	client base.Client

	isSessionOn bool
	muSessionOn sync.RWMutex

	number    int32
	startTime time.Time
}

func (c *clientAgent) sendProto(p *ffProto.Proto) {
	c.muSessionOn.RLock()
	defer c.muSessionOn.RUnlock()

	p.SetExtraData(ffProto.ExtraDataTypeNormal, 0)
	if c.isSessionOn {
		c.client.SendProto(p)
		return
	}
	ffProto.BackProtoAfterSend(p)
}

func (c *clientAgent) onConnect() {
	c.isSessionOn = true

	var targetServerID int32 = 1
	var uuid uint64 = 0x123456

	p := ffProto.ApplyProtoForSend(ffProto.MessageType_EnterGameWorld)
	message := p.Message().(*ffProto.MsgEnterGameWorld)
	message.ServerID = &targetServerID
	message.UUIDAccount = &uuid
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
		if protoID == ffProto.MessageType_EnterGameWorld {
			proto := data.(*ffProto.Proto)
			if err := proto.Unmarshal(); err != nil {
				log.FatalLogger.Println(err)
				return
			}

			message, _ := proto.Message().(*ffProto.MsgEnterGameWorld)
			if *message.Result != *ffError.ErrNone.Code() {
				log.FatalLogger.Println(ffError.ErrByCode(*message.Result))
				return
			}
			// c.number = 0
			// c.startTime = time.Now()

			// p := ffProto.ApplyProtoForSend(ffProto.MessageType_CountNumber)
			// m, _ := p.Message().(*ffProto.MsgCountNumber)
			// m.Number = &c.number
			// c.sendProto(p)
			return
		}

		// if protoID == ffProto.MessageType_CountNumber {
		// 	proto := data.(*ffProto.Proto)
		// 	if err := proto.Unmarshal(); err != nil {
		// 		log.RunLogger.Println(err)
		// 	}

		// 	message, _ := proto.Message().(*ffProto.MsgCountNumber)
		// 	if *message.Number != c.number+1 {
		// 		log.FatalLogger.Printf("message.Number[%d] != c.number[%d]+1", *message.Number, c.number)
		// 	}
		// 	c.number++
		// 	if c.number < 1000 {
		// 		message.Number = &c.number
		// 		c.sendProto(proto)
		// 		return
		// 	}
		// 	delta := time.Now().Sub(c.startTime)
		// 	// fmt.Println(c.number, delta.Nanoseconds(), delta)
		// 	logfile.Init(logfile.DefaultLogFileRelativePath, logfile.DefaultLogFileLengthLimit, true, "run", true, logfile.DefaultLogFileFatalPrefix)
		// 	log.RunLogger.Println(c.number, delta.Nanoseconds(), delta)
		// }
	}
}

func (c *clientAgent) start(addr string, autoReconnect bool) (err error) {
	log.RunLogger = log.NewLoggerEmpty()
	// logfile.Init(logfile.DefaultLogFileRelativePath, logfile.DefaultLogFileLengthLimit, true, "run", true, logfile.DefaultLogFileFatalPrefix)
	c.client, err = tcpclient.New(addr, autoReconnect)
	c.client.Start(c, ffProto.ExtraDataTypeNormal)
	return
}

var agent clientAgent
