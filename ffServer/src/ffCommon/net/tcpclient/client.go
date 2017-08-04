package tcpclient

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

// tcpClient connect Server
type tcpClient struct {
	tcpAddr *net.TCPAddr // 地址信息

	uuid uuid.UUID // 唯一标识

	chNewSession   chan base.Session // 外界接收新连接被创建事件的管道
	chClientClosed chan struct{}     // 完成关闭时, 向外界通知

	chNtfWorkExit chan struct{} // 退出
	chReConnect   chan struct{} // 重连

	onceClose util.Once // 用于只执行一次关闭
}

// Start 开始连接Server, 只执行一次, 异步
func (c *tcpClient) Start(chNewSession chan base.Session, chClientClosed chan struct{}) {
	c.chNewSession, c.chClientClosed = chNewSession, chClientClosed

	c.chNtfWorkExit = make(chan struct{})
	c.chReConnect = make(chan struct{}, 1)

	log.RunLogger.Printf("tcpClient.Start: %v", c)

	go util.SafeGo(c.mainLoop, c.mainLoopEnd)
}

// Stop 停止Client
func (c *tcpClient) Stop() {
	log.RunLogger.Printf("tcpClient.Stop: %v", c)

	c.onceClose.Do(func() {
		close(c.chNtfWorkExit)
		c.chNtfWorkExit = nil
	})
}

// Back 回收Client资源, 只应在外界通过chServerClose接收到可回收事件之后下执行
func (c *tcpClient) Back() {
	log.RunLogger.Printf("tcpClient.Back: %v", c)

	close(c.chReConnect)
	c.chReConnect = nil

	mutexClient.Lock()
	defer mutexClient.Unlock()
	delete(mapClients, c.uuid)
}

// UUID 唯一标识
func (c *tcpClient) UUID() uuid.UUID {
	return c.uuid
}

// String 返回Client的自我描述
func (c *tcpClient) String() string {
	return fmt.Sprintf(`%p uuidClient[%v]`, c, c.uuid)
}

func (c *tcpClient) mainLoop(params ...interface{}) {
	for {
		conn, err := net.DialTCP("tcp", nil, c.tcpAddr)
		if err == nil {
			log.RunLogger.Printf("tcpClient.mainLoop connect success: %v", c)

			c.chNewSession <- tcpsession.Apply(conn)

			// 等待退出或重连逻辑
			{
				select {
				case <-c.chNtfWorkExit: // 等待退出通知
					// 退出
					return

				case <-c.chReConnect: // 等待上一连接关闭且回收完毕
					break
				}
			}

		} else {
			log.RunLogger.Printf("tcpClient.mainLoop err[%v]: %v", err, c)

			// 连接失败, 自动重连
			<-time.After(time.Second)

			// 检查退出逻辑
			{
				select {
				case <-c.chNtfWorkExit:
					// 退出
					return

				default:
					break
				}
			}
		}
	}
}
func (c *tcpClient) mainLoopEnd() {
	log.RunLogger.Printf("tcpClient.mainLoopEnd: %v", c)

	c.chClientClosed <- struct{}{}
}

// newClient 新建一个 tcpClient
func newClient(addr string, uuid uuid.UUID) (s *tcpClient, err error) {
	log.RunLogger.Printf("tcpclient.newClient: addr[%v] uuid[%v]", addr, uuid)

	// 监听地址有效性
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcpClient.newClient ResolveTCPAddr failed, uuid[%d] addr[%v] err[%v]",
			s.uuid, addr, err)
	}

	client := &tcpClient{
		tcpAddr: tcpAddr,

		uuid: uuid,
	}

	mapClients[uuid] = client

	return client, nil
}
