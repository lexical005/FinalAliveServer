package main

import (
	"ffCommon/log/log"
	"ffCommon/log/logfile"
	"ffCommon/util"
	"time"

	"net"
	"net/http"
)

var appConfig = &applicationConfig{}
var mysql = &mysqlManager{}

func main() {
	// 异常保护
	defer util.PanicProtect(nil)

	// 读取配置文件
	var err error
	err = readAppToml()
	if err != nil {
		panic(err)
	}
	log.RunLogger.Println("appConfig:")
	log.RunLogger.Println(appConfig)

	logfile.InitRunLog(logfile.DefaultLogFileRelativePath, logfile.DefaultLogFileLengthLimit, "")

	// 数据库配置
	mysql.start()

	// 响应客户端的请求
	http.HandleFunc("/client", onClientRequest)
	http.HandleFunc("/serverIAP", onServerRequestIAP)
	http.HandleFunc("/"+vivo.payKey, vivo.onPay)

	// 建立监听服务
	l, err := net.Listen("tcp", appConfig.Net.ListenAddr+":"+appConfig.Net.ListenPort)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	l = LimitListener(l, appConfig.Net.ConnectionLimit)

	srv := &http.Server{
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2 * 1024 * 1024,
	}
	srv.SetKeepAlivesEnabled(false)
	srv.Serve(l)
}
