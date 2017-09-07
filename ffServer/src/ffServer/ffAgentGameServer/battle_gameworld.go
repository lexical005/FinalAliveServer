package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"sync"
	"time"
)

const (
	worldTimerInterval = 20 * time.Microsecond // 逻辑50帧
)

type battleGameWorld struct {
	mapScene    map[uuid.UUID]*battleScene
	lockerScene sync.Mutex

	lastTimer time.Time
}

// 新战场
func (world *battleGameWorld) NewScene(message *ffProto.MsgServerNewBattle) {
	log.RunLogger.Println("battleGameWorld.mainLoop")

	world.lockerScene.Lock()
	defer world.lockerScene.Unlock()

	scene := &battleScene{
		uuidBattle: uuid.NewUUID(message.UUIDBattle),
	}
	scene.Init(message.UserTokens)

	log.RunLogger.Printf("battleGameWorld.NewScene scene[%v:%v]",
		scene.uuidBattle, scene.uuidTokens)

	world.mapScene[scene.uuidBattle] = scene
}

// 战场新增用户
func (world *battleGameWorld) OnUserEnterScene(message *ffProto.MsgServerBattleUserEnter) {
	world.lockerScene.Lock()
	defer world.lockerScene.Unlock()

	uuidBattle := uuid.NewUUID(message.UUIDBattle)
	uuidToken := uuid.NewUUID(message.UserToken)

	log.RunLogger.Printf("battleGameWorld.OnUserEnterScene scene[%v] token[%v]",
		uuidBattle, uuidToken)

	if scene, ok := world.mapScene[uuidBattle]; ok {
		scene.uuidTokens = append(scene.uuidTokens, uuidToken)
	}
}

// 战场移除用户
func (world *battleGameWorld) OnUserLeaveScene(message *ffProto.MsgServerBattleUserLeave) {
	world.lockerScene.Lock()
	defer world.lockerScene.Unlock()

	uuidBattle := uuid.NewUUID(message.UUIDBattle)
	uuidToken := uuid.NewUUID(message.UserToken)

	log.RunLogger.Printf("battleGameWorld.OnUserLeaveScene scene[%v] token[%v]",
		uuidBattle, uuidToken)

	if scene, ok := world.mapScene[uuidBattle]; ok {
		for i, token := range scene.uuidTokens {
			if token == uuidToken {
				scene.uuidTokens = append(scene.uuidTokens[:i], scene.uuidTokens[i+1:]...)
				break
			}
		}
	}
}

func (world *battleGameWorld) CheckBattle(agent *battleUser) (*battleScene, error) {
	battle, ok := world.mapScene[agent.uuidBattle]
	if !ok {
		return nil, fmt.Errorf("battleGameWorld.CheckBattle failed, agent[%v] uniqueid[%v] uuidBattle[%v]",
			agent.agent.UUID(), agent.uniqueid, agent.uuidBattle)
	}
	return battle, nil
}

func (world *battleGameWorld) update(passTime time.Duration) {
}

func (world *battleGameWorld) Start() {
	go util.SafeGo(world.mainLoop, world.mainLoopEnd)
}

func (world *battleGameWorld) mainLoop(params ...interface{}) {
	log.RunLogger.Println("battleGameWorld.mainLoop")

	world.lastTimer = time.Now()

	{
		waitTime := worldTimerInterval
		for {
			<-time.After(waitTime)

			// 时间驱动
			world.lastTimer = world.lastTimer.Add(worldTimerInterval)
			waitTime = worldTimerInterval - time.Now().Sub(world.lastTimer)
			world.update(worldTimerInterval)
		}
	}
}

func (world *battleGameWorld) mainLoopEnd(isPanic bool) {
	log.RunLogger.Println("battleGameWorld.mainLoopEnd", isPanic)
}

var instBattleGameWorld = &battleGameWorld{
	mapScene: make(map[uuid.UUID]*battleScene, 1),
}
