// Package tcpsession 实现了 tcpSession, 使用步骤:
// 0.	进程启动时, 执行本模块的 Init 方法, 以初始化本模块
// 1.	底层连接建立后, 通过本模块的 Apply 方法, 申请 tcpSession 实例
// 2.	调用 tcpSession 实例的 Start 方法, 开始 tcpSession 的收发, 异步
// 3.	tcpSession 实例 向外抛出 NetEventOn 事件, 外界处理到该事件后, 方可通过 tcpSession 实例的 SendProto 方法, 发送 Proto
// 4.	调用 tcpSession 实例的 SendProto 方法, 发送 Proto, 异步
// 5.	tcpSession 实例 向外抛出 NetEventProto 事件, 外界响应处理 Proto
// 6.	调用 tcpSession 实例的 Close 方法, 关闭 tcpSession 实例, 异步
// 7.	tcpSession 实例在底层连接断开时或者网络异常时(接收或发送失败), 进入退出流程
// 7.1	tcpSession 实例 向外抛出 NetEventOff 事件, 外界处理到该事件后, 不可再发送 Proto
// 7.2	tcpSession 实例 等待发送和接收协程退出
// 7.3	tcpSession 实例 向外抛出 NetEventEnd 事件, 外界处理到该事件后, 不可再引用 tcpSession 实例
package tcpsession
