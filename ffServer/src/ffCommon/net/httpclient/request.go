package httpclient

// Request 一次请求
type Request interface {
	// Data 待发送的json格式的字节流
	Data() ([]byte, error)

	// OnError 请求过程中发生错误
	OnError(err error)

	// OnResponse 接收到请求反馈, data为反馈的json格式的字节流
	OnResponse(data []byte)
}
