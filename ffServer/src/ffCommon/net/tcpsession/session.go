package tcpsession

import (
	"fmt"
	"io"
	"net"
	"time"

	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"
)

// tcpSession 客户端与服务端之间的连接
type tcpSession struct {
	// 数据发送相关
	sendLeft        []byte               // 上次发送剩余字节流, 固化
	sendProtoHeader *ffProto.ProtoHeader // 创建tcpSession时创建, 固化

	// 数据接收相关
	recvHeaderBuf   []byte               // 接收协议头所需的缓冲区, 固化
	recvProtoHeader *ffProto.ProtoHeader // 创建tcpSession时创建, 固化

	conn           net.Conn               // 底层连接
	chNetEventData chan base.NetEventData // 外界接收事件数据的管道

	uuid    uuid.UUID // UUID
	working bool      // 是否正常工作状态

	chSendPool                  chan *ffProto.Proto // 待发送协议管道
	chWaitRecvSendGoroutineExit chan struct{}       // 等待发送/接收协程退出
	chNtfRecvSendGoroutineExit  chan struct{}       // 通知发送/接收协程退出
	onceClose                   util.Once           // 用于只执行一次关闭tcpsession
}

// UUID UUID
func (s *tcpSession) UUID() uuid.UUID {
	return s.uuid
}

// Close 异步关闭Session
// delayMillisecond: 延迟多少毫秒关闭
func (s *tcpSession) Close(delayMillisecond int64) {
	log.RunLogger.Printf("tcpSession.Close delayMillisecond[%d]: %v", delayMillisecond, s)

	// 立即标识停止工作
	s.working = false

	go util.SafeGo(func(params ...interface{}) {
		if delayMillisecond > 0 {
			select {
			case <-time.After(time.Duration(delayMillisecond) * time.Millisecond):
				s.onceClose.Do(func() {
					s.doClose(true)
				})
			}
		} else {
			s.onceClose.Do(func() {
				s.doClose(true)
			})
		}
	})
}

// SendProto 发送Proto到对端, 外界只应该在收到连接建立完成事件之后再调用此接口, 异步
func (s *tcpSession) SendProto(proto *ffProto.Proto) {
	log.RunLogger.Printf("tcpSession.SendProto proto[%v]: %v", proto, s)

	s.chSendPool <- proto
}

// Start 主循环
func (s *tcpSession) Start(conn net.Conn, chNetEventData chan base.NetEventData, recvProtoExtraDataType ffProto.ExtraDataType) {
	// 重新开始
	if s.sendLeft != nil {
		s.sendLeft = s.sendLeft[:0]
	} else {
		s.sendLeft = make([]byte, 0)
	}
	s.sendProtoHeader.ResetForSend()
	s.recvProtoHeader.ResetForRecv(recvProtoExtraDataType)

	s.chSendPool = make(chan *ffProto.Proto, 4)
	s.chWaitRecvSendGoroutineExit = make(chan struct{}, 2)
	s.chNtfRecvSendGoroutineExit = make(chan struct{})
	s.onceClose.Reset()

	s.conn, s.chNetEventData = conn, chNetEventData

	s.working = true

	log.RunLogger.Printf("tcpSession.Start: %v", s)

	s.chNetEventData <- newSessionNetEventOn(s)

	// start recv goroutine
	go util.SafeGo(s.mainRecv)

	// start send goroutine
	go util.SafeGo(s.mainSend)
}

func (s *tcpSession) String() string {
	return fmt.Sprintf(`uuid[%v] working[%v] sendLeft[%v] chSendPool[%v]`,
		s.uuid, s.working, len(s.sendLeft), len(s.chSendPool))
}

func (s *tcpSession) mainSend(params ...interface{}) {
	defer func() {
		log.RunLogger.Printf("tcpSession.mainSend end: %v", s)

		if err := recover(); err != nil {
			util.PrintPanicStack(err, "tcpSession.mainSend", s)
		}

		s.chWaitRecvSendGoroutineExit <- struct{}{}

		s.onceClose.Do(func() {
			s.doClose(false)
		})
	}()

	for {
		select {
		case <-s.chNtfRecvSendGoroutineExit:
			// 关闭
			return
		case p := <-s.chSendPool:
			// 发送
			if p == nil || !s.doSend(p) {
				return
			}
		}

		// 有未发送完毕的数据, 且当前没有等待发送的协议
		for len(s.sendLeft) > 0 && len(s.chSendPool) == 0 {
			// 等待2毫秒
			<-time.After(2 * time.Microsecond)

			if !s.doSendBuf(s.sendLeft) {
				return
			}

			select {
			case <-s.chNtfRecvSendGoroutineExit:
				// 关闭
				return
			default:
				break
			}
		}
	}
}

func (s *tcpSession) doSend(p *ffProto.Proto) bool {
	defer p.BackAfterSend()

	err := p.Marshal(s.sendProtoHeader)
	if err != nil {
		log.RunLogger.Printf("tcpSession.doSend err[%v]: %v", err, s)
		return false
	}

	buf := p.BytesForSend()

	// 之前的发送有遗留
	if len(s.sendLeft) > 0 {
		s.sendLeft = append(s.sendLeft, buf...)
		buf = s.sendLeft
	}

	return s.doSendBuf(buf)
}

func (s *tcpSession) doSendBuf(buf []byte) bool {
	s.sendLeft = s.sendLeft[:0]

	// 发送
	n, err := s.conn.Write(buf)

	log.RunLogger.Printf("tcpSession.doSend send buf[%v] send real[%d] err[%v]: %v", buf, n, err, s)

	// 全部发送
	if n == len(buf) {
		return true
	}

	log.RunLogger.Printf("tcpSession.doSend send expect[%v] send real[%v] err[%v]: %v", len(buf), n, err, s)

	// 发送了部分
	if n > 0 {
		s.sendLeft = append(s.sendLeft, buf[n:]...)
		return true
	}

	// 一点都没发送出去, 则认为连接出现了问题
	return false
}

func (s *tcpSession) mainRecv(params ...interface{}) {
	defer func() {
		log.RunLogger.Printf("tcpSession.mainRecv end: %v", s)

		if err := recover(); err != nil {
			util.PrintPanicStack(err, "tcpSession.mainRecv", s)
		}

		s.chWaitRecvSendGoroutineExit <- struct{}{}

		s.onceClose.Do(func() {
			s.doClose(false)
		})
	}()

	var err error
	for {
		// 接收
		if err = s.doRecv(); err != nil {
			log.RunLogger.Printf("tcpSession.mainRecv err[%v]: %v", err, s)
			break
		}

		// 是否关闭
		select {
		case <-s.chNtfRecvSendGoroutineExit:
			return
		default:
		}
	}
}

func (s *tcpSession) doRecv() error {
	// solve dead link problem:
	// physical disconnection without any communcation between client and server
	// will cause the read to block FOREVER, so a timeout is a rescue.
	s.conn.SetReadDeadline(time.Now().Add(time.Duration(readDeadtime) * time.Second))

	// read Proto header
	_, err := io.ReadFull(s.conn, s.recvHeaderBuf)
	if err != nil {
		return err
	}

	log.RunLogger.Printf("tcpSession.doRecv recvHeaderBuf[%v]: %v", s.recvHeaderBuf, s)

	// 协议头解析
	err = s.recvProtoHeader.Unmarshal(s.recvHeaderBuf)
	if err != nil {
		return err
	}

	// 接收剩余部分
	p := ffProto.ApplyProtoForRecv(s.recvProtoHeader)

	buf := p.BytesForRecv()
	_, err = io.ReadFull(s.conn, buf)
	if err != nil {
		return err
	}

	log.RunLogger.Printf("tcpSession.doRecv recvProtoData[%v]: %v", buf, s)

	// 数据接收完毕, 通知校验
	err = p.OnRecvAllBytes(s.recvProtoHeader)
	if err != nil {
		return err
	}

	// 协议事件
	s.chNetEventData <- newSessionNetEventProto(s, p)

	return nil
}

// doClose Session本次有效期间, 只会被执行一次
func (s *tcpSession) doClose(manual bool) {
	log.RunLogger.Printf("tcpSession.doClose: %v", s)

	// 连接异常导致关闭时, 标识停止工作
	if !manual {
		s.working = false
	}

	// 关闭结束管道, 触发发送/接收协程退出
	close(s.chNtfRecvSendGoroutineExit)

	// 关闭底层连接
	if s.conn != nil {
		s.conn.Close()
	}

	// 连接断开事件
	s.chNetEventData <- newSessionNetEventOff(s, manual)

	// 等待发送和接收协程退出
	<-s.chWaitRecvSendGoroutineExit
	<-s.chWaitRecvSendGoroutineExit

	// 连接结束事件
	s.chNetEventData <- newSessionNetEventEnd(s)
}

// back 外界已停止引用Session, 可安全回收了
func (s *tcpSession) back() {
	log.RunLogger.Printf("tcpSession.back: %v", s)

	// 清理内部数据
	close(s.chSendPool)
	for p := range s.chSendPool {
		p.BackAfterSend()
	}
	s.chSendPool = nil
	s.sendLeft = s.sendLeft[:0]
	s.conn, s.chNetEventData = nil, nil

	// 归还
	sessPool.back(s)
}

func newSession() base.Session {
	return &tcpSession{
		sendProtoHeader: ffProto.NewProtoHeader(),

		recvHeaderBuf:   ffProto.NewProtoHeaderBuf(),
		recvProtoHeader: ffProto.NewProtoHeader(),
	}
}
