// Package tcpserver 实现了 tcpServer, 使用步骤:
// tcpServer: 对象通过管道转交给其他goroutine后, 不再对其做任何操作, 不与其他goroutine同时共享同一个对象
// 	外界执行Start, 开启接受连接请求goroutine
// 		listener.Accept接受连接请求后, 通过chNewSession向外界通知新连接Session建立, Server不维护Session
// 	外界执行StopAccept, 开始关停Server, 只应在Start成功之后希望关停服务器时执行
// 		标记退出, 关闭listener对象
// 		接受连接请求goroutine内, listener.Accept失败, 检查退出标记, 退出goroutine
// 		接受连接请求goroutine退出时
// 			通过chServerClose像外界通知通知Server关闭完成, 可执行清理操作Back
// 外界执行Back, 回收Server资源
//		只应在Start失败或者外界通过chServerClose接收到可回收事件之后下执行
package tcpserver
