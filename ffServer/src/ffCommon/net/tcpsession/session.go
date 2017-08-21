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

	uuid uuid.UUID // UUID

	chSendProto                 chan *ffProto.Proto // 待发送协议管道
	chWaitRecvSendGoroutineExit chan struct{}       // 等待发送/接收协程退出
	chNtfRecvSendGoroutineExit  chan struct{}       // 通知发送/接收协程退出
	onceClose                   util.Once           // 用于只执行一次关闭tcpsession

	manualClose bool // 是否主动断开连接
}

// UUID UUID
func (s *tcpSession) UUID() uuid.UUID {
	return s.uuid
}

// Start 主循环
func (s *tcpSession) Start(chSendProto chan *ffProto.Proto, chNetEventData chan base.NetEventData, recvProtoExtraDataType ffProto.ExtraDataType) {
	// 重新开始
	if s.sendLeft != nil {
		s.sendLeft = s.sendLeft[:0]
	} else {
		s.sendLeft = make([]byte, 0)
	}
	s.sendProtoHeader.ResetForSend()
	s.recvProtoHeader.ResetForRecv(recvProtoExtraDataType)

	s.chWaitRecvSendGoroutineExit = make(chan struct{}, 2)
	s.chNtfRecvSendGoroutineExit = make(chan struct{})
	s.onceClose.Reset()

	s.manualClose = false

	s.chNetEventData, s.chSendProto = chNetEventData, chSendProto

	log.RunLogger.Printf("tcpSession.Start: %v", s)

	s.chNetEventData <- newSessionNetEventOn(s)

	// start recv goroutine
	go util.SafeGo(s.mainRecv, s.mainRecvEnd, s.uuid)

	// start send goroutine
	go util.SafeGo(s.mainSend, s.mainSendEnd, s.uuid)
}

// Close 在执行Start之前, 就直接关闭连接, 用于外界已决定关闭服务时新建立的连接需要立即关闭
func (s *tcpSession) Close() {
	s.onceClose.Do(func() {
		// 关闭底层连接
		s.conn.Close()

		// 归还
		sessPool.back(s)
	})
}

func (s *tcpSession) String() string {
	return fmt.Sprintf(`%p uuidSession[%v]`, s, s.uuid)
}

// setConn 设置底层连接对象
func (s *tcpSession) setConn(conn net.Conn) {
	s.conn = conn
}

func (s *tcpSession) mainSend(params ...interface{}) {
	log.RunLogger.Printf("tcpSession.mainSend: %v", s)

	for {
		select {
		case <-s.chNtfRecvSendGoroutineExit:
			// 关闭
			return
		case proto := <-s.chSendProto:
			// 发送
			if proto == nil {
				// 标识主动退出
				s.manualClose = true
				log.RunLogger.Printf("tcpSession.mainSend start quit: %v", s)
				return
			} else if !s.doSend(proto) {
				return
			}
		}

		// 有未发送完毕的数据, 且当前没有等待发送的协议
	loopSendLeft:
		for len(s.sendLeft) > 0 && len(s.chSendProto) == 0 {
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
				break loopSendLeft
			}
		}
	}
}

func (s *tcpSession) mainSendEnd(isPanic bool) {
	log.RunLogger.Printf("tcpSession.mainSendEnd isPanic[%v]: %v", isPanic, s)

	s.chWaitRecvSendGoroutineExit <- struct{}{}

	s.onceClose.Do(func() {
		s.doClose()
	})
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
	log.RunLogger.Printf("tcpSession.mainRecv: %v", s)

	var err error
	for {
		// 接收
		if err = s.doRecv(); err != nil {
			log.RunLogger.Printf("tcpSession.mainRecv err[%v] manualClose[%v]: %v", err, s.manualClose, s)
			return
		}

		// 是否关闭
		select {
		case <-s.chNtfRecvSendGoroutineExit:
			return
		default:
			break
		}
	}
}
func (s *tcpSession) mainRecvEnd(isPanic bool) {
	log.RunLogger.Printf("tcpSession.mainRecvEnd isPanic[%v]: %v", isPanic, s)

	s.chWaitRecvSendGoroutineExit <- struct{}{}

	s.onceClose.Do(func() {
		s.doClose()
	})
}

func (s *tcpSession) doRecv() (err error) {
	// 接收协议头
	err = s.doRecvHeader()
	if err != nil {
		return
	}

	// 接收协议内容
	data, err := s.doRecvContent()
	if err != nil {
		return
	}

	// 协议事件
	s.chNetEventData <- data

	return nil
}

// doRecvHeader 接收协议头
func (s *tcpSession) doRecvHeader() error {
	// solve dead link problem:
	// physical disconnection without any communcation between client and server
	// will cause the read to block FOREVER, so a timeout is a rescue.
	s.conn.SetReadDeadline(time.Now().Add(time.Duration(readDeadtime) * time.Second))

	// read Proto header
	_, err := io.ReadFull(s.conn, s.recvHeaderBuf)
	if err != nil {
		log.RunLogger.Printf("tcpSession.doRecvHeader ReadFull header error[%v]: %v", err, s)
		return err
	}

	log.RunLogger.Printf("tcpSession.doRecvHeader recvHeaderBuf[%v]: %v", s.recvHeaderBuf, s)

	// 协议头解析
	err = s.recvProtoHeader.Unmarshal(s.recvHeaderBuf)
	if err != nil {
		log.RunLogger.Printf("tcpSession.doRecvHeader recvProtoHeader.Unmarshal error[%v]: %v", err, s)
		return err
	}

	return nil
}

// doRecvContent 接收协议内容
func (s *tcpSession) doRecvContent() (data base.NetEventData, err error) {
	// 接收剩余部分
	p := ffProto.ApplyProtoForRecv(s.recvProtoHeader)

	defer p.BackAfterRecv()

	buf := p.BytesForRecv()

	// solve dead link problem:
	// physical disconnection without any communcation between client and server
	// will cause the read to block FOREVER, so a timeout is a rescue.
	s.conn.SetReadDeadline(time.Now().Add(time.Duration(readDeadtime) * time.Second))

	_, err = io.ReadFull(s.conn, buf)
	if err != nil {
		log.RunLogger.Printf("tcpSession.doRecv ReadFull content error[%v]: %v", err, s)
		return
	}

	log.RunLogger.Printf("tcpSession.doRecv recvProtoData[%v]: %v", buf, s)

	// 数据接收完毕, 通知校验
	err = p.OnRecvAllBytes(s.recvProtoHeader)
	if err != nil {
		log.RunLogger.Printf("tcpSession.doRecv proto[%v] OnRecvAllBytes error[%v]: %v", p, err, s)
		return
	}

	// 设置Proto状态为等待分发
	p.SetCacheWaitDispatch()

	return newSessionNetEventProto(s, p), nil
}

// doClose Session本次有效期间, 只会被执行一次
func (s *tcpSession) doClose() {
	log.RunLogger.Printf("tcpSession.doClose: %v", s)

	// 关闭结束管道, 触发发送/接收协程退出
	close(s.chNtfRecvSendGoroutineExit)

	// 关闭底层连接
	s.conn.Close()

	// 等待发送和接收协程退出
	<-s.chWaitRecvSendGoroutineExit
	<-s.chWaitRecvSendGoroutineExit

	// 连接断开事件
	s.chNetEventData <- newSessionNetEventOff(s, s.manualClose)
}

// back 外界已停止引用Session, 可安全回收了
func (s *tcpSession) back() {
	log.RunLogger.Printf("tcpSession.back: %v", s)

	// 清理内部数据
	s.chSendProto = nil
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
