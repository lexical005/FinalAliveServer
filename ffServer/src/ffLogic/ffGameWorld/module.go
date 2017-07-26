package ffGameWorld

import (
	"ffLogic/ffDef"
)

// 游戏世界
var world = &gameWorld{}

// 游戏世界框架
var worldFrame ffDef.IGameWorldFrame

// pool of account
var poolOfAccount = &accountPool{}

// manager of account
var managerOfAccount = &accountManager{}

// NewGameWorld 新建并初始化游戏世界
func NewGameWorld(frame ffDef.IGameWorldFrame) (ffDef.IGameWorld, error) {
	worldFrame = frame

	err := world.init()
	if err != nil {
		return nil, err
	}

	return world, nil
}
