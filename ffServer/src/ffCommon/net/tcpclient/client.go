package tcpclient

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

// Client connect Server
type Client struct {
	tcpAddr       *net.TCPAddr // 地址信息
	autoReconnect bool         // 是否自动重连

	agent                  base.Agent            // Client's agent
	recvProtoExtraDataType ffProto.ExtraDataType // 此客户端接收到的协议, 附加数据类型限定

	sess base.Session // 当前与服务器的连接

	chClose chan struct{}  // 关闭Client
	wgClose sync.WaitGroup // 等待关闭Client完成
}

// Close Client
func (c *Client) Close() {
	// 不再重连
	c.autoReconnect = false

	close(c.chClose)
	c.wgClose.Wait()

	c.chClose = nil
	c.sess = nil
	c.agent = nil
}

// SendProto SendProto
func (c *Client) SendProto(p *ffProto.Proto) {
	c.sess.SendProto(p)
}

// Start Client
func (c *Client) Start(agent base.Agent, recvProtoExtraDataType ffProto.ExtraDataType) error {
	c.agent = agent
	c.recvProtoExtraDataType = recvProtoExtraDataType
	go util.SafeGo(c.doConnnect)
	return nil
}

func (c *Client) doConnnect(params ...interface{}) {
	c.wgClose.Add(1)
	defer c.wgClose.Done()

	for c.autoReconnect {
		conn, err := net.DialTCP("tcp", nil, c.tcpAddr)
		if err == nil {
			c.onConnect(conn)
		} else {
			log.RunLogger.Println(err)
		}

		// 自动重连
		if c.autoReconnect {
			<-time.After(time.Second)
		}
	}
}

func (c *Client) onConnect(conn net.Conn) {
	// 启用Session
	c.sess = session.Apply()
	c.sess.Start(conn, c, c.recvProtoExtraDataType)

	// 等待连接关闭
	select {
	case <-c.sess.WaitCloseChan():
	case <-c.chClose:
		c.sess.Close(0)
	}
}

// OnEvent base.Agent'OnEvent
func (c *Client) OnEvent(protoID ffProto.MessageType, data interface{}) {
	// 上层逻辑必须执行Done
	if protoID == ffProto.MessageType_SessionDisConnect {
		wg, _ := data.(*sync.WaitGroup)
		wg.Add(1)
	}

	c.agent.OnEvent(protoID, data)

	if protoID == ffProto.MessageType_SessionDisConnect {
		c.sess = nil

		wg, _ := data.(*sync.WaitGroup)
		wg.Done()
	}
}

// String 返回Client的自我描述
func (c *Client) String() string {
	return fmt.Sprintf(`tcpAddr[%v] autoReconnect[%v] recvProtoExtraDataType[%v]`,
		c.tcpAddr, c.autoReconnect, c.recvProtoExtraDataType)
}

// New create new Client
func New(addr string, autoReconnect bool) (c base.Client, err error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("tcpclient.New: ResolveTCPAddr failed. addr[%v] err[%v]", addr, err)
	}

	return &Client{
		tcpAddr:       tcpAddr,
		autoReconnect: autoReconnect,

		chClose: make(chan struct{}),
	}, nil
}
