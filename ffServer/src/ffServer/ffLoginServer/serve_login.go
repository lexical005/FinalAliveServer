package main

import (
	"encoding/json"
	"ffCommon/log/log"
	"ffCommon/util"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sync/atomic"
	"time"
)

type serveLogin struct {
	s *http.Server

	listener net.Listener
}

// handleLoginCustom 响应自定义登录验证
func (serve *serveLogin) handleLoginCustom(rw http.ResponseWriter, req *http.Request) {
	// 只响应 POST 方式的请求
	if req.Method != "POST" {
		return
	}

	// 携带的数据格式
	type customLoginData struct {
		UUIDPlatform string
		UUIDAccount  uint64
		Result       int32 // 0 无错 1 异常返回 2 UUIDPlatform长度不对
	}
	data := &customLoginData{
		Result: 1,
	}

	// 异常保护
	defer util.PanicProtect(nil)

	// 反馈
	defer func() {
		t, err := json.Marshal(data)
		if err != nil {
			log.RunLogger.Printf("serveLogin.handleLoginCustom json.Marshal data[%v] error[%v]", data, err)
			return
		}

		log.RunLogger.Printf("serveLogin.handleLoginCustom login result data[%v]", data)

		rw.Write(t)
	}()

	// 解析参数
	req.ParseForm()

	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.RunLogger.Printf("serveLogin.handleLoginCustom ioutil.ReadAll error[%v]", err)
		return
	}
	defer req.Body.Close()

	// 解析数据
	err = json.Unmarshal(content, data)
	if err != nil {
		log.RunLogger.Printf("serveLogin.handleLoginCustom json.Unmarshal content[%v] get error[%v]", string(content), err)
		return
	}

	// UUIDPlatform 有效性判定
	UUIDPlatform := []byte(data.UUIDPlatform)
	if len(UUIDPlatform) < 1 || len(UUIDPlatform) > 16 {
		data.Result = 2
		log.RunLogger.Printf("serveLogin.handleLoginCustom invalid UUIDPlatform length content[%v] get error[%v]", string(content), err)
		return
	}

	// 生成 UUIDAccount
	data.Result = 0
	data.UUIDAccount = 0
	for i := 0; i < len(UUIDPlatform); i++ {
		data.UUIDAccount += (uint64(UUIDPlatform[i])) << (uint(i) * 4)
	}
}

// mainLoop
func (serve *serveLogin) mainLoop(params ...interface{}) {
	log.RunLogger.Printf("serveLogin.mainLoop start")

	atomic.AddInt32(&waitApplicationQuit, 1)

	// 响应客户端的请求
	http.HandleFunc("/login", serve.handleLoginCustom)

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
