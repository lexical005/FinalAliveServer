package main

import (
	"ffCommon/log/log"
	"ffCommon/log/logfile"
	"ffCommon/net/tcpsession"
	"ffCommon/uuid"

	"github.com/davecgh/go-spew/spew"
)

// 初始化并启动
func startup() (err error) {
	// 解析
	err = readToml()
	if err != nil {
		return err
	}

	// 输出配置文件
	log.RunLogger.Printf("startup appConfig:\n%v", spew.Sdump(appConfig))

	// 检查配置
	err = appConfig.Check()
	if err != nil {
		return err
	}

	// 初始化log
	if appConfig.Logger.LoggerType == "file" {
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
	}
	log.RunLogger.Printf("application Config:\n%v", spew.Sdump(appConfig))

	// 初始化Session
	uuidSessionGenerator, err := uuid.NewGeneratorSafe(uint64(appConfig.Server.ServerID))
	if err != nil {
		return err
	}
	err = tcpsession.Init(appConfig.Session, uuidSessionGenerator)
	if err != nil {
		return err
	}

	// 启动
	err = instAgentUserServer.Start()
	if err != nil {
		return
	}

	// 启动
	err = instMatchServerClient.Start()
	if err != nil {
		return
	}

	// 启动
	err = instHTTPLoginClient.Start()
	if err != nil {
		return
	}

	return err
}
