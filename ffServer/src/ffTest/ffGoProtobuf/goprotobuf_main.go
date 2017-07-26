package main

import (
	"ffCommon/log/log"

	// 辅助库
	"github.com/golang/protobuf/proto"
	// 协议路径
	"ffProto"
)

func testExtraDataUUID() {
	sendHeader := ffProto.NewProtoHeader()

	sendProto := ffProto.ApplyProtoForSend(ffProto.MessageType_MT_MsgPrepareLoginPlatformUniqueID)
	messageForSend := sendProto.Message().(*ffProto.MsgPrepareLoginPlatformUniqueID)
	messageForSend.SubChannel = proto.String("SubChannel")
	messageForSend.GameUserid = proto.String("GameUserid")
	messageForSend.Timestamp = proto.Int32(32)

	sendProto.SetExtraData(false, ffProto.ExtraDataTypeUUID, 0x12345678)

	// 进行编码
	err := sendProto.Marshal(sendHeader)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	log.RunLogger.Printf("%x\n", sendProto.ExtraData())

	sendBytes := sendProto.BytesForSend()

	log.RunLogger.Println(sendHeader)
	log.RunLogger.Println(sendProto)
	log.RunLogger.Println(messageForSend)
	log.RunLogger.Println("")

	// 进行解码
	recvHeader := ffProto.NewProtoHeader()
	recvHeaderBuff := ffProto.NewProtoHeaderBuf()
	copy(recvHeaderBuff, sendBytes[:cap(recvHeaderBuff)])

	err = recvHeader.Unmarshal(recvHeaderBuff)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	log.RunLogger.Println(recvHeaderBuff)
	log.RunLogger.Println(recvHeader)
	log.RunLogger.Println("")

	recvProto := ffProto.ApplyProtoForRecv(recvHeader)

	copy(recvProto.BytesForRecv(), sendBytes[cap(recvHeaderBuff):])
	recvProto.OnRecvAllBytes(recvHeader)

	messageForRecv, err := recvProto.Unmarshal()
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	log.RunLogger.Printf("%x\n", recvProto.ExtraData())

	log.RunLogger.Println(recvHeader)
	log.RunLogger.Println(recvProto)
	log.RunLogger.Println(messageForRecv)
}

func testExtraDataNormal() {
	sendHeader := ffProto.NewProtoHeader()

	sendProto := ffProto.ApplyProtoForSend(ffProto.MessageType_MT_MsgPrepareLoginPlatformUniqueID)
	messageForSend := sendProto.Message().(*ffProto.MsgPrepareLoginPlatformUniqueID)
	messageForSend.SubChannel = proto.String("SubChannel")
	messageForSend.GameUserid = proto.String("GameUserid")
	messageForSend.Timestamp = proto.Int32(32)

	// 进行编码
	err := sendProto.Marshal(sendHeader)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	sendBytes := sendProto.BytesForSend()

	log.RunLogger.Println(sendHeader)
	log.RunLogger.Println(sendProto)
	log.RunLogger.Printf("%x\n", sendProto.ExtraData())
	log.RunLogger.Println(messageForSend)
	log.RunLogger.Println("")

	// ---------------------------------------
	// 进行接收
	recvHeader := ffProto.NewProtoHeader()
	recvHeaderBuff := ffProto.NewProtoHeaderBuf()
	copy(recvHeaderBuff, sendBytes[:cap(recvHeaderBuff)])

	err = recvHeader.Unmarshal(recvHeaderBuff)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	recvProto := ffProto.ApplyProtoForRecv(recvHeader)

	copy(recvProto.BytesForRecv(), sendBytes[cap(recvHeaderBuff):])
	recvProto.OnRecvAllBytes(recvHeader)

	messageForRecv, err := recvProto.Unmarshal()
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	log.RunLogger.Println(recvHeader)
	log.RunLogger.Println(recvProto)
	log.RunLogger.Printf("%x\n", recvProto.ExtraData())
	log.RunLogger.Println(messageForRecv)
	log.RunLogger.Println("")
}

func testExtraDataNormalToUUID() {
	sendHeader := ffProto.NewProtoHeader()

	sendProto := ffProto.ApplyProtoForSend(ffProto.MessageType_MT_MsgPrepareLoginPlatformUniqueID)
	messageForSend := sendProto.Message().(*ffProto.MsgPrepareLoginPlatformUniqueID)
	messageForSend.SubChannel = proto.String("SubChannel")
	messageForSend.GameUserid = proto.String("GameUserid")
	messageForSend.Timestamp = proto.Int32(32)

	// 进行编码
	err := sendProto.Marshal(sendHeader)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	sendBytes := sendProto.BytesForSend()

	log.RunLogger.Println(sendHeader)
	log.RunLogger.Println(sendProto)
	log.RunLogger.Println(messageForSend)
	log.RunLogger.Println("")

	// ---------------------------------------
	// 进行接收
	recvHeader := ffProto.NewProtoHeader()
	recvHeaderBuff := ffProto.NewProtoHeaderBuf()
	copy(recvHeaderBuff, sendBytes[:cap(recvHeaderBuff)])

	err = recvHeader.Unmarshal(recvHeaderBuff)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	recvProto := ffProto.ApplyProtoForRecv(recvHeader)

	copy(recvProto.BytesForRecv(), sendBytes[cap(recvHeaderBuff):])
	recvProto.OnRecvAllBytes(recvHeader)

	log.RunLogger.Println(recvHeader)
	log.RunLogger.Println(recvProto)
	log.RunLogger.Println("")

	// ---------------------------------------
	// 进行转发
	recvProto.SetExtraData(true, ffProto.ExtraDataTypeUUID, 0x12345678)

	err = recvProto.Marshal(recvHeader)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	sendBytes = recvProto.BytesForSend()

	log.RunLogger.Println(recvHeader)
	log.RunLogger.Println(recvProto)
	log.RunLogger.Printf("%x\n", recvProto.ExtraData())
	log.RunLogger.Println("")

	// ---------------------------------------
	// 接收转发
	recvHeader2 := ffProto.NewProtoHeader()
	recvHeader2Buff := ffProto.NewProtoHeaderBuf()
	copy(recvHeader2Buff, sendBytes[:cap(recvHeader2Buff)])

	err = recvHeader2.Unmarshal(recvHeader2Buff)
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	recvProto2 := ffProto.ApplyProtoForRecv(recvHeader2)

	copy(recvProto2.BytesForRecv(), sendBytes[cap(recvHeader2Buff):])
	recvProto2.OnRecvAllBytes(recvHeader2)

	messageForRecv, err := recvProto2.Unmarshal()
	if err != nil {
		log.RunLogger.Println(err)
		return
	}

	log.RunLogger.Println(recvHeader2)
	log.RunLogger.Println(recvProto2)
	log.RunLogger.Printf("%x\n", recvProto2.ExtraData())
	log.RunLogger.Println(messageForRecv)
	log.RunLogger.Println("")
}

func main() {
	testExtraDataNormal()
	log.RunLogger.Println("---------------")
	testExtraDataUUID()
	log.RunLogger.Println("---------------")
	testExtraDataNormalToUUID()
}
