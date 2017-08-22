package httpclient

import (
	"ffCommon/uuid"
	"fmt"
)

// PostRequest 封装一个请求数据
type PostRequest struct {
	url  string
	data []byte
	uuid uuid.UUID

	chResponseData chan *ResponseData
}

func (request *PostRequest) String() string {
	return fmt.Sprintf("uuid[%v] data[%v]", request.uuid, string(request.data))
}

// UUID 请求唯一标识
func (request *PostRequest) UUID() uuid.UUID {
	return request.uuid
}

//
func (request *PostRequest) onError(err error, response *ResponseData) {
	response.onError(request.uuid, err)
	request.chResponseData <- response
}

//
func (request *PostRequest) onResponse(resp []byte, response *ResponseData) {
	response.onResponse(request.uuid, resp)
	request.chResponseData <- response
}

//
func (request *PostRequest) back() {
	request.data, request.chResponseData = nil, nil
}

//
func (request *PostRequest) init(uuid uuid.UUID, url string, data []byte, chResponseData chan *ResponseData) {
	request.uuid, request.url, request.data, request.chResponseData = uuid, url, data, chResponseData
}

func newPostRequest() *PostRequest {
	return &PostRequest{}
}
