package main

import (
	"ffCommon/log/log"
	"ffCommon/net/httpclient"
	"ffCommon/util"
	"fmt"

	"sync/atomic"
)

type httpClientLogin struct {
	url         string
	workerCount int
	client      *httpclient.Client

	status util.Worker // 可使用性状态管理, 内含一次性关闭

	chPostRequest chan httpclient.Request

	chExit chan struct{} // 需要收集齐config.RequestWorkerCount个退出完成事件, 才表明底层彻底退出了
}

func (c *httpClientLogin) String() string {
	return fmt.Sprintf("url[%v] workerCount[%v] status[%v] c.chPostRequest[%v]",
		c.url, c.workerCount, c.status.String(), len(c.chPostRequest))
}

func (c *httpClientLogin) Start() error {
	c.url, c.workerCount = appConfig.ConnectLogin.URL, appConfig.ConnectLogin.RequestWorkerCount

	c.chExit = make(chan struct{}, 1)
	c.chPostRequest = make(chan httpclient.Request, appConfig.ConnectLogin.RequestCountCache)

	c.client = httpclient.NewClient(c.url, c.chExit, c.chPostRequest)
	c.client.Start(c.workerCount)

	c.status.Ready()

	go util.SafeGo(c.mainLoop, c.mainLoopEnd)

	return nil
}

func (c *httpClientLogin) mainLoop(params ...interface{}) {
	log.RunLogger.Println("httpClientLogin.mainLoop")

	atomic.AddInt32(&waitApplicationQuit, 1)

	// 等待进程退出
	<-chApplicationQuit

	c.status.Close()

	// 通知底层退出
	c.chPostRequest <- nil

	// 等待所有工作协程退出
	for i := 0; i < c.workerCount; i++ {
		<-c.chExit
	}
}

func (c *httpClientLogin) mainLoopEnd(isPanic bool) {
	log.RunLogger.Println("httpClientLogin.mainLoopEnd", isPanic)

	// 等待使用完毕
	c.status.WaitWorkEnd(10)

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
	work := c.status.EnterWork()

	defer func() {
		c.status.LeaveWork(work)
	}()

	if work {
		c.chPostRequest <- request
	}

	return work
}
