package netmanager

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/uuid"
	"fmt"
	"sync/atomic"
	"time"
)

// Manager Agent管理器
type Manager struct {
	name string

	// 接收KeepAlive超时, 通过Session的ReadDeadTime实现
	// sendKeepAliveLeftTime 隔多久再次向对端发送KeepAlive协议, 只有当Manager作为Client时才需要
	sendKeepAliveLeftTime time.Duration
	// nextSendKeepAliveTime 下一次向对端发送KeepAlive协议的时间, 只有当Manager作为Client时才需要
	nextSendKeepAliveTime time.Time
	// sendKeepAliveInterval 定时发送KeepAlive协议的间隔
	sendKeepAliveInterval time.Duration

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
	agent.init(sess, mgr.net, mgr.chAgentClosed, mgr.sendKeepAliveInterval == 0)

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

	//
	mgr.net.onAgentClosed()
}

// 通知现有每个连接向对端发送KeepAlive, 只有Manager作为Client时才会触发
func (mgr *Manager) sendKeepAlive() {
	log.RunLogger.Printf("%v.sendKeepAlive sendKeepAliveLeftTime[%v]", mgr.name, mgr.sendKeepAliveLeftTime)

	mgr.sendKeepAliveLeftTime += mgr.sendKeepAliveInterval
	mgr.nextSendKeepAliveTime = time.Now().Add(mgr.sendKeepAliveLeftTime)

	for _, agent := range mgr.mapAgent {
		agent.keepAlive()
	}
}

// mainLoop
func (mgr *Manager) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("%v.mainLoop start", mgr.name)

	atomic.AddInt32(mgr.countApplicationQuit, 1)

	// 主循环
	{
		if mgr.sendKeepAliveInterval == 0 { //无KeepAlive

		mainLoopNoKeepAlive:
			for {
				select {
				case sess := <-mgr.net.NewSessionChan(): // 新连接
					mgr.onNewSession(sess)

				case agent := <-mgr.chAgentClosed: // 连接结束
					mgr.onAgentClosed(agent)

				case <-mgr.chApplicationQuit: // 进程退出
					mgr.net.Stop()
					break mainLoopNoKeepAlive
				}
			}

		} else { // 有KeepAlive

			// 最近发送KeepAlive时间
			mgr.sendKeepAliveLeftTime = mgr.sendKeepAliveInterval
			mgr.nextSendKeepAliveTime = time.Now().Add(mgr.sendKeepAliveLeftTime)

		mainLoopKeepAlive:
			for {
				triggerKeepAlive := false
				select {
				case sess := <-mgr.net.NewSessionChan(): // 新连接
					mgr.onNewSession(sess)

				case agent := <-mgr.chAgentClosed: // 连接结束
					mgr.onAgentClosed(agent)

				case <-time.After(mgr.sendKeepAliveLeftTime):
					triggerKeepAlive = true
					mgr.sendKeepAliveLeftTime = time.Now().Sub(mgr.nextSendKeepAliveTime)

				case <-mgr.chApplicationQuit: // 进程退出
					mgr.net.Stop()
					break mainLoopKeepAlive
				}

				// 发送KeepAlive
				if !triggerKeepAlive {
					mgr.sendKeepAliveLeftTime = mgr.nextSendKeepAliveTime.Sub(time.Now())
					triggerKeepAlive = mgr.sendKeepAliveLeftTime < 1
				}

				if triggerKeepAlive {
					mgr.sendKeepAlive()
				}
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
