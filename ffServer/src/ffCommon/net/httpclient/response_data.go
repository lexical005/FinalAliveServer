package httpclient

import "ffCommon/uuid"
import "fmt"

// ResponseData 封装一个结果数据
type ResponseData struct {
	uuid uuid.UUID

	data []byte
	err  error
}

func (response *ResponseData) String() string {
	return fmt.Sprintf("uuid[%v] err[%v] data[%v]", response.uuid, response.err, string(response.data))
}

// UUID 请求唯一标识
func (response *ResponseData) UUID() uuid.UUID {
	return response.uuid
}

// Response 获取结果
func (response *ResponseData) Response() ([]byte, error) {
	return response.data, response.err
}

//
func (response *ResponseData) onError(uuid uuid.UUID, err error) {
	response.uuid, response.err = uuid, err
}

//
func (response *ResponseData) onResponse(uuid uuid.UUID, data []byte) {
	response.uuid, response.data = uuid, data
}

//
func (response *ResponseData) back() {
	response.data = nil
}

func newResponseData() *ResponseData {
	return &ResponseData{}
}
