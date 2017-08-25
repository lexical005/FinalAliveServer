package main

import (
	"ffCommon/log/log"
	"ffCommon/net/base"
	"ffCommon/util"

	"github.com/lexical005/toml"
)

// matchConfig 匹配配置
type matchConfig struct {
	// InitMatchCount 默认初始同时多少用户在匹配
	InitMatchCount int

	// ExpectMaxPlayerCount 期望的最大匹配人数
	ExpectMaxPlayerCount int

	// BattleMinPlayerCount 进行比赛所需的最少人数
	BattleMinPlayerCount int

	// StopEnterTime 达成比赛所需的最少人数后, 最大等待多久, 就不再添加新玩家进入此战场
	StopEnterTime int

	// WaitBattleMaxTime 达成比赛所需的最少人数后, 最大等待多久, 就战斗开启
	WaitBattleMaxTime int

	// MixMatchWaitTime 等待多久后, 允许进行混合匹配(如:单人进入多人战场)
	MixMatchWaitTime int
}

// 服务器配置
type applicationConfig struct {
	// Server 服务器自身描述
	Server *base.ServerInfo

	// Session 连接配置
	Session *base.SessionConfig

	// ServeAgentGameServer 服务AgentGameServer的配置
	ServeAgentGameServer *base.ServeConfig

	// ServeAgentBattleServer 服务BattleGameServer的配置
	ServeAgentBattleServer *base.ServeConfig

	// Logger 日志配置
	Logger *log.LoggerConfig

	// Match 匹配配置
	Match *matchConfig
}

// Check 配置检查
func (config *applicationConfig) Check() (err error) {
	err = config.Session.Check()
	if err != nil {
		return
	}

	return nil
}

func readToml() error {
	// 读取文件内容
	fileContent, err := util.ReadFile("toml/config.toml")
	if err != nil {
		return err
	}

	// 解析
	err = toml.Unmarshal(fileContent, appConfig)
	if err != nil {
		return err
	}

	return nil
}
