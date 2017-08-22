package main

import (
	"encoding/json"
	"ffCommon/log/log"
	"ffCommon/net/httpclient"
	"ffCommon/util"
	"time"
)

// 携带的数据格式
type customLoginData struct {
	UUIDPlatform string
	UUIDAccount  uint64
	Result       int32 // 0 无错 1 异常返回 2 UUIDPlatform长度不对
}

type httpClient struct {
	client         *httpclient.Client
	chResponseData chan *httpclient.ResponseData

	chExit chan struct{}
}

func (c *httpClient) start() {
	c.chResponseData = make(chan *httpclient.ResponseData, 2)
	c.chExit = make(chan struct{}, 1)

	c.client = httpclient.NewClient(c.chExit, &waitApplicationQuit, chApplicationQuit)

	c.client.Start(2)

	go util.SafeGo(c.mainSendLoop, c.mainSendLoopEnd)

	go util.SafeGo(c.mainRecvLoop, c.mainRecvLoopEnd)

	<-c.chExit
}

func (c *httpClient) mainRecvLoop(params ...interface{}) {
	log.RunLogger.Println("mainRecvLoop")

	{
	deadLoop:
		for {
			select {
			case response := <-c.chResponseData:
				log.RunLogger.Printf("response[%v]", response)
				c.client.BackResponseData(response)
			case <-c.chExit:
				close(c.chExit)
				break deadLoop
			}
		}
	}
}

func (c *httpClient) mainRecvLoopEnd(isPanic bool) {
	log.RunLogger.Println("mainRecvLoopEnd", isPanic)

	close(c.chResponseData)
	c.chResponseData = nil
}

func (c *httpClient) mainSendLoop(params ...interface{}) {
	log.RunLogger.Println("mainSendLoop")

deadLoop:
	for {
		s := ""
		for i := 1; i < 16; i++ {
			select {
			case <-time.After(time.Second):
				break
			case <-c.chExit:
				close(c.chExit)
				break deadLoop
			}

			s += "1"
			data := &customLoginData{
				UUIDPlatform: s,
			}

			bytes, err := json.Marshal(data)
			if err != nil {
				log.RunLogger.Printf("mainSendLoop json.Marshal data[%v] get error[%v]", data, err)
				continue
			}

			request := c.client.ApplyPostRequest("http://127.0.0.1:15011/login", bytes, c.chResponseData)
			log.RunLogger.Printf("mainSendLoop request[%v]", request.UUID())
			c.client.Post(request)
		}
	}
}

func (c *httpClient) mainSendLoopEnd(isPanic bool) {
	log.RunLogger.Println("mainSendLoopEnd", isPanic)
}

func main() {
	defer util.PanicProtect(func(isPanic bool) {
		if isPanic {
			log.RunLogger.Println("异常退出, 以上是错误堆栈")
			<-time.After(time.Hour)
		}
	}, "httpClient")

	c := &httpClient{}
	c.start()
	// tr := &http.Transport{
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	// }
	// client := &http.Client{Transport: tr}
	// resp, err := client.Get("https://localhost:8081")

	// if err != nil {
	// 	fmt.Println("error:", err)
	// 	return
	// }
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))
}
