// Package tcpclient 实现了 tcpClient, 使用步骤:
// tcpClient: 对象通过管道转交给其他goroutine后, 不再对其做任何操作, 不与其他goroutine同时共享同一个对象
// 	外界执行Start, 开启连接到指定Server的goroutine
// 		连接失败
// 			延迟一秒
// 			检查chNtfWorkExit是否关闭了, 如果是, 则退出goroutine
// 		连接成功
// 			通过chNewSession向外界通知新连接Session建立, Client不维护Session
// 			等待chNtfWorkExit关闭或者chReConnect重连
// 				chNtfWorkExit关闭: 外界调用了Close, 关停Client
// 				chReConnect重连:  上一连接关闭了, 继续重连
// 	外界执行Stop, 开始关停Client
// 		连接到指定Server的goroutine, 总是会退出
// 			阅读上一节Start, 即可知晓缘由
// 		连接到指定Server的goroutine退出时
// 			通过chClientClosed像外界通知通知Client关闭完成, 可执行清理操作Back
// 	外界执行Back, 回收Client资源
// 		只应在外界通过chServerClose接收到可回收事件之后下执行
package tcpclient
