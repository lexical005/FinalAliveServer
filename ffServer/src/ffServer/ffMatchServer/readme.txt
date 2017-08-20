代理服务器
直接与用户交互的服务器，后端服务器与用户之间交流的中转站

Client1->AgentServer->LoginServer1
Client2->AgentServer->LoginServer1

Client1->AgentServer->GameServer1
Client2->AgentServer->GameServer1

GameServer1->AgentServer->Client1
GameServer1->AgentServer->Client2

Client连接建立
Client请求登录校验
    将Client与鉴权服务器绑定, 转发协议
Client请求进入游戏
    将Client与游戏服务器绑定, 转发协议

目前只实现一端监听客户端(允许N个同时连接), 一端监听特定类型服务器(允许N个同时连接)

客户端断开连接, 要告知服务端侧

测试延迟关闭Session       --pass

DisConnect事件            --pass

listen client 时, extra data type must be normal

当与GameServer的连接断开时, 进入到该GameServer的Client, 不会立即被踢出, 直到收到其下一条需转发到该GameServer的协议时, 那时才会被踢出
当与GameServer的连接断开时, 会立即从tunnelServerAgentManager.agents内移除该GameServer
