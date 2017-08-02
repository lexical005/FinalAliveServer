package main

import (
	"ffCommon/log/log"
	"ffCommon/net/tcpsession"
)

func startup() (err error) {
	// 解析
	err = readToml()
	if err != nil {
		return err
	}

	// 输出配置文件
	log.RunLogger.Printf("startup appConfig:\n%v", appConfig)

	// 初始化log
	// relativePath := appConfig.Logger.RelativePath
	// if relativePath == "" {
	// 	relativePath = logfile.DefaultLogFileRelativePath
	// }

	// fileLenLimit := appConfig.Logger.FileLenLimit
	// if fileLenLimit == 0 {
	// 	fileLenLimit = logfile.DefaultLogFileLengthLimit
	// }
	// err = logfile.Init(
	// 	relativePath, fileLenLimit,
	// 	appConfig.Logger.RunLogger, appConfig.Logger.RunLoggerPrefix,
	// 	true, logfile.DefaultLogFileFatalPrefix)
	// if err != nil {
	// 	return err
	// }

	// 初始化Session
	readDeadTime := tcpsession.DefaultReadDeadTime
	if appConfig.Session.ReadDeadTime > 0 {
		readDeadTime = appConfig.Session.ReadDeadTime
	}
	initNetEventDataCount := appConfig.Session.InitNetEventDataCount
	if initNetEventDataCount == 0 {
		initNetEventDataCount = appConfig.Session.InitOnlineCount / 4
	}
	if initNetEventDataCount < 2 {
		initNetEventDataCount = 2
	}
	err = tcpsession.Init(readDeadTime, appConfig.Session.InitOnlineCount, initNetEventDataCount)
	if err != nil {
		return err
	}

	return err
}
