package ffProto

import (
	"ffCommon/log/log"

	"testing"
)

//                                                 接收到的协议字节流           发送的协议字节流     接收协议重发送
// AgentServer接收到Client的协议                   ExtraDataTypeNormal
// AgentServer将Client的协议返回给Client           ExtraDataTypeNormal         ExtraDataTypeNormal        true
// AgentServer新发协议到Client                                                 ExtraDataTypeNormal        false
// AgentServer将GameServer的协议转发给Client       ExtraDataTypeUUID           ExtraDataTypeNormal        true

// AgentServer接收到GameServer的协议               ExtraDataTypeUUID
// AgentServer将Client的协议转发给GameServer       ExtraDataTypeNormal         ExtraDataTypeUUID          true
// AgentServer新发协议到GameServer                                             ExtraDataTypeUUID          false

// GameServer将收到的Client协议返回给AgentServer   ExtraDataTypeUUID           ExtraDataTypeUUID          true
// GameServer新发协议到GameServer                                              ExtraDataTypeUUID          false

var clientSendHeader *ProtoHeader
var clientSendProto *Proto
var clientSendBytes []byte

var clientRecvHeader *ProtoHeader
var clientRecvHeaderBuff []byte
var clientRecvProto *Proto
var clientRecvBytes []byte

var agentServerSendToClientHeader *ProtoHeader
var agentServerSendToClientProto *Proto
var agentServerSendToClientBytes []byte

var agentServerRecvClientHeader *ProtoHeader
var agentServerRecvClientHeaderBuff []byte
var agentServerRecvClientProto *Proto
var agentServerRecvClientBytes []byte

var gameServerRecvClientHeader *ProtoHeader
var gameServerRecvClientHeaderBuff []byte
var gameServerRecvClientProto *Proto
var gameServerRecvClientBytes []byte

// 测试--Client创建发送到AgentServer的协议
func testClientSendToAgentServerProto(t *testing.T) {
	clientSendHeader = NewProtoHeader()
	clientSendHeader.ResetForSend()

	clientSendProto = ApplyProtoForSend(MessageType_PrepareLoginPlatformUniqueID)

	messageForSend := clientSendProto.Message().(*MsgPrepareLoginPlatformUniqueID)
	messageForSend.SubChannel = "SubChannel"
	messageForSend.UUIDPlatformLogin = "Client->AgentServer"
	messageForSend.Timestamp = 32

	// 进行编码
	err := clientSendProto.Marshal(clientSendHeader)
	if err != nil {
		t.Error(err)
		return
	}

	clientSendBytes = clientSendProto.BytesForSend()

	log.RunLogger.Println("testClientSendToAgentServerProto")
	log.RunLogger.Println(clientSendHeader)
	log.RunLogger.Println(clientSendProto)
}

// 测试--AgentServer接收Client发来的协议
func testAgentServerRecvClientProto(t *testing.T) {
	agentServerRecvClientHeader = NewProtoHeader()
	agentServerRecvClientHeader.ResetForRecv(ExtraDataTypeNormal)

	agentServerRecvClientHeaderBuff = NewProtoHeaderBuf()
	copy(agentServerRecvClientHeaderBuff, clientSendBytes[:cap(agentServerRecvClientHeaderBuff)])

	err := agentServerRecvClientHeader.Unmarshal(agentServerRecvClientHeaderBuff)
	if err != nil {
		t.Error(err)
		return
	}

	agentServerRecvClientProto = ApplyProtoForRecv(agentServerRecvClientHeader)

	copy(agentServerRecvClientProto.BytesForRecv(), clientSendBytes[cap(agentServerRecvClientHeaderBuff):])
	err = agentServerRecvClientProto.OnRecvAllBytes(agentServerRecvClientHeader)
	if err != nil {
		t.Error(err)
		return
	}

	// err = agentServerRecvClientProto.Unmarshal()
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	// agentServerRecvClientBytes = agentServerRecvClientProto.BytesForSend()

	log.RunLogger.Println("testAgentServerRecvClientProto 1")
	log.RunLogger.Println(agentServerRecvClientHeader)
	log.RunLogger.Println(agentServerRecvClientProto)

	// 向GameServer转发
	agentServerRecvClientProto.SetExtraDataUUID(0x12345678)
	sendToGameServerHeader := NewProtoHeader()
	sendToGameServerHeader.ResetForSend()

	err = agentServerRecvClientProto.Marshal(sendToGameServerHeader)
	if err != nil {
		t.Error(err)
		return
	}

	agentServerRecvClientBytes = agentServerRecvClientProto.BytesForSend()

	log.RunLogger.Println("testAgentServerRecvClientProto 2")
	log.RunLogger.Println(agentServerRecvClientHeader)
	log.RunLogger.Println(agentServerRecvClientProto)
	log.RunLogger.Printf("ExtraData:%x\n", agentServerRecvClientProto.ExtraData())
}

// 测试--GameServer接收Client发来的协议
func testGameServerRecvClientProto(t *testing.T) {
	gameServerRecvClientHeader = NewProtoHeader()
	gameServerRecvClientHeader.ResetForRecv(ExtraDataTypeUUID)

	gameServerRecvClientHeaderBuff = NewProtoHeaderBuf()
	copy(gameServerRecvClientHeaderBuff, agentServerRecvClientBytes[:cap(gameServerRecvClientHeaderBuff)])

	err := gameServerRecvClientHeader.Unmarshal(gameServerRecvClientHeaderBuff)
	if err != nil {
		t.Error(err)
		return
	}

	gameServerRecvClientProto = ApplyProtoForRecv(gameServerRecvClientHeader)

	copy(gameServerRecvClientProto.BytesForRecv(), agentServerRecvClientBytes[cap(gameServerRecvClientHeaderBuff):])
	err = gameServerRecvClientProto.OnRecvAllBytes(gameServerRecvClientHeader)
	if err != nil {
		t.Error(err)
		return
	}

	err = gameServerRecvClientProto.Unmarshal()
	if err != nil {
		t.Error(err)
		return
	}

	gameServerRecvClientBytes = gameServerRecvClientProto.BytesForSend()

	log.RunLogger.Println("testGameServerRecvClientProto 1")
	log.RunLogger.Println(gameServerRecvClientHeader)
	log.RunLogger.Println(gameServerRecvClientProto)
	log.RunLogger.Printf("ExtraData:%x\n", gameServerRecvClientProto.ExtraData())

	// 返回给AgentServer
	message, _ := gameServerRecvClientProto.Message().(*MsgPrepareLoginPlatformUniqueID)
	GameUserid := "GameServer->Client"
	message.UUIDPlatformLogin = GameUserid

	gameServerRecvClientProto.SetExtraDataUUID(gameServerRecvClientProto.ExtraData())

	sendToAgentServerHeader := NewProtoHeader()
	sendToAgentServerHeader.ResetForSend()

	err = gameServerRecvClientProto.Marshal(sendToAgentServerHeader)
	if err != nil {
		t.Error(err)
		return
	}

	gameServerRecvClientBytes = gameServerRecvClientProto.BytesForSend()

	log.RunLogger.Println("testGameServerRecvClientProto 2")
	log.RunLogger.Println(gameServerRecvClientHeader)
	log.RunLogger.Println(gameServerRecvClientProto)
	log.RunLogger.Printf("ExtraData:%x\n", gameServerRecvClientProto.ExtraData())
}

// 测试--AgentServer接收GameServer返回的Client发来的协议
func testAgentServerRecvGameServerProto(t *testing.T) {
	agentServerRecvClientHeader = NewProtoHeader()
	agentServerRecvClientHeader.ResetForRecv(ExtraDataTypeUUID)

	agentServerRecvClientHeaderBuff = NewProtoHeaderBuf()
	copy(agentServerRecvClientHeaderBuff, gameServerRecvClientBytes[:cap(agentServerRecvClientHeaderBuff)])

	err := agentServerRecvClientHeader.Unmarshal(agentServerRecvClientHeaderBuff)
	if err != nil {
		t.Error(err)
		return
	}

	agentServerRecvClientProto = ApplyProtoForRecv(agentServerRecvClientHeader)

	copy(agentServerRecvClientProto.BytesForRecv(), gameServerRecvClientBytes[cap(agentServerRecvClientHeaderBuff):])
	err = agentServerRecvClientProto.OnRecvAllBytes(agentServerRecvClientHeader)
	if err != nil {
		t.Error(err)
		return
	}

	// err = agentServerRecvClientProto.Unmarshal()
	// if err != nil {
	// 	t.Error(err)
	// 	return
	// }

	agentServerRecvClientBytes = agentServerRecvClientProto.BytesForSend()

	log.RunLogger.Println("testAgentServerRecvGameServerProto")
	log.RunLogger.Println(agentServerRecvClientHeader)
	log.RunLogger.Println(agentServerRecvClientProto)
	log.RunLogger.Printf("ExtraData:%x\n", agentServerRecvClientProto.ExtraData())

	// 返回给Client
	agentServerRecvClientProto.SetExtraDataNormal()

	sendToClientHeader := NewProtoHeader()
	sendToClientHeader.ResetForSend()

	err = agentServerRecvClientProto.Marshal(sendToClientHeader)
	if err != nil {
		t.Error(err)
		return
	}

	agentServerRecvClientBytes = agentServerRecvClientProto.BytesForSend()

	log.RunLogger.Println("testAgentServerRecvGameServerProto 2")
	log.RunLogger.Println(agentServerRecvClientHeader)
	log.RunLogger.Println(agentServerRecvClientProto)
	log.RunLogger.Printf("ExtraData:%x\n", agentServerRecvClientProto.ExtraData())
}

// 测试--Client接收AgentServer返回的Client上传的协议
func testClientRecvClientProto(t *testing.T) {
	clientRecvHeader = NewProtoHeader()
	clientRecvHeader.ResetForRecv(ExtraDataTypeNormal)

	clientRecvHeaderBuff = NewProtoHeaderBuf()
	copy(clientRecvHeaderBuff, agentServerRecvClientBytes[:cap(clientRecvHeaderBuff)])

	err := clientRecvHeader.Unmarshal(clientRecvHeaderBuff)
	if err != nil {
		t.Error(err)
		return
	}

	clientRecvProto = ApplyProtoForRecv(clientRecvHeader)

	copy(clientRecvProto.BytesForRecv(), agentServerRecvClientBytes[cap(clientRecvHeaderBuff):])
	err = clientRecvProto.OnRecvAllBytes(clientRecvHeader)
	if err != nil {
		t.Error(err)
		return
	}

	err = clientRecvProto.Unmarshal()
	if err != nil {
		t.Error(err)
		return
	}

	clientRecvBytes = clientRecvProto.BytesForSend()

	log.RunLogger.Println("testClientRecvClientProto")
	log.RunLogger.Println(clientRecvHeader)
	log.RunLogger.Println(clientRecvProto)
}

// 测试--AgentServer创建发送到Client的协议
func testAgentServerSendToClientProto(t *testing.T) {
	agentServerSendToClientHeader = NewProtoHeader()
	agentServerSendToClientHeader.ResetForSend()

	agentServerSendToClientProto = ApplyProtoForSend(MessageType_PrepareLoginPlatformUniqueID)

	messageForSend := agentServerSendToClientProto.Message().(*MsgPrepareLoginPlatformUniqueID)
	messageForSend.SubChannel = "SubChannel"
	messageForSend.UUIDPlatformLogin = "AgentServer->Client"
	messageForSend.Timestamp = 32

	// 进行编码
	err := agentServerSendToClientProto.Marshal(agentServerSendToClientHeader)
	if err != nil {
		t.Error(err)
		return
	}

	agentServerSendToClientBytes = agentServerSendToClientProto.BytesForSend()

	log.RunLogger.Println("testAgentServerSendToClientProto")
	log.RunLogger.Println(clientSendHeader)
	log.RunLogger.Println(clientSendProto)
}

// 测试--Client接收AgentServer发来的协议
func testClientRecvAgentServerProto(t *testing.T) {
	clientRecvHeader = NewProtoHeader()
	clientRecvHeader.ResetForRecv(ExtraDataTypeNormal)

	clientRecvHeaderBuff = NewProtoHeaderBuf()
	copy(clientRecvHeaderBuff, agentServerSendToClientBytes[:cap(clientRecvHeaderBuff)])

	err := clientRecvHeader.Unmarshal(clientRecvHeaderBuff)
	if err != nil {
		t.Error(err)
		return
	}

	clientRecvProto = ApplyProtoForRecv(clientRecvHeader)

	copy(clientRecvProto.BytesForRecv(), agentServerSendToClientBytes[cap(clientRecvHeaderBuff):])
	err = clientRecvProto.OnRecvAllBytes(clientRecvHeader)
	if err != nil {
		t.Error(err)
		return
	}

	err = clientRecvProto.Unmarshal()
	if err != nil {
		t.Error(err)
		return
	}

	clientRecvBytes = clientRecvProto.BytesForSend()

	log.RunLogger.Println("testClientRecvAgentServerProto")
	log.RunLogger.Println(agentServerRecvClientHeader)
	log.RunLogger.Println(agentServerRecvClientProto)
}

func Test_RecvClientProto(t *testing.T) {

	{
		log.RunLogger.Println("Client发送->AgentServer接收->转到GameServer->AgentServer接收结果->AgentServer返回给Client")
		testClientSendToAgentServerProto(t)
		log.RunLogger.Println("")

		testAgentServerRecvClientProto(t)
		log.RunLogger.Println("")

		testGameServerRecvClientProto(t)
		log.RunLogger.Println("")

		testAgentServerRecvGameServerProto(t)
		log.RunLogger.Println("")

		testClientRecvClientProto(t)
		log.RunLogger.Println("Client发送->AgentServer接收->转到GameServer->AgentServer接收结果->AgentServer返回给Client")
		log.RunLogger.Println("")
		log.RunLogger.Println("")
	}

	// {
	// 	log.RunLogger.Println("AgentServer发送->Client接收")
	// 	testAgentServerSendToClientProto(t)
	// 	log.RunLogger.Println("")

	// 	testClientRecvAgentServerProto(t)
	// 	log.RunLogger.Println("AgentServer发送->Client接收")
	// 	log.RunLogger.Println("")
	// 	log.RunLogger.Println("")
	// }
}
