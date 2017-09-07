package tcpserver

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpsession"
	"ffCommon/util"
	"ffCommon/uuid"

	"fmt"
	"net"
	"time"
)

// tcpServer 监听客户端连接
type tcpServer struct {
	tcpAddr  *net.TCPAddr // 地址信息
	listener net.Listener // listener

	chNewSession   chan base.Session // 外界接收新连接被创建事件的管道
	chServerClosed chan struct{}     // 完成关闭时, 向外界通知

	uuid uuid.UUID // 唯一标识

	working bool // 是否正常工作状态

	onceClose util.Once // 用于只执行一次关闭
}

// Start 开始监听, 等候客户端连接, 只应调用一次
func (s *tcpServer) Start(chNewSession chan base.Session, chServerClosed chan struct{}) (err error) {
	// 建立监听
	listener, err := net.ListenTCP("tcp", s.tcpAddr)
	if err != nil {
		return fmt.Errorf("tcpServer[%v].Start ListenTCP failed err[%v]", s, err)
	}

	s.listener = listener

	s.chNewSession, s.chServerClosed = chNewSession, chServerClosed

	s.onceClose.Reset()

	s.working = true

	log.RunLogger.Printf("tcpServer[%v].Start", s)

	// tcpServer loop
	go util.SafeGo(s.mainAccept, s.mainAcceptEnd)

	return
}

// StopAccept 关闭服务器, 只应在关闭服务器进程时调用
func (s *tcpServer) StopAccept() {
	log.RunLogger.Printf("tcpServer[%v].Close", s)

	// 立即标识停止工作
	s.working = false

	// 停止监听, 将导致mainAccept开始退出
	s.onceClose.Do(func() {
		if s.listener != nil {
			s.listener.Close()
			s.listener = nil
		}
	})
}

// Back 回收Server资源, 只应在Start失败或者外界通过chServerClose接收到可回收事件之后下执行
func (s *tcpServer) Back() {
	log.RunLogger.Printf("tcpServer[%v].mainSession Back", s)

	// 不再引用外界管道
	s.chNewSession, s.chServerClosed = nil, nil

	// 数据清理
	s.tcpAddr = nil

	// 移除记录
	mutexServer.Lock()
	defer mutexServer.Unlock()
	delete(mapServers, s.uuid)
}

// UUID 唯一标识
func (s *tcpServer) UUID() uuid.UUID {
	return s.uuid
}

// String 返回Server的自我描述
func (s *tcpServer) String() string {
	return fmt.Sprintf(`%p:%v`, s, s.uuid)
}

// mainAccept 接受客户端连接请求
func (s *tcpServer) mainAccept(params ...interface{}) {
	log.RunLogger.Printf("tcpServer[%v].mainAccept", s)

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

			log.RunLogger.Printf("tcpServer[%v].mainAccept Accept error[%v], retry in [%v] milisecond", s, err, tempDelay)

			time.Sleep(tempDelay)

			continue
			// }

			// 其他错误
			// return
		}
		tempDelay = 0

		// 创建session
		sess := tcpsession.Apply(conn)

		log.RunLogger.Printf("tcpServer[%v].mainAccept accept session[%v]", s, sess)

		// 向外界通知
		s.chNewSession <- sess
	}
}

// mainAcceptEnd 接受客户端连接彻底退出了
func (s *tcpServer) mainAcceptEnd(isPanic bool) {
	log.RunLogger.Printf("tcpServer[%v].mainAcceptEnd isPanic[%v]", s, isPanic)

	// 退出完成
	s.chServerClosed <- struct{}{}
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
