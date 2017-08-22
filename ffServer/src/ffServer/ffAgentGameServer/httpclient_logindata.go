package main

import (
	"encoding/json"
	"ffCommon/uuid"
)

// 携带的数据格式
type httpClientCustomLoginData struct {
	// 请求数据/反馈数据
	UUIDPlatform string
	UUIDAccount  uint64
	Result       int32 // 0 无错 1 异常返回 2 UUIDPlatform长度不对

	// 请求者
	uuidRequester uuid.UUID

	// 记录请求过程中, 出现的错误
	err error
}

// Data 待发送的json格式的字节流
func (logindata *httpClientCustomLoginData) Data() ([]byte, error) {
	return json.Marshal(logindata)
}

// OnError 请求出现错误
func (logindata *httpClientCustomLoginData) OnError(err error) {
	logindata.err = err

	instAgentUserServer.OnCustomLoginResult(logindata)
}

// OnResponse 接收到请求反馈, data为反馈的json格式的字节流
func (logindata *httpClientCustomLoginData) OnResponse(data []byte) {
	logindata.err = json.Unmarshal(data, logindata)

	instAgentUserServer.OnCustomLoginResult(logindata)
}
