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

// AgentUser 管理器
type agentUserManager struct {
	config            *serveConfig          // 配置
	sendExtraDataType ffProto.ExtraDataType // 发送的协议的附加数据类型
	recvExtraDataType ffProto.ExtraDataType // 接收的协议的附加数据类型

	server     base.Server // server 底层server
	uuidServer uuid.UUID   // uuidServer server的UUID

	chNewSession   chan base.Session // 用于接收新连接事件
	chServerClosed chan struct{}     // 用于接收服务器退出事件

	chAgentClosed chan *agentUser          // 用于接收Agent关闭事件
	mapUserAgent  map[uuid.UUID]*agentUser // 所有连接用户
	agentPool     *agentUserPool           // 所有用户缓存
}

func (mgr *agentUserManager) Status() string {
	return fmt.Sprintf("uuid[%v] chNewSession[%v] chAgentClosed[%v] mapUserAgent[%v] agentPool[%v]",
		mgr.uuidServer, len(mgr.chNewSession), len(mgr.chAgentClosed), len(mgr.mapUserAgent), mgr.agentPool)
}

func (agentManager *agentUserManager) String() string {
	return fmt.Sprintf("uuid[%v]", agentManager.uuidServer)
}

func (agentManager *agentUserManager) doClear() {
	close(agentManager.chNewSession)
	agentManager.chNewSession = nil

	close(agentManager.chServerClosed)
	agentManager.chServerClosed = nil

	close(agentManager.chAgentClosed)
	agentManager.chAgentClosed = nil
}

func (agentManager *agentUserManager) onBaseServerClosed() {
	log.RunLogger.Printf("agentUserManager.onBaseServerClosed: %v", agentManager)

	agentManager.server.Back()
	agentManager.server = nil
}

// onNewSession 新连接
func (agentManager *agentUserManager) onNewSession(sess base.Session) {
	log.RunLogger.Printf("agentUserManager.onNewSession sess[%v]: %v", sess, agentManager)

	agent := agentManager.agentPool.apply()
	agentManager.mapUserAgent[sess.UUID()] = agent
	agent.Start(sess, agentManager)
}

// onAgentClosed Agent关闭
func (agentManager *agentUserManager) onAgentClosed(agent *agentUser) {
	log.RunLogger.Printf("agentUserManager.onAgentClosed %v: %v", agent, agentManager)

	delete(agentManager.mapUserAgent, agent.uuidSession)

	// 回收清理
	agent.Back()

	// 缓存
	agentManager.agentPool.back(agent)
}

// mainLoop
func (agentManager *agentUserManager) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("agentUserManager.mainLoop start: %v", agentManager)

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

	log.RunLogger.Printf("agentUserManager.mainLoop start application quit: %v", agentManager)

	// 等待底层服务器退出完成
	{
		<-agentManager.chServerClosed

		agentManager.onBaseServerClosed()
	}

	log.RunLogger.Printf("agentUserManager.mainLoop application quit step 1: recv base.Server closed: %v", agentManager)

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

	log.RunLogger.Printf("agentUserManager.mainLoop application quit step 2: close all new wait session: %v", agentManager)

	// 关闭所有已建立的连接
	{
		if len(agentManager.mapUserAgent) > 0 {
			// 向其通知退出
			for _, agent := range agentManager.mapUserAgent {
				agent.Close()
			}

			log.RunLogger.Printf("agentUserManager.mainLoop application quit step 3: notify user agent close: %v", agentManager)

			// 等待全部退出
		endSession:
			for {
				select {
				case agent := <-agentManager.chAgentClosed: // 连接结束
					agentManager.onAgentClosed(agent)

					// 全关闭了
					if len(agentManager.mapUserAgent) == 0 {
						break endSession
					}
				}
			}
		}

		log.RunLogger.Printf("agentUserManager.mainLoop application quit step 4: all user agent closed: %v", agentManager)
	}
}

// mainLoopEnd
func (agentManager *agentUserManager) mainLoopEnd() {
	log.RunLogger.Printf("agentUserManager.mainLoopEnd")

	atomic.AddInt32(&waitServerQuit, -1)
}

// init 根据配置初始化Server
func (agentManager *agentUserManager) start(config *serveConfig) (err error) {
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

	agentManager.chAgentClosed = make(chan *agentUser, 2)
	agentManager.mapUserAgent = make(map[uuid.UUID]*agentUser, config.InitOnlineCount)
	agentManager.agentPool = newAgentUserPool(config.InitOnlineCount)

	go util.SafeGo(agentManager.mainLoop, agentManager.mainLoopEnd)

	return nil
}
