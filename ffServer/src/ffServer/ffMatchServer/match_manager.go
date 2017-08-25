package main

import (
	"ffAutoGen/ffError"
	"ffCommon/log/log"
	"ffCommon/util"
	"ffCommon/uuid"
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

	mgr.matchProto = make(chan *ffProto.Proto, appConfig.Match.InitMatchCount)

	go util.SafeGo(mgr.mainLoop, mgr.mainLoopEnd)

	return nil
}

// GetMatchGroup 获取匹配模式
func (mgr *matchManager) GetMatchGroup(mode matchMode) *matchGroup {
	return mgr.groups[mode-1]
}

// OnPlayerMatchProto 用户匹配相关协议
func (mgr *matchManager) OnPlayerMatchProto(proto *ffProto.Proto) bool {
	mgr.matchProto <- proto
	return true
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
					if group != nil {
						group.Match()
					}
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
	log.RunLogger.Printf("matchManager.onProto proto[%v]", proto)
	player := instMatchPlayerMgr.GetPlayer(uuid.NewUUID(proto.ExtraData()))

	switch proto.ProtoID() {
	case ffProto.MessageType_StartMatch:
		mgr.doStartMatch(player, proto)
	case ffProto.MessageType_StopMatch:
		mgr.doStopMatch(player, proto)
	}

	ffProto.SendProtoExtraDataUUID(player.sourceServer, player.uuidPlayerKey, proto, true)
}

// doStartMatch 开始匹配
func (mgr *matchManager) doStartMatch(player *matchPlayer, proto *ffProto.Proto) {
	message, _ := proto.Message().(*ffProto.MsgStartMatch)

	result := false
	if player != nil {
		mode := matchMode(message.MatchMode)

		if matchModeSingle == mode || matchModeDouble == mode || matchModeFour == mode {
			result = player.StartMatch(mode)
		}
	}

	if !result {
		message.Result = ffError.ErrUnknown.Code()
	}
}

// doStopMatch 停止匹配
func (mgr *matchManager) doStopMatch(player *matchPlayer, proto *ffProto.Proto) {
	message, _ := proto.Message().(*ffProto.MsgStopMatch)

	result := false
	if player != nil {
		result = player.StopMatch()
	}

	if !result {
		message.Result = ffError.ErrUnknown.Code()
	}
}

// mainLoopEnd
func (mgr *matchManager) mainLoopEnd(isPanic bool) {
	log.RunLogger.Printf("matchGroupManager.mainLoopEnd")

	atomic.AddInt32(&waitApplicationQuit, -1)
}
