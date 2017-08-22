package main

import (
	"ffCommon/log/log"
	"ffCommon/net/httpclient"
	"ffCommon/util"
	"fmt"
	"time"

	"sync/atomic"
)

type httpClientLogin struct {
	url         string
	workerCount int
	client      *httpclient.Client

	status int32 // 运行状态  0初始状态 1可使用 2使用中 -1关闭中(阻碍进入使用状态) -2完成了所有关闭工作(阻碍进入使用状态) -3关闭完成

	chPostRequest chan httpclient.Request

	chExit chan struct{} // 需要收集齐config.RequestWorkerCount个退出完成事件, 才表明底层彻底退出了
}

func (c *httpClientLogin) String() string {
	return fmt.Sprintf("url[%v] workerCount[%v] status[%v] c.chPostRequest[%v]",
		c.url, c.workerCount, c.status, len(c.chPostRequest))
}

func (c *httpClientLogin) Start() error {
	c.url, c.workerCount = appConfig.ConnectLogin.URL, appConfig.ConnectLogin.RequestWorkerCount

	c.chExit = make(chan struct{}, 1)
	c.chPostRequest = make(chan httpclient.Request, appConfig.ConnectLogin.RequestCountCache)

	c.client = httpclient.NewClient(c.url, c.chExit, c.chPostRequest)
	c.client.Start(c.workerCount)

	// 0初始状态 ==> 1可使用
	atomic.CompareAndSwapInt32(&c.status, 0, 1)

	go util.SafeGo(c.mainLoop, c.mainLoopEnd)

	return nil
}

func (c *httpClientLogin) mainLoop(params ...interface{}) {
	log.RunLogger.Println("httpClientLogin.mainLoop")

	atomic.AddInt32(&waitApplicationQuit, 1)

	// 等待进程退出
	<-chApplicationQuit

	// 2使用中 ==> -1关闭中(阻碍进入使用状态)
	if !atomic.CompareAndSwapInt32(&c.status, 2, -1) {
		// -2完成了所有关闭工作(阻碍进入使用状态)
		atomic.StoreInt32(&c.status, -2)
	}

	// 通知底层退出
	c.chPostRequest <- nil

	// 等待所有工作协程退出
	for i := 0; i < c.workerCount; i++ {
		<-c.chExit
	}
}

func (c *httpClientLogin) mainLoopEnd(isPanic bool) {
	log.RunLogger.Println("httpClientLogin.mainLoopEnd", isPanic)

	// 等待发送方法执行完毕
	{
		waitCount, maxWaitCount := 0, 10
		for {
			// -2完成了所有关闭工作(阻碍进入使用状态) ==> -3关闭完成
			if atomic.CompareAndSwapInt32(&c.status, -2, -3) {
				break
			}

			// 等待1秒
			<-time.After(time.Second)

			waitCount++
			if waitCount > maxWaitCount {
				log.FatalLogger.Printf("httpClientLogin.mainLoopEnd wait status change to -3 too long time[%v] to exit", waitCount)
				break
			}
		}
	}

	// 清理
	close(c.chPostRequest)
	for request := range c.chPostRequest {
		if request != nil {
			// todo: do what?
		}
	}
	c.chPostRequest = nil

	close(c.chExit)
	c.chExit = nil

	atomic.AddInt32(&waitApplicationQuit, -1)
}

func (c *httpClientLogin) PostCustomLogin(request *httpClientCustomLoginData) bool {
	// 1可使用 ==> 2使用中
	if !atomic.CompareAndSwapInt32(&c.status, 1, 2) {
		return false
	}

	c.chPostRequest <- request

	// 2使用中 ==> 1可使用
	if !atomic.CompareAndSwapInt32(&c.status, 2, 1) {
		// -2完成了所有关闭工作(阻碍进入使用状态)
		atomic.StoreInt32(&c.status, -2)
	}

	return true
}
