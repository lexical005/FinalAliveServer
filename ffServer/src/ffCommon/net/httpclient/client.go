package httpclient

import (
	"bytes"
	"ffCommon/log/log"
	"ffCommon/util"
	"io/ioutil"
	"net/http"
)

// Client http Client
//	从请求Close开始, 已发起了的请求, 依然会返回结果, 管道内尚未处理的, 则全部忽略
type Client struct {
	url    string
	client *http.Client

	chExit    chan struct{} // 退出完成时, 向外界通知, 仅有使用权
	chRequest chan Request  // 待发送列表
}

// Start 启动, countGoroutine发送请求协程数量
func (c *Client) Start(workerCount int) {
	// 工作协程
	for index := 0; index < workerCount; index++ {
		go util.SafeGo(c.mainRequestLoop, c.mainRequestLoopEnd)
	}
}

func (c *Client) doPost(request Request) {
	var err error
	var content []byte

	defer func() {
		if err != nil {
			request.OnError(err)
		} else if content != nil {
			request.OnResponse(content)
		}
	}()

	data, err := request.Data()
	if err != nil {
		return
	}

	resp, err := c.client.Post(c.url, postContentTypeJSON, bytes.NewBuffer(data))
	if err != nil {
		return
	}

	defer resp.Body.Close()

	content, err = ioutil.ReadAll(resp.Body)
}

func (c *Client) doClear() {
	c.chExit = nil
	c.chRequest = nil
}

// mainRequestLoop
func (c *Client) mainRequestLoop(params ...interface{}) {
	log.RunLogger.Printf("Client.mainRequestLoop start")

	for {
		select {
		case request := <-c.chRequest:
			if request != nil {
				c.doPost(request)
			} else {
				return
			}
		}
	}
}

// mainRequestLoopEnd
func (c *Client) mainRequestLoopEnd(isPanic bool) {
	log.RunLogger.Printf("Client.mainRequestLoopEnd isPanic[%v]", isPanic)

	c.chExit <- struct{}{}
}
