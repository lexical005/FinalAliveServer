package main

import (
	"ffCommon/log/log"
	"ffCommon/util"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"
)

type serveLogin struct {
	s *http.Server

	listener net.Listener
}

// handleLoginRequest 响应登录验证
func (serve *serveLogin) handleLoginRequest(rw http.ResponseWriter, req *http.Request) {
	// 只响应 POST 方式的请求
	if req.Method != "POST" {
		return
	}

	// 异常保护
	defer util.PanicProtect(nil)

	// 成功
	defer func() {
		rw.Write([]byte(serverIAPResponse))
	}()

	// 解析参数
	req.ParseForm()

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.RunLogger.Printf("serveLogin.handleLoginRequest err[%v]", err)
		return
	}

	defer req.Body.Close()

	s := string(content)
	s = strings.Replace(s, "\r\n", "", -1)
	s = strings.Replace(s, "\n", "", -1)

	log.RunLogger.Println("onServerRequestIAP:", s)
}

// mainLoop
func (serve *serveLogin) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("serveLogin.mainLoop start")

	atomic.AddInt32(&waitApplicationQuit, 1)

	// 响应客户端的请求
	http.HandleFunc("/login", serve.handleLoginRequest)

	serve.s.Serve(serve.listener)
}

// mainLoopEnd
func (serve *serveLogin) mainLoopEnd(isPanic bool) {
	log.RunLogger.Printf("serveLogin.mainLoopEnd isPanic[%v]", isPanic)

	atomic.AddInt32(&waitApplicationQuit, -1)
}

func (serve *serveLogin) start() (err error) {
	// 建立监听
	serve.listener, err = net.Listen("tcp", appConfig.ServeLogin.ListenAddr)
	if err != nil {
		return fmt.Errorf("serveLogin.start err[%v]", err)
	}

	serve.s = &http.Server{
		Handler:        nil,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 10 * 1024,
	}

	// 登录服务
	go util.SafeGo(serve.mainLoop, serve.mainLoopEnd)

	// 登录退出
	go util.SafeGo(serve.waitExit, nil)

	return nil
}

func (serve *serveLogin) waitExit(params ...interface{}) {
	select {
	case <-chApplicationQuit:
		// 返回时, 即优雅了的关闭了http服务器
		serve.s.Shutdown(nil)
		break
	}
}
