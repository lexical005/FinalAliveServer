package main

import (
	"ffCommon/log/log"
	"ffCommon/log/logfile"
	"ffCommon/net/tcpsession"
	"fmt"
)

func startup() (err error) {
	// 解析
	err = readToml()
	if err != nil {
		return err
	}

	// 输出配置文件
	log.RunLogger.Println(appConfig)

	// 初始化log
	relativePath := appConfig.Logger.RelativePath
	if relativePath == "" {
		relativePath = logfile.DefaultLogFileRelativePath
	}

	fileLenLimit := appConfig.Logger.FileLenLimit
	if fileLenLimit == 0 {
		fileLenLimit = logfile.DefaultLogFileLengthLimit
	}
	err = logfile.Init(
		relativePath, fileLenLimit,
		appConfig.Logger.RunLogger, appConfig.Logger.RunLoggerPrefix,
		true, logfile.DefaultLogFileFatalPrefix)
	if err != nil {
		return err
	}

	// 服务器编号
	if appConfig.Server.ServerID < 1 || appConfig.Server.ServerID >= (1<<uuidRequesterBitServerID) {
		return fmt.Errorf("startup: appConfig.Server.ServerID[%d] must in range[1-%d)", appConfig.Server.ServerID, 1<<uuidRequesterBitServerID)
	}

	// 初始化Session
	err = tcpsession.Init(appConfig.Session.ReadDeadTime, appConfig.Session.OnlineCount, tcpsession.DefaultInitSessionNetEventDataCount)
	if err != nil {
		return err
	}

	err = agentServerMgr.init()
	if err != nil {
		return err
	}

	// 初始化游戏世界框架
	worldFrame.init()

	return err
}
