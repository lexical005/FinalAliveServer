package httpclient

import (
	"net/http"
)

const (
	// postContentTypeJSON Post发送的数据的类型
	postContentTypeJSON = "application/json;charset=utf-8"
)

// NewClient 返回一个Client, 供外界向指定url通讯使用
func NewClient(chExit chan struct{}, chRequest chan Request) *Client {
	return &Client{
		client: &http.Client{},

		chExit:    chExit,
		chRequest: chRequest,
	}
}
