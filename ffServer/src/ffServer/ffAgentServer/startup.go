package main

import (
	"ffCommon/log/log"
	"ffCommon/log/logfile"
	"ffCommon/net/session"
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

	// 初始化Session
	err = session.Init(appConfig.Session.ReadDeadTime, appConfig.Session.OnlineCount)
	if err != nil {
		return err
	}

	err = clientAgentMgr.init()
	if err != nil {
		return err
	}

	err = serverAgentMgr.init()
	if err != nil {
		return err
	}

	err = clientAgentPool.init()
	if err != nil {
		return err
	}

	return err
}
