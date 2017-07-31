package tcpserver

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpsession"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"

	"fmt"
	"net"
	"sync"
	"time"
)

// tcpServer 监听客户端连接, 并维持与客户端之间的连接
type tcpServer struct {
	tcpAddr                *net.TCPAddr          // 地址信息
	listener               net.Listener          // listener
	recvProtoExtraDataType ffProto.ExtraDataType // 此服务器接收到的协议, 附加数据类型限定

	uuid uuid.UUID // 唯一标识

	chNetEventDataOuter chan base.NetEventData // 外界接收事件数据的管道
	chNetEventDataInner chan base.NetEventData // 自身接收Session事件数据的管道

	working bool // 是否正常工作状态

	chWaitWorkExit chan struct{} // 等待协程退出

	chNewSession       chan base.Session // 接受连接请求的协程 ==> session处理协程: 新连接
	chAllSessionClosed chan struct{}     // 所有连接都已断开

	mutexSession sync.RWMutex               // 连接管理锁
	mapSessions  map[uuid.UUID]base.Session // 所有的连接

	onceClose util.Once // 用于只执行一次关闭
}

// Start 开始监听, 等候客户端连接, 只应调用一次
func (s *tcpServer) Start(chNetEventData chan base.NetEventData, recvProtoExtraDataType ffProto.ExtraDataType) (err error) {
	log.RunLogger.Printf("tcpServer.Start: %v", s)

	// 建立监听
	listener, err := net.ListenTCP("tcp", s.tcpAddr)
	if err != nil {
		return fmt.Errorf("tcpServer.Start ListenTCP failed err[%v]: %v", err, s)
	}

	s.listener = listener
	s.recvProtoExtraDataType = recvProtoExtraDataType

	s.chWaitWorkExit = make(chan struct{}, 1)

	s.chNewSession = make(chan base.Session, 4)
	s.chAllSessionClosed = make(chan struct{})

	s.working = true

	// tcpServer loop
	go util.SafeGo(s.mainAccept)
	go util.SafeGo(s.mainSession)

	return
}

// SendProto 发送Proto到指定对端, 异步
func (s *tcpServer) SendProto(uuidSession uuid.UUID, proto *ffProto.Proto) (err error) {
	log.RunLogger.Printf("tcpServer.SendProto uuidSession[%v] proto[%v]: %v", uuidSession, proto, s)

	s.mutexSession.RLock()
	defer s.mutexSession.RUnlock()
	if sess, ok := s.mapSessions[uuidSession]; ok {
		sess.SendProto(proto)
		return
	}

	return fmt.Errorf("tcpServer.SendProto invalid uuidSession[%v]: %v", uuidSession, s)
}

// CloseSession 断开指定连接
func (s *tcpServer) CloseSession(uuidSession uuid.UUID, delayMillisecond int64) {
	log.RunLogger.Printf("tcpServer.CloseSession uuidSession[%v] delayMillisecond[%d]: %v", uuidSession, delayMillisecond, s)

	s.mutexSession.RLock()
	defer s.mutexSession.RUnlock()

	if sess, ok := s.mapSessions[uuidSession]; ok {
		sess.Close(delayMillisecond)
	} else {
		log.FatalLogger.Printf("tcpServer.CloseSession uuidSession[%v] not exist: %s", uuidSession, s)
	}
}

// Close tcpServer, 只应在关闭服务器进程时调用
func (s *tcpServer) Close(delayMillisecond int64) {
	log.RunLogger.Printf("tcpServer.Close delayMillisecond[%d]: %v", delayMillisecond, s)

	// 立即标识停止工作
	s.working = false

	go util.SafeGo(func(params ...interface{}) {
		if delayMillisecond > 0 {
			select {
			case <-time.After(time.Duration(delayMillisecond) * time.Millisecond):
				s.onceClose.Do(func() {
					s.doClose()
				})
			}
		} else {
			s.onceClose.Do(func() {
				s.doClose()
			})
		}
	})
}

// doClose 关闭服务器
func (s *tcpServer) doClose() {
	log.RunLogger.Printf("tcpServer.doClose: %v", s)

	s.listener.Close()
	s.listener = nil

	// 等待接受客户端连接的协程退出
	<-s.chWaitWorkExit

	// 等待session处理协程退出
	<-s.chWaitWorkExit

	// 结束了
	s.chNetEventDataOuter <- newClientNetEventDataEnd(s)
}

// back 外界处理完毕Server关闭事件
func (s *tcpServer) back() {
	log.RunLogger.Printf("tcpServer.mainSession back: %v", s)

	close(s.chNetEventDataInner)
	s.chNetEventDataInner = nil

	close(s.chNewSession)
	s.chNewSession = nil

	close(s.chWaitWorkExit)
	s.chWaitWorkExit = nil

	s.chAllSessionClosed = nil

	mutexServer.Lock()
	defer mutexServer.Unlock()
	delete(mapServers, s.uuid)
}

// 接受客户端连接请求
func (s *tcpServer) mainAccept(params ...interface{}) {
	// 协程退出时记录
	defer func() {
		log.RunLogger.Printf("tcpServer.mainAccept end: %v", s)

		if err := recover(); err != nil {
			util.PrintPanicStack(err, "tcpServer.mainAccept", s)
		}

		// 完成了退出
		s.chWaitWorkExit <- struct{}{}
	}()

	var tempDelay time.Duration
	for {
		conn, err := s.listener.Accept()

		// 出错时, 只要不是主动关闭, 都将继续监听
		if err != nil {
			if !s.working {
				return
			}

			// 临时错误
			// 最长等待1秒, 然后再次尝试接受连接请求
			// if ne, ok := err.(net.Error); ok && ne.Temporary() {
			if tempDelay == 0 {
				tempDelay = 50 * time.Millisecond
			} else {
				tempDelay *= 2
			}
			if max := 1 * time.Second; tempDelay > max {
				tempDelay = max
			}

			log.RunLogger.Printf("tcpServer.mainAccept Accept error[%v], retry in [%v] milisecond: %v", err, tempDelay, s)

			time.Sleep(tempDelay)

			continue
			// }

			// 其他错误
			// return
		}
		tempDelay = 0

		// 新连接激活
		sess := tcpsession.Apply()
		sess.Start(conn, s.chNetEventDataInner, s.recvProtoExtraDataType)

		s.chNewSession <- sess
	}
}

// 处理连接事宜
func (s *tcpServer) mainSession(params ...interface{}) {
	// 协程退出时记录
	defer func() {
		log.RunLogger.Printf("tcpServer.mainSession end: %v", s)

		if err := recover(); err != nil {
			util.PrintPanicStack(err, "tcpServer.mainSession", s)
		}

		// 完成了退出
		s.chWaitWorkExit <- struct{}{}
	}()

	// 主循环
	{
	mainLoop:
		for {
			select {
			case sess := <-s.chNewSession:
				if sess != nil {
					// 新连接建立
					s.onSessionCreated(sess)
				} else {
					// 停止接受连接请求了
					break mainLoop
				}
			case dataSession := <-s.chNetEventDataInner: // 处理session事件
				s.onSessionEvent(dataSession)
			}
		}
	}

	// 通知: 关闭所有Session
	s.closeAllSession()

	// 收尾: 处理完毕所有事件
	{
	endLoop:
		for {
			select {
			case dataSession := <-s.chNetEventDataInner: // 继续处理session事件
				s.onSessionEvent(dataSession)
			case <-s.chAllSessionClosed: // 收尾工作完成
				break endLoop
			}
		}
	}
}

// closeAllSession 关闭所有连接
func (s *tcpServer) closeAllSession() {
	log.RunLogger.Printf("tcpServer.closeAllSession: %v", s)

	s.mutexSession.RLock()
	defer s.mutexSession.RUnlock()

	for _, sess := range s.mapSessions {
		sess.Close(0)
	}

	if len(mapServers) == 0 {
		close(s.chAllSessionClosed)
	}
}

// onSessionCreated 新连接建立
func (s *tcpServer) onSessionCreated(sess base.Session) {
	log.RunLogger.Printf("tcpServer.onSessionCreated uuid[%v]: %v", sess.UUID(), s)

	s.mutexSession.Lock()
	defer s.mutexSession.Unlock()
	s.mapSessions[sess.UUID()] = sess
}

// onSessionClosed 外界处理完毕session断开事件
func (s *tcpServer) onSessionClosed(uuidSession uuid.UUID) {
	log.RunLogger.Printf("tcpServer.onSessionClosed uuid[%v]: %v", uuidSession, s)

	s.mutexSession.Lock()
	defer s.mutexSession.Unlock()
	delete(s.mapSessions, uuidSession)
	if !s.working && len(s.mapSessions) == 0 {
		close(s.chAllSessionClosed)
	}
}

// onSessionEvent 处理session事件
func (s *tcpServer) onSessionEvent(dataSession base.NetEventData) {
	eventType := dataSession.NetEventType()
	if eventType == base.NetEventEnd {
		// 不向外界抛出此事件, 直接回收
		dataSession.Back()
		return
	}

	data, _ := dataSession.(base.SessionNetEventData)
	s.chNetEventDataOuter <- newServerNetEventDataFromSessionNetEventData(s, data)
}

// String 返回Server的自我描述
func (s *tcpServer) String() string {
	return fmt.Sprintf(`uuid[%v] working[%v] tcpAddr[%v] recvProtoExtraDataType[%v]`,
		s.uuid, s.working, s.tcpAddr, s.recvProtoExtraDataType)
}

// newServer 新建一个 tcpServer
func newServer(addr string, uuid uuid.UUID) (s *tcpServer, err error) {
	log.RunLogger.Printf("tcpserver.newServer: addr[%v] uuid[%v]", addr, uuid)

	// 监听地址有效性
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcpServer.newServer ResolveTCPAddr failed, addr[%v] err[%v]",
			addr, err)
	}

	return &tcpServer{
		tcpAddr: tcpAddr,

		uuid: uuid,
	}, nil
}
