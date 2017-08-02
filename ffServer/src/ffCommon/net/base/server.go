package base

// Server 在指定配置上监听用户连接, 有新连接建立时, 通过管道, 向外界汇报
//	关闭服务器流程:
//		先执行StopAccept, 让Server结束监听
//		外界关闭现有所有连接
//		当所有连接都已关闭后, 执行Back, 回收Server
type Server interface {
	// Start 启动Server, 开始监听用户连接, 一次性, 同步
	//	chNewSession: 外界接收新连接被创建事件的管道, Server仅有写入权
	//	chServerClosed: 用于向外界通知关闭完成的管道, Server仅有写入权
	Start(chNewSession chan Session, chServerClosed chan struct{}) error

	// StopAccept 停止接受连接请求, 只应在Start成功前提下希望关闭服务器时执行, 一次性, 同步
	StopAccept()

	// Back 回收服务器资源, 只应在Start失败或者所有连接均已完成关闭情况下执行, 一次性, 同步
	Back()

	String() string
}
