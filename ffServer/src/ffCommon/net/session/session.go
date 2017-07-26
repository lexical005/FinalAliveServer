package session

import (
	"io"
	"net"
	"sync"
	"time"

	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"
)

// Session is a connnect between Server and Client
type Session struct {
	// 数据发送相关
	sendLeft        []byte               // 上次发送剩余字节流
	sendProtoHeader *ffProto.ProtoHeader // 创建Session时创建，固化

	// 数据接收相关
	recvHeaderBuf   []byte               // 接收协议头所需的缓冲区
	recvProtoHeader *ffProto.ProtoHeader // 创建Session时创建，固化

	conn  net.Conn   // 底层连接
	agent base.Agent // 上层逻辑对象
	uuid  uuid.UUID  // UUID

	chSendPool chan *ffProto.Proto // 待发送协议管道
	chExit     chan bool           // 发送/接收协程退出控制
	chClose    chan bool           // 通知发送/接收协程退出
	onceClose  util.Once           // 用于只执行一次关闭
}

// UUID UUID
func (s *Session) UUID() uuid.UUID {
	return s.uuid
}

// Close Session
// close send and recv goroutine(SendProto is called by high logic)
// delayMillisecond: 延迟多少毫秒关闭
func (s *Session) Close(delayMillisecond int64) {
	log.RunLogger.Printf("Session.Close: Session[%x] delayMillisecond[%d]", s.UUID(), delayMillisecond)
	if delayMillisecond > 0 {
		go util.SafeGo(func(params ...interface{}) {
			select {
			case <-time.After(time.Duration(delayMillisecond) * time.Millisecond):
				s.onceClose.Do(func() {
					// 关闭结束管道，触发发送/接收协程退出
					close(s.chClose)

					s.doClose()
				})
			}
		})
	} else {
		s.onceClose.Do(func() {
			// 关闭结束管道，触发发送/接收协程退出
			close(s.chClose)

			s.doClose()
		})
	}
}

// SendProto 将待发送协议追加到待发送管道
func (s *Session) SendProto(p *ffProto.Proto) {
	// 追加到待发送数据的末尾
	s.chSendPool <- p
}

// Start 主循环
func (s *Session) Start(conn net.Conn, agent base.Agent, recvProtoExtraDataType ffProto.ExtraDataType) {
	log.RunLogger.Printf("Session.Start: Session[%x]", s.UUID())

	// 重新开始
	if s.sendLeft != nil {
		s.sendLeft = s.sendLeft[:0]
	} else {
		s.sendLeft = make([]byte, 0)
	}
	s.sendProtoHeader.ResetForSend()
	s.recvProtoHeader.ResetForRecv(recvProtoExtraDataType)

	s.chSendPool = make(chan *ffProto.Proto, 4)
	s.chExit = make(chan bool, 2)
	s.chClose = make(chan bool)
	s.onceClose.Reset()

	s.conn, s.agent = conn, agent

	s.agent.OnEvent(ffProto.MessageType_MT_Connect, nil)

	// start recv goroutine
	go util.SafeGo(s.doRecvWork)

	// start send goroutine, and send goroutine is main goroutine of session: clear session is called on send goroutine exit
	go util.SafeGo(s.doSendWork)
}

// WaitCloseChan 外界调用此接口后，将返回一个管道，供外界等候Session关闭
func (s *Session) WaitCloseChan() <-chan bool {
	return s.chClose
}

func (s *Session) doRecvWork(params ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			util.PrintPanicStack(err, "Session.doRecvWork")
		}

		s.chExit <- true

		s.onceClose.Do(func() {
			// 关闭结束管道，触发发送/接收协程退出
			close(s.chClose)

			s.doClose()
		})
	}()

	var err error
	for {
		// 接收
		if err = s.doRecv(); err != nil {
			log.RunLogger.Printf("Session.doRecvWork: Session[%x] err[%v]", s.UUID(), err)
			break
		}

		// 是否关闭
		select {
		case <-s.chClose:
			return
		default:
		}
	}
}

func (s *Session) doSendWork(params ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			util.PrintPanicStack(err, "Session.doSendWork")
		}

		s.chExit <- true

		s.onceClose.Do(func() {
			// 关闭结束管道，触发发送/接收协程退出
			close(s.chClose)

			s.doClose()
		})
	}()

	for {
		select {
		case <-s.chClose:
			// 关闭
			return
		case p := <-s.chSendPool:
			// 发送
			if p == nil || !s.doSend(p) {
				return
			}
		}
	}
}

func (s *Session) doSend(p *ffProto.Proto) bool {
	defer ffProto.BackProtoAfterSend(p)

	err := p.Marshal(s.sendProtoHeader)
	if err != nil {
		log.RunLogger.Printf("Session.doSend: Session[%x] err[%v]", s.UUID(), err)
		return false
	}

	buf := p.BytesForSend()

	// 之前的发送有遗留
	if len(s.sendLeft) > 0 {
		s.sendLeft = append(s.sendLeft, buf...)
		buf = s.sendLeft
	}

	// 发送
	n, err := s.conn.Write(buf)

	log.RunLogger.Printf("Session.doSend: Session[%x] send buf[%v] send real[%d] err[%v]", s.UUID(), buf, n, err)

	// 全部发送
	if n == len(buf) {
		s.sendLeft = s.sendLeft[:0]
		return true
	}

	log.RunLogger.Printf("Session.doSend: Session[%x] send expect[%v] send real[%v] err[%v]", s.UUID(), len(buf), n, err)

	// 发送了部分
	if n > 0 {
		s.sendLeft = s.sendLeft[:0]
		s.sendLeft = append(s.sendLeft, buf[n:]...)
		return true
	}

	// 一点都没发送出去, 则认为连接出现了问题
	return false
}

func (s *Session) doRecv() error {
	// solve dead link problem:
	// physical disconnection without any communcation between client and server
	// will cause the read to block FOREVER, so a timeout is a rescue.
	s.conn.SetReadDeadline(time.Now().Add(time.Duration(readDeadtime) * time.Second))

	// read Proto header
	_, err := io.ReadFull(s.conn, s.recvHeaderBuf)
	if err != nil {
		return err
	}

	log.RunLogger.Printf("Session.doRecv: Session[%x] recvHeaderBuf[%v]", s.UUID(), s.recvHeaderBuf)

	// 协议头解析
	err = s.recvProtoHeader.Unmarshal(s.recvHeaderBuf)
	if err != nil {
		return err
	}

	// 接收剩余部分
	p := ffProto.ApplyProtoForRecv(s.recvProtoHeader)
	defer ffProto.BackProtoAfterRecv(p)

	buf := p.BytesForRecv()
	_, err = io.ReadFull(s.conn, buf)
	if err != nil {
		return err
	}

	log.RunLogger.Printf("Session.doRecv: Session[%x] recvProtoContent[%v]", s.UUID(), buf)

	// 数据接收完毕, 通知校验
	err = p.OnRecvAllBytes(s.recvProtoHeader)
	if err != nil {
		return err
	}

	// 通知 agent(此处不立即反序列化，由agent选择合适的时机自行调用Proto的反序列化方法)
	s.agent.OnEvent(p.ProtoID(), p)

	return nil
}

func (s *Session) doClose() {
	log.RunLogger.Printf("Session.doClose: Session[%x]", s.UUID())

	// 关闭底层连接
	if s.conn != nil {
		s.conn.Close()
	}

	// 等待发送和接收协程退出
	<-s.chExit
	<-s.chExit

	// 通知上层逻辑对象
	var wg sync.WaitGroup
	wg.Add(1)
	s.agent.OnEvent(ffProto.MessageType_MT_DisConnect, &wg)
	wg.Wait()

	// 清理内部数据
	close(s.chSendPool)
	for p := range s.chSendPool {
		ffProto.BackProtoAfterSend(p)
	}
	s.chSendPool = nil
	s.sendLeft = s.sendLeft[:0]
	s.conn = nil
	s.agent = nil

	// 归还
	sp.back(s)
}

func newSession() base.Session {
	return &Session{
		sendProtoHeader: ffProto.NewProtoHeader(),

		recvHeaderBuf:   ffProto.NewProtoHeaderBuf(),
		recvProtoHeader: ffProto.NewProtoHeader(),
	}
}
