package httpclient

import (
	"bytes"
	"ffCommon/log/log"
	"ffCommon/util"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

// Client http Client
//	从请求Close开始, 已发起了的请求, 依然会返回结果, 管道内尚未处理的, 则全部忽略
type Client struct {
	client *http.Client

	chRequest        chan *PostRequest // 待发送列表
	postRequestPool  *postRequestPool  // post请求对象缓存
	responseDataPool *responseDataPool // 响应结果缓存

	status int32     // 运行状态  0初始状态 1可使用 2使用中 -1关闭中(阻碍进入使用状态) -2完成了所有关闭工作(阻碍进入使用状态) -3关闭完成
	once   util.Once // 一次性关闭

	workGoroutineCount int           // 工作协程数量
	chWaitWorkerExit   chan struct{} // 等待工作协程退出

	chExit               chan struct{} // 退出完成时, 向外界通知, 仅有使用权
	countApplicationQuit *int32        // countApplicationQuit 退出时计数
	chApplicationQuit    chan struct{} // chApplicationQuit 外界通知退出
}

// Start 启动, countGoroutine发送请求协程数量
func (c *Client) Start(countGoroutine int) {
	// 可使用
	atomic.StoreInt32(&c.status, 1)

	// 系统模块增加一个
	atomic.AddInt32(c.countApplicationQuit, 1)

	// 工作协程
	c.workGoroutineCount = countGoroutine
	c.chRequest = make(chan *PostRequest, c.workGoroutineCount*2)
	c.postRequestPool = newPostRequestPool("httpClient", c.workGoroutineCount*2)
	c.responseDataPool = newResponseDataPool("httpClient", c.workGoroutineCount*2)
	for index := 0; index < countGoroutine; index++ {
		go util.SafeGo(c.mainRequestLoop, c.mainRequestLoopEnd)
	}

	// 等待工作协程全部退出
	c.chWaitWorkerExit = make(chan struct{}, c.workGoroutineCount)
	go util.SafeGo(c.waitWorkerExitLoop, c.waitWorkerExitLoopEnd)
}

// ApplyPostRequest 申请一个Post请求
//	data是原始数据经json.Unmarshal序列化后的字节流
func (c *Client) ApplyPostRequest(url string, data []byte, chResponseData chan *ResponseData) *PostRequest {
	uuid := uuidGenerator.Gen()
	request := c.postRequestPool.apply()
	request.init(uuid, url, data, chResponseData)
	return request
}

// BackResponseData 归还请求结果
func (c *Client) BackResponseData(response *ResponseData) {
	response.back()

	c.responseDataPool.back(response)
}

// Post 以Post方式, 发送json数据
//	返回值: 是否加入了待Post的缓存管道内
func (c *Client) Post(request *PostRequest) bool {
	// 1可使用 ==> 2使用中
	if !atomic.CompareAndSwapInt32(&c.status, 1, 2) {
		return false
	}

	c.chRequest <- request

	// 2使用中 ==> 1可使用
	if !atomic.CompareAndSwapInt32(&c.status, 2, 1) {
		// -2完成了所有关闭工作(阻碍进入使用状态)
		atomic.StoreInt32(&c.status, -2)
	}

	return true
}

func (c *Client) doPost(request *PostRequest) {
	var err error

	defer func() {
		request.back()
		c.postRequestPool.back(request)
	}()

	response := c.responseDataPool.apply()

	resp, err := c.client.Post(request.url, postContentTypeJSON, bytes.NewBuffer(request.data))
	if err != nil {
		request.onError(err, response)
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		request.onError(err, response)
		return
	}

	request.onResponse(content, response)
}

func (c *Client) doClose() {
	// 2使用中 ==> -1关闭中(阻碍进入使用状态)
	if !atomic.CompareAndSwapInt32(&c.status, 2, -1) {
		// -2完成了所有关闭工作(阻碍进入使用状态)
		atomic.StoreInt32(&c.status, -2)
	}
}

func (c *Client) doClear() {
	c.chExit = nil
	c.chApplicationQuit = nil

	close(c.chWaitWorkerExit)
	c.chWaitWorkerExit = nil

	close(c.chRequest)
	for data := range c.chRequest {
		if data != nil {
			c.postRequestPool.back(data)
		} else {
			break
		}
	}
	c.chRequest = nil
}

// mainRequestLoop
func (c *Client) mainRequestLoop(params ...interface{}) {
	log.RunLogger.Printf("Client.mainRequestLoop start")

	for {
		select {
		case <-c.chApplicationQuit:
			c.once.Do(c.doClose)
			return
		case data := <-c.chRequest:
			c.doPost(data)
		}
	}
}

// mainRequestLoopEnd
func (c *Client) mainRequestLoopEnd(isPanic bool) {
	log.RunLogger.Printf("Client.mainRequestLoopEnd isPanic[%v]", isPanic)

	c.chWaitWorkerExit <- struct{}{}
}

func (c *Client) waitWorkerExitLoop(params ...interface{}) {
	log.RunLogger.Printf("Client.waitWorkerExitLoop start")

	// 等待工作协程退出
	for index := 0; index < c.workGoroutineCount; index++ {
		<-c.chWaitWorkerExit
	}
}

func (c *Client) waitWorkerExitLoopEnd(isPanic bool) {
	log.RunLogger.Printf("Client.waitWorkerExitLoopEnd isPanic[%v]", isPanic)

	// 全部工作协程全部退出了, 通知外界
	c.chExit <- struct{}{}

	// 清理
	c.doClear()

	// 系统模块减少一个
	atomic.AddInt32(c.countApplicationQuit, -1)
}
