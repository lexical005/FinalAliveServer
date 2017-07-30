// Package tcpclient 实现了 tcpClient, 使用步骤:
// 0.	进程启动时, 执行本模块的 Init 方法, 以初始化本模块
// 1.	使用者 通过 本模块的 NewClient 方法, 创建 tcpClient 实例
// 2.	使用者 调用 tcpClient 实例的 Start 方法, 开始连接服务器, 异步
// 3.	使用者 处理 事件回调
// 4.	使用者 处理到 NetEventOn 事件后, 表明连接已建立, 可发送 Proto
// 5.	使用者 调用 tcpClient 实例的 SendProto 方法, 发送协议, 异步
// 6.	使用者 处理 NetEventProto 事件, 以响应 Proto 处理
// 7.	使用者 处理到 NetEventOff 事件时, 如果需要重连, 则可执行 tcpClient 实例的 ReConnect 方法
// 8.	使用者 调用 tcpClient 实例的 Close 方法, 关闭 tcpClient 实例, 异步
// 9.	使用者 处理到 NetEventEnd 事件时, 表明 tcpClient 关闭完成, 使用者不应再引用该 tcpClient 实例, 该实例也将随着事件处理完成(事件.Back())而回收
package tcpclient
