package httpclient

// Request 一次请求
type Request interface {
	// URL 请求发往的url
	URL() string

	// Unmarshal json格式的字节流反序列化为请求的结果
	Unmarshal(data []byte) error

	// Data 待发送的json格式的字节流
	Data() ([]byte, error)

	// OnError 请求出现错误
	OnError(err error)

	// OnResponse 接收到请求反馈, data为反馈的json格式的字节流
	OnResponse(data []byte)
}
