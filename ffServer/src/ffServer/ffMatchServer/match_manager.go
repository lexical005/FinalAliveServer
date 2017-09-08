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

// matchManager 匹配管理器
type matchManager struct {
	groups []*matchGroup

	matchProto chan *ffProto.Proto // 开始/取消匹配

	uuidBattleGenerator      uuid.Generator
	uuidBattleTokenGenerator uuid.Generator
}

// Start 初始化
func (mgr *matchManager) Start() (err error) {
	mgr.groups = []*matchGroup{
		newMatchGroup(matchModeSingle),
		newMatchGroup(matchModeDouble),
		nil,
		newMatchGroup(matchModeFour),
	}

	mgr.matchProto = make(chan *ffProto.Proto, appConfig.Match.InitMatchCount)
	mgr.uuidBattleGenerator, err = uuid.NewGeneratorSafe(uint64(appConfig.Server.ServerID))
	if err != nil {
		return err
	}
	mgr.uuidBattleTokenGenerator, err = uuid.NewGeneratorSafe(0)
	if err != nil {
		return err
	}

	go util.SafeGo(mgr.mainLoop, mgr.mainLoopEnd)

	return nil
}

// GetMatchGroup 获取匹配模式
func (mgr *matchManager) GetMatchGroup(mode matchMode) *matchGroup {
	return mgr.groups[mode-1]
}

// OnPlayerMatchProto 用户匹配相关协议
func (mgr *matchManager) OnPlayerMatchProto(proto *ffProto.Proto) bool {
	proto.SetCacheWaitDispatch()
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

	changedToSendState := false

	proto.SetCacheDispatched()
	defer func() {
		if !changedToSendState {
			proto.BackAfterDispatch()
		}
	}()

	uuidPlayerKey := uuid.NewUUID(proto.ExtraData())
	player := instMatchPlayerMgr.GetPlayer(uuidPlayerKey)
	if player == nil {
		log.RunLogger.Printf("matchManager.onProto, can not find player[%v]", uuidPlayerKey)
		return
	}

	switch proto.ProtoID() {
	case ffProto.MessageType_StartMatch:
		mgr.doStartMatch(player, proto)
	case ffProto.MessageType_StopMatch:
		mgr.doStopMatch(player, proto)
	case ffProto.MessageType_LeaveMatchServer:
		mgr.leaveMatchServer(player, proto)
	}

	changedToSendState = instAgentGameServerMgr.SendProtoExtraDataUUID(player, proto, true)
}

// 用户离开匹配服务器
func (mgr *matchManager) leaveMatchServer(player *matchPlayer, proto *ffProto.Proto) bool {
	player.StopMatch()

	instMatchPlayerMgr.RemovePlayer(player.uuidPlayerKey)
	return false
}

// doStartMatch 开始匹配
func (mgr *matchManager) doStartMatch(player *matchPlayer, proto *ffProto.Proto) {
	message, _ := proto.Message().(*ffProto.MsgStartMatch)

	result := false
	mode := matchMode(message.MatchMode)
	if matchModeSingle == mode || matchModeDouble == mode || matchModeFour == mode {
		result = player.StartMatch(mode)
	}

	if !result {
		message.Result = ffError.ErrUnknown.Code()
	}
}

// doStopMatch 停止匹配
func (mgr *matchManager) doStopMatch(player *matchPlayer, proto *ffProto.Proto) {
	message, _ := proto.Message().(*ffProto.MsgStopMatch)

	result := player.StopMatch()
	if !result {
		message.Result = ffError.ErrUnknown.Code()
	}
}

// mainLoopEnd
func (mgr *matchManager) mainLoopEnd(isPanic bool) {
	log.RunLogger.Printf("matchGroupManager.mainLoopEnd")

	atomic.AddInt32(&waitApplicationQuit, -1)
}
