package netmanager

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/uuid"
	"fmt"
	"sync/atomic"
)

// Manager Agent管理器
type Manager struct {
	name string

	// handlerManager
	handlerManager INetSessionHandlerManager

	// net 底层实现
	net inet

	// chAgentClosed 用于接收Agent关闭事件
	chAgentClosed chan *agentSession
	// mapAgent 所有连接
	mapAgent map[uuid.UUID]*agentSession
	// agentPool 所有连接缓存
	agentPool *agentSessionPool

	// countApplicationQuit 退出时计数
	countApplicationQuit *int32
	// chApplicationQuit 外界通知退出
	chApplicationQuit chan struct{}
}

// Status 内部状态
func (mgr *Manager) Status() string {
	return fmt.Sprintf("%v net[%v] chAgentClosed[%v]",
		mgr.name, mgr.net, len(mgr.chAgentClosed))
}

func (mgr *Manager) String() string {
	return mgr.name
}

// 底层退出
func (mgr *Manager) onNetExit() {
	log.RunLogger.Printf("%v.onNetExit", mgr.name)

	mgr.net.BackNet()
}

// onNewSession 新连接
func (mgr *Manager) onNewSession(sess base.Session) {
	log.RunLogger.Printf("%v.onNewSession sess[%v]", mgr.name, sess)

	//
	agent := mgr.agentPool.apply()
	agent.init(sess, mgr.net, mgr.chAgentClosed)

	//
	handler := mgr.handlerManager.Create(agent)

	//
	mgr.mapAgent[agent.UUID()] = agent

	agent.Start(sess, mgr.net, handler)
}

// onAgentClosed Agent关闭
func (mgr *Manager) onAgentClosed(agent *agentSession) {
	log.RunLogger.Printf("%v.onAgentClosed %v", mgr.name, agent)

	delete(mgr.mapAgent, agent.UUID())

	mgr.handlerManager.Back(agent.handler)

	// 回收清理
	agent.Back()

	// 缓存
	mgr.agentPool.back(agent)
}

// mainLoop
func (mgr *Manager) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("%v.mainLoop start", mgr.name)

	atomic.AddInt32(mgr.countApplicationQuit, 1)

	// 主循环
	{
	mainLoop:
		for {
			select {
			case sess := <-mgr.net.NewSessionChan(): // 新连接
				mgr.onNewSession(sess)

			case agent := <-mgr.chAgentClosed: // 连接结束
				mgr.onAgentClosed(agent)

			case <-mgr.chApplicationQuit: // 进程退出
				mgr.net.Stop()
				break mainLoop
			}
		}
	}

	log.RunLogger.Printf("%v.mainLoop start application quit", mgr.name)

	// 等待底层服务器退出完成
	{
		mgr.net.WaitNetExit()

		mgr.onNetExit()
	}

	log.RunLogger.Printf("%v.mainLoop application quit step 1: recv base.Server closed", mgr.name)

	// 继续处理新连接(直接关闭)
	{
	endNewSession:
		for {
			select {
			case sess := <-mgr.net.NewSessionChan(): // 新连接
				// todo: 此分支需要测试
				// 直接关闭
				sess.Close()
			default:
				break endNewSession
			}
		}
	}

	log.RunLogger.Printf("%v.mainLoop application quit step 2: close all new wait session", mgr.name)

	// 关闭所有已建立的连接
	{
		if len(mgr.mapAgent) > 0 {
			// 向其通知退出
			for _, agent := range mgr.mapAgent {
				agent.Close()
			}

			log.RunLogger.Printf("%v.mainLoop application quit step 3: notify user agent close", mgr.name)

			// 等待全部退出
		endSession:
			for {
				select {
				case agent := <-mgr.chAgentClosed: // 连接结束
					mgr.onAgentClosed(agent)

					// 全关闭了
					if len(mgr.mapAgent) == 0 {
						break endSession
					}
				}
			}
		}

		log.RunLogger.Printf("%v.mainLoop application quit step 4: all user agent closed", mgr.name)
	}
}

// mainLoopEnd
func (mgr *Manager) mainLoopEnd(isPanic bool) {
	log.RunLogger.Printf("%v.mainLoopEnd isPanic[%v]", mgr.name, isPanic)

	mgr.net.Clear()

	close(mgr.chAgentClosed)
	mgr.chAgentClosed = nil

	handlerManager := mgr.handlerManager
	mgr.handlerManager = nil

	atomic.AddInt32(mgr.countApplicationQuit, -1)

	handlerManager.End()
}
