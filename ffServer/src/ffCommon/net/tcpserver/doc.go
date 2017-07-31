// Package tcpserver 实现了 tcpServer, 使用步骤:
// 0.	进程启动时, 执行本模块的 Init 方法, 以初始化本模块
// 1.	执行本模块的 NewServer 方法, 创建 tcpServer 实例
// 2.	调用 tcpServer 实例的 Start 方法, 开始 tcpServer 的监听和收发, 异步
// 3.	tcpServer 实例在监听到新连接建立后, 将session通过新session管道, 发送到session协程维护
// 3.1	立即将新建立的session维护到session列表内
// 4.	tcpServer 实例 向外抛出 NetEventOn 事件, 外界处理到该事件后, 即认为新连接建立(携带指定连接的UUID), 外界可透过此连接收发 Proto
// 5.	调用 tcpServer 实例的 SendProto 方法(携带指定连接的UUID), 发送 Proto, 异步
// 6.	tcpServer 实例 向外抛出 NetEventProto 事件(携带指定连接的UUID), 外界响应处理 Proto
// 7.	调用 tcpServer 实例的 CloseSession 方法, 结束指定session, 异步
// 7.1	外界处理到 NetEventDateOff 事件(携带指定连接的UUID), 即可认为指定session连接断开了(携带是否主动关闭参数)
// 7.2	外界处理完毕 NetEventDateOff 事件时, tcpServer 实例 执行 onSessionClosed, 并从session列表内移除
// 8.	调用 tcpServer 实例的 Close 方法, 关闭 tcpServer 实例, 异步
// 8.1	tcpServer 实例 结束监听, 不再建立新的连接
// 8.2	tcpServer 实例 监听协程结束后, 关闭新session管道
// 8.3	tcpServer 实例 session处理协程处理到session管道关闭时, 结束主循环, 进入关闭处理
// 8.4	tcpServer 实例 向所有记录的session 通知 关闭, 如果没有session已建立, 则直接抛出所有session均已断开事件
// 8.5	tcpServer 实例 继续处理session事件 同时等待 所有session 均已断开事件
// 8.6	tcpServer 实例 在 onSessionClosed 处理过程中, 如果所有session均已断开, 则抛出所有session均已断开事件
// 8.7	tcpServer 实例 处理到 所有session均已断开事件, 则退出session协程
// 9	tcpSession 实例 向外抛出 NetEventEnd 事件, 外界处理到该事件后, 不可再引用 tcpServer 实例
// 9.1	外界处理完毕 NetEventEnd 事件后, 执行 tcpSession 实例的 back 方法, tcpServer 实例在此被回收
package tcpserver
