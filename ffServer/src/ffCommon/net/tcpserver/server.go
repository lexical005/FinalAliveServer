package tcpserver

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/session"
	"ffCommon/util"
	"ffProto"

	"fmt"
	"net"
	"sync"
	"time"
)

// Server listen on ip:port to support clients to connect and commulicate
type Server struct {
	tcpAddr  *net.TCPAddr // 地址信息
	listener net.Listener // listener

	serverid int // 服务唯一标识

	wgClose     sync.WaitGroup // 用于关闭 Server
	manualClose bool           // 是不是手动关闭服务器

	agentCreator           base.AgentCreator     // agent 创建工厂， 外界提供
	recvProtoExtraDataType ffProto.ExtraDataType // 此服务器接收到的协议, 附加数据类型限定
}

// Start Server, 外界只应调用一次
func (s *Server) Start(agentCreator base.AgentCreator, recvProtoExtraDataType ffProto.ExtraDataType) (err error) {
	// 建立监听
	listener, err := net.ListenTCP("tcp", s.tcpAddr)
	if err != nil {
		return err
	}
	s.listener = listener
	s.agentCreator = agentCreator
	s.recvProtoExtraDataType = recvProtoExtraDataType

	// Server loop
	go util.SafeGo(s.mainLoop)

	return
}

// Close Server, 只应在关闭服务器时调用
func (s *Server) Close() {
	if !s.manualClose {
		s.manualClose = true
		s.listener.Close()

		s.wgClose.Wait()

		s.listener = nil
		s.agentCreator = nil

		// todo: 关闭现有连接, 属于逻辑曾事务!
	}
}

// 在独立的协程中运行Server的主逻辑: 接受连接请求, 创建/查找空闲Session, 并在独立的协程中运行Session的主逻辑
func (s *Server) mainLoop(params ...interface{}) {
	s.wgClose.Add(1)
	defer s.wgClose.Done()

	var tempDelay time.Duration
	for {
		conn, err := s.listener.Accept()

		// 出错时, 只要不是主动关闭, 都将继续监听
		if err != nil {
			// 服务器手动关闭了
			if s.manualClose {
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

			log.RunLogger.Printf("Server.mainLoop: Accept error[%v] retry in [%v] milisecond", err, tempDelay)

			time.Sleep(tempDelay)

			continue
			// }

			// 其他错误
			// return
		}
		tempDelay = 0

		// 新连接激活
		sess := session.Apply()
		agent := s.agentCreator.Create(sess)
		sess.Start(conn, agent, s.recvProtoExtraDataType)
	}
}

// String 返回Server的自我描述
func (s *Server) String() string {
	return fmt.Sprintf(`id[%v] tcpAddr[%v] manualClose[%v] recvProtoExtraDataType[%v]`,
		s.serverid, s.tcpAddr, s.manualClose, s.recvProtoExtraDataType)
}

// newServer 新建一个Server
func newServer(addr string, serverid int) (s *Server, err error) {
	// 监听地址有效性
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcpserver.newServer: ResolveTCPAddr failed. addr[%v] err[%v]", addr, err)
	}

	return &Server{
		tcpAddr: tcpAddr,

		serverid: serverid,
	}, nil
}

// NewServer 新建一个Server
// addr: 监听地址
func NewServer(addr string) (server base.Server, err error) {
	serverCount++
	return newServer(addr, serverCount)
}
