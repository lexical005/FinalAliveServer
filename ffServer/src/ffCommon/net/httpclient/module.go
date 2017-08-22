package httpclient

import (
	"ffCommon/log/log"
	"ffCommon/uuid"
	"net/http"
)

const (
	// postContentTypeJSON Post发送的数据的类型
	postContentTypeJSON = "application/json;charset=utf-8"
)

var uuidGenerator uuid.Generator

func init() {
	var err error
	uuidGenerator, err = uuid.NewGeneratorSafe(0)
	if err != nil {
		log.FatalLogger.Printf("httpclient uuid.NewGeneratorSafe failed: %v", err)
	}
}

// NewClient 返回一个Client, 供外界向指定url通讯使用
func NewClient(
	chExit chan struct{},
	countApplicationQuit *int32,
	chApplicationQuit chan struct{}) *Client {
	return &Client{
		client: &http.Client{},

		chExit:               chExit,
		countApplicationQuit: countApplicationQuit,
		chApplicationQuit:    chApplicationQuit,
	}
}
