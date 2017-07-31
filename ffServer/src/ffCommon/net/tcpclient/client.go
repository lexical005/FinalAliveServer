package tcpclient

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpsession"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"

	"fmt"
	"net"
	"time"
)

// tcpClient connect Server
type tcpClient struct {
	tcpAddr *net.TCPAddr // 地址信息

	uuid uuid.UUID // 唯一标识

	recvProtoExtraDataType ffProto.ExtraDataType // 此客户端接收到的协议, 附加数据类型限定

	sess base.Session // 当前与服务器的连接

	chNetEventDataOuter chan base.NetEventData // 外界接收事件数据的管道
	chNetEventDataInner chan base.NetEventData // 自身接收Session事件数据的管道

	working bool // 是否正常工作状态

	chNtfWorkExit chan struct{} // 退出

	chReConnect chan struct{} // 重连

	onceClose util.Once // 用于只执行一次关闭
}

// Close tcpClient
func (c *tcpClient) Close(delayMillisecond int64) {
	log.RunLogger.Printf("tcpClient.Close delayMillisecond[%d]: %v", delayMillisecond, c)

	// 立即标识停止工作
	c.working = false

	go util.SafeGo(func(params ...interface{}) {
		if delayMillisecond > 0 {
			select {
			case <-time.After(time.Duration(delayMillisecond) * time.Millisecond):
				c.onceClose.Do(func() {
					c.doClose()
				})
			}
		} else {
			c.onceClose.Do(func() {
				c.doClose()
			})
		}
	})
}

// SendProto 发送Proto到对端, 只应该在收到连接建立事件之后再调用, 异步
func (c *tcpClient) SendProto(p *ffProto.Proto) {
	c.sess.SendProto(p)
}

// Start 开始连接Server, 只执行一次, 异步
func (c *tcpClient) Start(chNetEventData chan base.NetEventData, recvProtoExtraDataType ffProto.ExtraDataType) error {
	c.chNetEventDataOuter, c.recvProtoExtraDataType = chNetEventData, recvProtoExtraDataType

	c.chNetEventDataInner = make(chan base.NetEventData, DefaultClientNetEventDataChanCount)
	c.chNtfWorkExit = make(chan struct{})
	c.chReConnect = make(chan struct{}, 2)

	c.working = true

	log.RunLogger.Printf("tcpClient.Start: %v", c)

	go util.SafeGo(c.mainLoop)

	return nil
}

// ReConnect 重连, 只能在外界处理到NetEventOff事件时, 如果需要恢复连接, 才可调用此接口, 一次NetEventOff事件对应一次ReConnect调用, 异步
func (c *tcpClient) ReConnect() {
	log.RunLogger.Printf("tcpClient.ReConnect: %v", c)

	c.chReConnect <- struct{}{}
}

// String 返回Client的自我描述
func (c *tcpClient) String() string {
	return fmt.Sprintf(`uuid[%v] tcpAddr[%v] recvProtoExtraDataType[%v]`,
		c.uuid, c.tcpAddr, c.recvProtoExtraDataType)
}

func (c *tcpClient) mainLoop(params ...interface{}) {
	// 协程退出时记录
	defer func() {
		log.RunLogger.Printf("tcpClient.mainLoop end: %v", c)

		if err := recover(); err != nil {
			util.PrintPanicStack(err, "tcpClient.mainLoop", c)
		}
	}()

	for {
		conn, err := net.DialTCP("tcp", nil, c.tcpAddr)
		if err == nil {
			log.RunLogger.Printf("tcpClient.mainLoop success: %v", c)

			// 本次连接主循环
			{
				// 创建Session
				c.sess = tcpsession.Apply()
				c.sess.Start(conn, c.chNetEventDataInner, c.recvProtoExtraDataType)

				sessionOn := true

			sessionLoop:
				for {
					select {
					case dataSession := <-c.chNetEventDataInner: // 转发事件
						eventType := dataSession.NetEventType()
						if eventType == base.NetEventOff {
							sessionOn = false
						} else if eventType == base.NetEventEnd {
							// 不向外界抛出此事件, 直接回收
							dataSession.Back()

							// 结束在此session上的事件接收
							break sessionLoop
						}

						c.chNetEventDataOuter <- newClientNetEventDataFromSessionNetEventData(c, dataSession)

					case <-c.chNtfWorkExit: // 等待退出通知
						if sessionOn {
							// 关闭连接
							c.sess.Close(0)

							// 继续分发事件, 直至处理到NetEventEnd
							for {
								select {
								case dataSession := <-c.chNetEventDataInner: // 转发事件
									eventType := dataSession.NetEventType()
									if eventType == base.NetEventEnd {
										// 不向外界抛出此事件, 直接回收
										dataSession.Back()

										// 退出
										return
									}

									c.chNetEventDataOuter <- newClientNetEventDataFromSessionNetEventData(c, dataSession)
								}
							}
						}

						// 退出
						return
					}
				}
			}

			// 退出和重连逻辑
			{
				select {
				case <-c.chNtfWorkExit: // 等待退出通知
					// 退出
					return

				case <-c.chReConnect: // 等待上一连接回收完毕
					break
				}
			}

		} else {
			log.RunLogger.Printf("tcpClient.mainLoop err[%v]: %v", err, c)

			// 连接失败, 自动重连
			<-time.After(time.Second)

			// 退出逻辑
			{
				select {
				case <-c.chNtfWorkExit:
					// 结束事件
					c.chNetEventDataOuter <- newClientNetEventDataEnd(c)
					// 退出
					return

				default:
					break
				}
			}
		}
	}
}

// doClose 执行外界的要求: 关闭client
func (c *tcpClient) doClose() {
	log.RunLogger.Printf("tcpClient.doClose: %v", c)

	// 通知退出
	close(c.chNtfWorkExit)
}

// onSessionClosed 外界处理完毕连接断开事件
func (c *tcpClient) onSessionClosed() {
	log.RunLogger.Printf("tcpClient.onSessionClosed: %v", c)

	// 当前连接结束
	c.sess = nil

	// 允许重连
	c.chReConnect <- struct{}{}
}

// back 回收
func (c *tcpClient) back() {
	log.RunLogger.Printf("tcpClient.back: %v", c)

	close(c.chNetEventDataInner)
	c.chNetEventDataInner = nil

	c.chNtfWorkExit = nil

	close(c.chReConnect)
	c.chReConnect = nil

	mutexClient.Lock()
	defer mutexClient.Unlock()
	delete(mapClients, c.uuid)
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
