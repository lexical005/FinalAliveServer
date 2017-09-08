package main

import (
	"ffCommon/log/log"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
	"sync"
)

type battleGameWorld struct {
	mapScene    map[uuid.UUID]*battleScene
	lockerScene sync.Mutex
}

// 新战场
func (world *battleGameWorld) NewScene(message *ffProto.MsgServerNewBattle) {
	log.RunLogger.Println("battleGameWorld.mainLoop")

	world.lockerScene.Lock()
	defer world.lockerScene.Unlock()

	scene := &battleScene{
		uuidBattle: uuid.NewUUID(message.UUIDBattle),

		chProto: make(chan *userProto, maxRoleCount),
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

func (world *battleGameWorld) CheckScene(agent *battleUser) (*battleScene, error) {
	battle, ok := world.mapScene[agent.uuidBattle]
	if !ok {
		return nil, fmt.Errorf("battleGameWorld.CheckBattle failed, agent[%v] uniqueid[%v] uuidBattle[%v]",
			agent.agent.UUID(), agent.uniqueid, agent.uuidBattle)
	}
	return battle, nil
}

func (world *battleGameWorld) Start() {
}

var instBattleGameWorld = &battleGameWorld{
	mapScene: make(map[uuid.UUID]*battleScene, 1),
}
