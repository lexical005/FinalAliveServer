package main

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/util"
	"ffProto"
	"sync/atomic"
	"time"
)

// matchManager 匹配组管理器
type matchManager struct {
	groups []*matchGroup

	matchProto chan *ffProto.Proto // 开始/取消匹配
}

// Start 初始化
func (mgr *matchManager) Start() error {
	mgr.groups = []*matchGroup{
		newMatchGroup(matchModeSingle),
		newMatchGroup(matchModeDouble),
		nil,
		newMatchGroup(matchModeFour),
	}

	go util.SafeGo(mgr.mainLoop, mgr.mainLoopEnd)

	return nil
}

// GetMatchGroup 获取匹配模式
func (mgr *matchManager) GetMatchGroup(mode matchMode) *matchGroup {
	return mgr.groups[mode]
}

// OnPlayerMatchProto 用户匹配相关协议
func (mgr *matchManager) OnPlayerMatchProto(proto *ffProto.Proto) {
	mgr.matchProto <- proto
}

// mainLoop
func (mgr *matchManager) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("matchGroupManager.mainLoop")

	atomic.AddInt32(&waitApplicationQuit, 1)

	{
	deadLoop:
		for {
			// 匹配
			select {
			case <-time.After(time.Second):
				// 匹配
				for _, group := range mgr.groups {
					group.Match()
				}

			case <-chApplicationQuit: // 进程退出
				break deadLoop
			}

			// 申请匹配处理
			{
			startMatch:
				for {
					select {
					case proto := <-mgr.matchProto:
						mgr.onProto(proto)

					default:
						break startMatch
					}
				}
			}
		}
	}
}

func (mgr *matchManager) onProto(proto *ffProto.Proto) {
	switch proto.ProtoID() {
	case ffProto.MessageType_EnterMatchServer:
		mgr.doStartMatch(proto)
	case ffProto.MessageType_LeaveMatchServer:
		mgr.doStartMatch(proto)
	case ffProto.MessageType_StartMatch:
		mgr.doStartMatch(proto)
	case ffProto.MessageType_StopMatch:
		mgr.doStopMatch(proto)
	}
}

// doStartMatch 开始匹配
func (mgr *matchManager) doStartMatch(proto *ffProto.Proto) {
	uuidPlayerKey := proto.ExtraData()
	message, _ := proto.Message().(*ffProto.MsgStartMatch)

	result := false
	player := instMatchPalyerMgr.GetPlayer(uuidPlayerKey)
	if player != nil {
		mode := matchMode(message.MatchMode)

		if matchModeSingle == mode || matchModeDouble == mode || matchModeFour == mode {
			result = player.StartMatch(mode)
		}
	}

	if !result {
		message.Result = ffError.ErrUnknown.Code()
	}

	player.sourceServer.SendProto(player.uuidPlayerKey, proto)
}

// doStopMatch 停止匹配
func (mgr *matchManager) doStopMatch(proto *ffProto.Proto) {
	uuidPlayerKey := proto.ExtraData()
	message, _ := proto.Message().(*ffProto.MsgStopMatch)

	result := false
	player := instMatchPalyerMgr.GetPlayer(uuidPlayerKey)
	if player != nil {
		result = player.StopMatch()
	}

	if !result {
		message.Result = ffError.ErrUnknown.Code()
	}

	player.sourceServer.SendProto(player.uuidPlayerKey, proto)
}

// mainLoopEnd
func (mgr *matchManager) mainLoopEnd(isPanic bool) {
	log.RunLogger.Printf("matchGroupManager.mainLoopEnd")

	atomic.AddInt32(&waitApplicationQuit, -1)
}
