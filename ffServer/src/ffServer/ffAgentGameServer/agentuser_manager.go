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

func (mgr *agentUserManager) String() string {
	return fmt.Sprintf("uuid[%v]", mgr.uuidServer)
}

func (mgr *agentUserManager) doClear() {
	close(mgr.chNewSession)
	mgr.chNewSession = nil

	close(mgr.chServerClosed)
	mgr.chServerClosed = nil

	close(mgr.chAgentClosed)
	mgr.chAgentClosed = nil
}

func (mgr *agentUserManager) onBaseServerClosed() {
	log.RunLogger.Printf("agentUserManager.onBaseServerClosed: %v", mgr)

	mgr.server.Back()
	mgr.server = nil
}

// onNewSession 新连接
func (mgr *agentUserManager) onNewSession(sess base.Session) {
	log.RunLogger.Printf("agentUserManager.onNewSession sess[%v]: %v", sess, mgr)

	agent := mgr.agentPool.apply()
	mgr.mapUserAgent[sess.UUID()] = agent
	agent.Start(sess, mgr)
}

// onAgentClosed Agent关闭
func (mgr *agentUserManager) onAgentClosed(agent *agentUser) {
	log.RunLogger.Printf("agentUserManager.onAgentClosed %v: %v", agent, mgr)

	delete(mgr.mapUserAgent, agent.uuidSession)

	// 回收清理
	agent.Back()

	// 缓存
	mgr.agentPool.back(agent)
}

// mainLoop
func (mgr *agentUserManager) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("agentUserManager.mainLoop start: %v", mgr)

	atomic.AddInt32(&waitServerQuit, 1)

	// 主循环
	{
	mainLoop:
		for {
			select {
			case sess := <-mgr.chNewSession: // 新连接
				mgr.onNewSession(sess)

			case agent := <-mgr.chAgentClosed: // 连接结束
				mgr.onAgentClosed(agent)

			case <-chApplicationQuit: // 进程退出
				mgr.server.StopAccept()
				break mainLoop
			}
		}
	}

	log.RunLogger.Printf("agentUserManager.mainLoop start application quit: %v", mgr)

	// 等待底层服务器退出完成
	{
		<-mgr.chServerClosed

		mgr.onBaseServerClosed()
	}

	log.RunLogger.Printf("agentUserManager.mainLoop application quit step 1: recv base.Server closed: %v", mgr)

	// 继续处理新连接(直接关闭)
	{
	endNewSession:
		for {
			select {
			case sess := <-mgr.chNewSession: // 新连接
				// todo: 此分支需要测试
				// 直接关闭
				sess.Close()
			default:
				break endNewSession
			}
		}
	}

	log.RunLogger.Printf("agentUserManager.mainLoop application quit step 2: close all new wait session: %v", mgr)

	// 关闭所有已建立的连接
	{
		if len(mgr.mapUserAgent) > 0 {
			// 向其通知退出
			for _, agent := range mgr.mapUserAgent {
				agent.Close()
			}

			log.RunLogger.Printf("agentUserManager.mainLoop application quit step 3: notify user agent close: %v", mgr)

			// 等待全部退出
		endSession:
			for {
				select {
				case agent := <-mgr.chAgentClosed: // 连接结束
					mgr.onAgentClosed(agent)

					// 全关闭了
					if len(mgr.mapUserAgent) == 0 {
						break endSession
					}
				}
			}
		}

		log.RunLogger.Printf("agentUserManager.mainLoop application quit step 4: all user agent closed: %v", mgr)
	}
}

// mainLoopEnd
func (mgr *agentUserManager) mainLoopEnd() {
	log.RunLogger.Printf("agentUserManager.mainLoopEnd")

	atomic.AddInt32(&waitServerQuit, -1)
}

// init 根据配置初始化Server
func (mgr *agentUserManager) start(config *serveConfig) (err error) {
	mgr.sendExtraDataType, err = ffProto.GetExtraDataType(config.SendExtraDataType)
	if err != nil {
		return err
	}

	mgr.recvExtraDataType, err = ffProto.GetExtraDataType(config.RecvExtraDataType)
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

	mgr.config = config

	mgr.server, mgr.uuidServer = server, server.UUID()
	mgr.chNewSession, mgr.chServerClosed = chNewSession, chServerClosed

	mgr.chAgentClosed = make(chan *agentUser, 2)
	mgr.mapUserAgent = make(map[uuid.UUID]*agentUser, config.InitOnlineCount)
	mgr.agentPool = newAgentUserPool(config.InitOnlineCount)

	go util.SafeGo(mgr.mainLoop, mgr.mainLoopEnd)

	return nil
}
