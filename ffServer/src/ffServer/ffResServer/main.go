package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"time"

	"net"
	"net/http"
)

func main() {
	var err error
	defer util.PanicProtect(func(isPanic bool) {
		if isPanic {
			log.RunLogger.Println("异常退出, 以上是错误堆栈")
			<-time.After(time.Hour)
		} else if err != nil {
			util.PrintPanicStack(err)
			log.RunLogger.Println("启动出错, 以上是错误堆栈")
			<-time.After(time.Hour)
		}
	}, "ffResServer")

	// 读取配置文件
	serverConfig, err = loadServerConfig("cfg/cfg.cfg")
	if err != nil {
		return
	}

	// 监听的地址
	log.RunLogger.Println("devChannel: " + serverConfig.devChannel)
	log.RunLogger.Println("SERVE_IP_PORT: " + serverConfig.listenIPPort)
	log.RunLogger.Println("NTF_CLIENT_IP_PORT: " + serverConfig.outerIPPort)
	log.RunLogger.Printf("MAX_CONNECTION_COUNT: %d\n", serverConfig.connectionLimit)
	log.RunLogger.Println("")
	log.RunLogger.Println("")

	// 渠道数据
	genAllChannels()

	// 远端资源
	genRemoteResMap()

	// 打印渠道数据
	for channelName, channelInfo := range globalChannelInfo {
		log.RunLogger.Println("channelName: " + channelName)
		log.RunLogger.Println("full_package_version: " + channelInfo.fullPackageVersion.String())
		log.RunLogger.Println("newest_version: " + channelInfo.newestVersion.String())
		log.RunLogger.Println("examine_version: " + channelInfo.examineVersion.String())
		log.RunLogger.Println("hotResFileMd5: " + channelInfo.hotResFileMd5)
		log.RunLogger.Println("select_server_ip: " + serverConfig.channelConfig[channelName]["SELECT_SERVER_IP"])
		log.RunLogger.Println("down_url: " + serverConfig.channelConfig[channelName]["FULL_DOWN_URL"])
		log.RunLogger.Println("")
		log.RunLogger.Println("")
	}

	// 响应客户端的请求
	http.HandleFunc("/check", getCheck)
	http.HandleFunc("/apk", getApk)
	// http.HandleFunc("/hotfix", get_hotfix)
	http.HandleFunc("/res", getRemoteResObsolete)
	http.HandleFunc("/remoteres", getRemoteRes)
	http.HandleFunc("/hotres", getHotRes)

	l, err := net.Listen("tcp", serverConfig.listenIPPort)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	l = LimitListener(l, serverConfig.connectionLimit)

	srv := &http.Server{
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 10*1024*1024 + 1024, // 10M + 1K
	}
	srv.SetKeepAlivesEnabled(false)
	srv.Serve(l)
}
