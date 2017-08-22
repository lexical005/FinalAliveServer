package main

import (
	"ffCommon/log/log"
	"ffCommon/log/logfile"
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

	// 启动服务
	err = serveLoginInst.start()
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	return err
}