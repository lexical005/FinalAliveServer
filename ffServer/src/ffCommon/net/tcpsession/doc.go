// Package tcpsession 实现了 tcpSession, 使用步骤:
// tcpSession
// 	外界通过tcpsession模块Apply方法, 申请Session对象
// 		 必须提供net.Conn, 在返回Session时, 内部已进行相关设置
// 	外界执行Start
// 		开启接收Proto的goroutine
// 			接收协议头
// 				如果出现网络错误, 则退出接收goroutine
// 			接收协议体
// 				申请Proto
// 				如果接收协议体字节流出错或者发序列化时出错, 则回收Proto, 且返回错误, 否则, 将Proto封装到网络事件--Proto
// 				返回网络事件-Proto
// 			如果接收协议体时出现错误, 则退出接受goroutine, 否则, 将网络事件--Proto, 通过chNetEventData管道, 向外界通知
// 			检查Session是否开始退出
// 				如果是, 则退出发送goroutine
// 		开启发送Proto的goroutine
// 			循环从chSendProto内取出待发送的Proto
// 				检查Session是否开始退出
// 					如果是, 则退出发送goroutine
// 				如果取到了nil, 则认为外界要求关闭Session, 设置主动关闭标志, 退出发送goroutine
// 				如果取到了有效Proto, 则尝试发送
// 					在发送操作后, 无论成功或者失败, 发送的Proto都会执行尝试回收BackAfterSend
// 					如果发送时出现网络错误, 则退出发送goroutine
// 					循环检查是否未完全发送
// 						等待2毫秒
// 						发送剩余内容
// 							如果成功, 则跳出循环发送
// 							检查Session是否开始退出
// 								如果是, 则退出发送goroutine
// 			在发送goroutine退出时
// 				设置发送goroutine已完成退出
// 				进入doClose逻辑
// 	外界执行Close
// 		外界决定关闭Client或Server后, 通过此方法, 关闭chNewSession内尚未处理的新建立的连接. 必须发生在Start之前
// 		直接关闭底层net.Conn
// 		将自身Session归还到sessionPool
// 	内部执行doClose
// 		设置关闭管道, 开始退出
// 			发送goroutine和接收goroutine会检查此管道, 如果此处于管道关闭状态, 则会退出goroutine
// 		关闭底层连接net.Conn
// 			会导致在此net.Conn进行的发送和读取操作失败, 进而导致接收Proto和发送Proto退出
// 		等待接受goroutine和发送goroutine退出完成
// 		通过chNetEventData管道, 向外界通知网络事件--Session断开事件
// 			外界处理完毕网络事件后, 必须执行回收操作
// 			在网络事件--Session断开事件回收处理时, 会执行Session的回收
package tcpsession
