protoc.exe --go_out=. ffProto.proto

http://blog.csdn.net/u011518120/article/details/54604615

protobuf v3:
1. 移除了 required和optional 关键字
2. 移除了 默认值, 改为字段类型默认零值
3. 移除了 unknown field
4. 支持 map
	map<string, string> fieldNames = 1; map 是 unordered_map
5. 注释 //
6. 基本数据类型
	
7. syntax syntax = "proto3";
8. tag1-15, 字段唯一数字标识, 1-15使用一个字节编码
9. 保留标识符 reserved

todo:
    客户端伪造附加数据控制位
    服务端发送失败, 协议重发

附加数据大小计算方式:
    extraDataUUID:      服务端与服务端之间, 附加数据为UUID, extraDataUUIDLength字节
    extraDataNormal:    服务端与客户端之间, 附加数据为自增序列号, extraDataNormalLength字节

附加数据机制:
接收:
    1. 接收协议头的字节流
    2. 校验协议头的字节流, 校验无误后, 转换为协议头(附加数据为extraDataNormal时, 需做额外校验)
    3. 根据协议头申请协议
    4. 申请协议时, 需求的缓冲区大小计算公式为: 协议头大小+协议体大小+附加数据大小
    5. 根据协议头, 初始化协议, 将协议头的内容, 记录到协议缓冲区
    6. 将本次通讯的剩余内容, 接收到协议缓冲区上, 协议缓冲区由三部分数据依次组成, 协议头, 协议体, 附加数据
    7. 上层逻辑, 则合适的时机, 将协议体字节流, 反序列化为协议结构
    8. 附加数据为extraDataUUID时, 将协议缓冲区内的附加数据字节流, 记录到独立的附加数据字节流内extraData
    9. 上层逻辑, 至此取得了具体的协议结构, 以及附加数据


服务端向客户端发送数据:
    1. 根据协议号申请协议, 协议缓冲区尽可能的大
    2. 填充协议结构体
    3. 将协议序列化为字节流
    4. 将协议结构体序列化为字节流
    5. 如果协议附加数据模式为extraDataNormal, 附加数据长度为extraDataNormalLength
    6. 如果协议附加数据模式为extraDataUUID, 附加数据长度为extraDataUUIDLength
    7. 需要时, 调整协议缓冲区, 以满足大小要求: 协议头大小+协议体大小+附加数据大小
    8. 将协议头序列化到协议缓冲区的协议头部分
    9. 如果协议附加数据模式为extraDataNormal, 则拷贝独立存储的附加数据, 到协议缓冲区的附加数据部分
    10.如果协议附加数据模式为extraDataUUID, 则拷贝协议头的附加数据, 到协议缓冲区的附加数据部分


                                                接收到的协议字节流           发送的协议字节流
AgentServer接收到Client的协议                   ExtraDataTypeNormal
AgentServer将Client的协议返回给Client           ExtraDataTypeNormal         ExtraDataTypeNormal
AgentServer新发协议到Client                                                 ExtraDataTypeNormal
AgentServer将GameServer的协议转发给Client       ExtraDataTypeUUID           ExtraDataTypeNormal

AgentServer接收到GameServer的协议               ExtraDataTypeUUID
AgentServer将Client的协议转发给GameServer       ExtraDataTypeNormal         ExtraDataTypeUUID
AgentServer新发协议到GameServer                                             ExtraDataTypeUUID

GameServer将收到的Client协议返回给AgentServer   ExtraDataTypeUUID           ExtraDataTypeUUID
GameServer新发协议到GameServer                                              ExtraDataTypeUUID


为发送而申请的协议, 其附加数据类型, 默认是ExtraDataTypeNormal, 可通过ChangeExtraData改变

一旦修改了协议数据, 就必须重新序列化		ok

GameServer新建协议发送到客户端

SetExtraData    发送任何协议前, 必须调用SetExtraData

协议状态 useState
    申请协议用于发送   resetForSend             设置  useState = useStateSend
    申请协议用于接收   resetForRecv             设置  useState = useStateRecv
    发送协议前        SetExtraData             设置  useState = useStateSend
    手动缓存等待分发   SetCacheWaitDispatch     设置  useState = useStateCacheWaitDispatch
    手动缓存等待发送   SetCacheWaitSend         设置  useState = useStateCacheWaitSend

    发送完毕尝试回收 BackAfterSend
    分发完毕尝试回收 BackAfterDispatchCache

    useStateSend, 协议在被发送时, 总是会被设置成此状态; 当底层完成发送操作时, 协议将被回收：
        新申请协议用于发送
        接收到的协议被缓存分发, 被逻辑处理后立即被返回或被转发
        接收到的协议被缓存分发, 被逻辑处理后被缓存等待发送, 查询到结果后被返回或被转发

        底层发送操作完成时, 协议将被回收

    useStateRecv, 接收到的协议会被设置成此状态; 临时状态, 接下立即进入状态useStateCacheWaitDispatch
        接收到的协议

    useStateCacheWaitDispatch, 接收到的协议进入等待分发, 由逻辑处理时, 会设置成此状态; 逻辑不更改状态时（发送或缓存）, 将被立即回收：
        接收到的协议, 进入等待分发给逻辑处理

        被分发给逻辑处理完毕后依然处于此状态的（未转发且未缓存等待发送）, 协议将被回收

    useStateCacheWaitSend, 逻辑处理协议时需要等待异步查询结果才能返回或转发时, 会设置成此状态：
        逻辑处理过程中, 需要等待异步查询结果然后再返回或转发

        没有对应的回收, 最终必须进入发送状态 useStateSend
