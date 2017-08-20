package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/net/tcpserver"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"sync/atomic"
)

// AgentGameServer 管理器
type agentGameServerManager struct {
	config            *serveConfig          // 配置
	sendExtraDataType ffProto.ExtraDataType // 发送的协议的附加数据类型
	recvExtraDataType ffProto.ExtraDataType // 接收的协议的附加数据类型

	server     base.Server // server 底层server
	uuidServer uuid.UUID   // uuidServer server的UUID

	chNewSession   chan base.Session // 用于接收新连接事件
	chServerClosed chan struct{}     // 用于接收服务器退出事件

	chAgentClosed chan *agentGameServer          // 用于接收agent关闭事件
	mapAgents     map[uuid.UUID]*agentGameServer // 所有连接用户
	agentPool     *agentGameServerPool           // 所有用户缓存
}

func (agentManager *agentGameServerManager) Status() string {
	return fmt.Sprintf("chNewSession[%v] chAgentClosed[%v] mapUserAgent[%v] agentPool[%v]",
		len(agentManager.chNewSession), len(agentManager.chAgentClosed), len(agentManager.mapAgents), agentManager.agentPool)
}

func (agentManager *agentGameServerManager) String() string {
	return fmt.Sprintf("uuidUserAgentServer[%v]", agentManager.uuidServer)
}

func (agentManager *agentGameServerManager) doClear() {
	close(agentManager.chNewSession)
	agentManager.chNewSession = nil

	close(agentManager.chServerClosed)
	agentManager.chServerClosed = nil

	close(agentManager.chAgentClosed)
	agentManager.chAgentClosed = nil
}

func (agentManager *agentGameServerManager) onBaseServerClosed() {
	log.RunLogger.Printf("agentGameServerManager.onBaseServerClosed: %v", agentManager)

	agentManager.server.Back()
	agentManager.server = nil
}

// onNewSession 新连接
func (agentManager *agentGameServerManager) onNewSession(sess base.Session) {
	log.RunLogger.Printf("agentGameServerManager.onNewSession sess[%v]: %v", sess, agentManager)

	agent := agentManager.agentPool.apply()
	agentManager.mapAgents[sess.UUID()] = agent
	agent.Start(sess, agentManager)
}

// onAgentClosed agent关闭
func (agentManager *agentGameServerManager) onAgentClosed(agent *agentGameServer) {
	log.RunLogger.Printf("agentGameServerManager.onAgentClosed %v: %v", agent, agentManager)

	delete(agentManager.mapAgents, agent.uuid)

	// 回收清理
	agent.Back()

	// 缓存
	agentManager.agentPool.back(agent)
}

// mainLoop
func (agentManager *agentGameServerManager) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("agentGameServerManager.mainLoop start: %v", agentManager)

	atomic.AddInt32(&waitServerQuit, 1)

	// 主循环
	{
	mainLoop:
		for {
			select {
			case sess := <-agentManager.chNewSession: // 新连接
				agentManager.onNewSession(sess)

			case agent := <-agentManager.chAgentClosed: // 连接结束
				agentManager.onAgentClosed(agent)

			case <-chApplicationQuit: // 进程退出
				agentManager.server.StopAccept()
				break mainLoop
			}
		}
	}

	log.RunLogger.Printf("agentGameServerManager.mainLoop start application quit: %v", agentManager)

	// 等待底层服务器退出完成
	{
		<-agentManager.chServerClosed

		agentManager.onBaseServerClosed()
	}

	log.RunLogger.Printf("agentGameServerManager.mainLoop application quit step 1: recv base.Server closed: %v", agentManager)

	// 继续处理新连接(直接关闭)
	{
	endNewSession:
		for {
			select {
			case sess := <-agentManager.chNewSession: // 新连接
				// todo: 此分支需要测试
				// 直接关闭
				sess.Close()
			default:
				break endNewSession
			}
		}
	}

	log.RunLogger.Printf("agentGameServerManager.mainLoop application quit step 2: close all new wait session: %v", agentManager)

	// 关闭所有已建立的连接
	{
		if len(agentManager.mapAgents) > 0 {
			// 向其通知退出
			for _, agent := range agentManager.mapAgents {
				agent.Close()
			}

			log.RunLogger.Printf("agentGameServerManager.mainLoop application quit step 3: notify user agent close: %v", agentManager)

			// 等待全部退出
		endSession:
			for {
				select {
				case agent := <-agentManager.chAgentClosed: // 连接结束
					agentManager.onAgentClosed(agent)

					// 全关闭了
					if len(agentManager.mapAgents) == 0 {
						break endSession
					}
				}
			}
		}

		log.RunLogger.Printf("agentGameServerManager.mainLoop application quit step 4: all user agent closed: %v", agentManager)
	}
}

// mainLoopEnd
func (agentManager *agentGameServerManager) mainLoopEnd() {
	log.RunLogger.Printf("agentGameServerManager.mainLoopEnd")

	atomic.AddInt32(&waitServerQuit, -1)
}

// init 根据配置初始化Server
func (agentManager *agentGameServerManager) start(config *serveConfig) (err error) {
	agentManager.sendExtraDataType, err = ffProto.GetExtraDataType(config.SendExtraDataType)
	if err != nil {
		return err
	}

	agentManager.recvExtraDataType, err = ffProto.GetExtraDataType(config.RecvExtraDataType)
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

	agentManager.config = config

	agentManager.server, agentManager.uuidServer = server, server.UUID()
	agentManager.chNewSession, agentManager.chServerClosed = chNewSession, chServerClosed

	agentManager.chAgentClosed = make(chan *agentGameServer, 2)
	agentManager.mapAgents = make(map[uuid.UUID]*agentGameServer, config.InitOnlineCount)
	agentManager.agentPool = newAgentGameServerPool(config.InitOnlineCount)

	go util.SafeGo(agentManager.mainLoop, agentManager.mainLoopEnd)

	return nil
}
