package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpserver"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"sync/atomic"
)

// 用户侧管理
type userAgentServer struct {
	config            *serverUserConfig     // 配置
	sendExtraDataType ffProto.ExtraDataType // 发送的协议的附加数据类型
	recvExtraDataType ffProto.ExtraDataType // 接收的协议的附加数据类型

	server base.Server // server 底层server

	chNewSession   chan base.Session // chNewSession 新连接
	chServerClosed chan struct{}     // 用于接收服务器退出事件
	chAgentClosed  chan *userAgent   // 用于接收userAgent关闭事件

	mapUserAgent map[uuid.UUID]*userAgent // 所有连接用户
	agentPool    *userAgentPool           // 所有用户缓存
}

func (agentServer *userAgentServer) Status() string {
	return fmt.Sprintf("chNewSession[%v] chAgentClosed[%v] mapUserAgent[%v] agentPool[%v]",
		len(agentServer.chNewSession), len(agentServer.chAgentClosed), len(agentServer.mapUserAgent), agentServer.agentPool)
}

// onNewSession 新连接
func (agentServer *userAgentServer) onNewSession(sess base.Session) {
	log.RunLogger.Printf("userAgentServer.onNewSession sess[%v]", sess)

	agent := agentServer.agentPool.apply()
	agentServer.mapUserAgent[sess.UUID()] = agent
	agent.Start(sess, agentServer)
}

// onAgentClosed userAgent关闭
func (agentServer *userAgentServer) onAgentClosed(agent *userAgent) {
	log.RunLogger.Printf("userAgentServer.onAgentClosed %v", agent)

	delete(agentServer.mapUserAgent, agent.uuidSessiont)

	// 回收清理
	agent.Back()

	// 缓存
	agentServer.agentPool.back(agent)
}

// mainLoop
func (agentServer *userAgentServer) mainLoop(params ...interface{}) {
	atomic.AddInt32(&waitServerQuit, 1)

	// 主循环
	{
	mainLoop:
		for {
			select {
			case sess := <-agentServer.chNewSession: // 新连接
				agentServer.onNewSession(sess)

			case agent := <-agentServer.chAgentClosed: // 连接结束
				agentServer.onAgentClosed(agent)

			case <-chApplicationQuit: // 进程退出
				agentServer.server.StopAccept()
				break mainLoop
			}
		}
	}

	// 等待底层服务器退出完成
	{
		<-agentServer.chServerClosed
	}

	// 继续处理新连接(直接关闭)
	{
	endNewSession:
		for {
			select {
			case sess := <-agentServer.chNewSession: // 新连接
				// todo: 此分支需要测试
				// 直接关闭
				sess.Close()
			default:
				break endNewSession
			}
		}
	}

	// 关闭所有已建立的连接
	{
		if len(agentServer.mapUserAgent) > 0 {
			// 向其通知退出
			for _, agent := range agentServer.mapUserAgent {
				agent.Close()
			}

			// 等待全部退出
		endSession:
			for {
				select {
				case agent := <-agentServer.chAgentClosed: // 连接结束
					agentServer.onAgentClosed(agent)

					// 全关闭了
					if len(agentServer.mapUserAgent) == 0 {
						break endSession
					}
				}
			}
		}
	}
}

// mainLoopEnd
func (agentServer *userAgentServer) mainLoopEnd() {
	log.RunLogger.Printf("userAgentServer.mainLoopEnd")

	atomic.AddInt32(&waitServerQuit, -1)
}

// init 根据配置初始化Server
func (agentServer *userAgentServer) start(config *serverUserConfig) (err error) {
	agentServer.sendExtraDataType, err = ffProto.GetExtraDataType(config.SendExtraDataType)
	if err != nil {
		return err
	}

	agentServer.recvExtraDataType, err = ffProto.GetExtraDataType(config.RecvExtraDataType)
	if err != nil {
		return err
	}

	server, err := tcpserver.NewServer(config.ListenAddr)
	if err != nil {
		return err
	}

	chNewSession := make(chan base.Session, config.AcceptNewSessionCache)
	chServerClosed := make(chan struct{}, 1)

	// 开启服务器
	if err = server.Start(chNewSession, chServerClosed); err != nil {
		close(chNewSession)
		close(chServerClosed)

		// 开启失败, 回收
		server.Back()
		return err
	}

	agentServer.config = config

	agentServer.server = server
	agentServer.chNewSession, agentServer.chServerClosed = chNewSession, chServerClosed

	agentServer.mapUserAgent = make(map[uuid.UUID]*userAgent, config.InitOnlineCount)
	agentServer.agentPool = newUserAgentPool(config.InitOnlineCount)

	return nil
}
