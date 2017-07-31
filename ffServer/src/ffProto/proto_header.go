package ffProto

import "fmt"

// ProtoHeader 协议头
type ProtoHeader struct {
	protoID       int32 // 协议号
	contentLength int   // 协议体的长度

	recvExtraDataType ExtraDataType // 用于接收时, 限定的附加数据的类型

	lastRecvExtraDataNormal []byte // 上一次接收到的附加数据
	lastSendExtraDataNormal []byte // 上一次发送的附加数据
}

// Unmarshal 反序列化协议头(bigEndian)
// 将底层收到的字节流, 记录到协议头内
func (ph *ProtoHeader) Unmarshal(buf []byte) error {
	// 附加数据类型限定
	if ph.recvExtraDataType.Value() != buf[protoHeaderExtraDataTypeOffset] {
		return fmt.Errorf("ffProto.ProtoHeader.Unmarshal: ExtraDataType not match[%v(%v):%v]",
			ph.recvExtraDataType, ph.recvExtraDataType.Value(), buf[protoHeaderExtraDataTypeOffset])
	}

	// 协议号有效性
	ph.protoID = int32(uint16(buf[protoHeaderIDOffset+1]) | uint16(buf[protoHeaderIDOffset+0])<<8)
	if _, ok := MessageType_name[ph.protoID]; !ok {
		return fmt.Errorf("ffProto.ProtoHeader.Unmarshal: invalid protoID[%v]", ph.protoID)
	}

	// 协议内容长度有效性
	ph.contentLength = int(uint16(buf[protoHeaderContentOffset+1]) | uint16(buf[protoHeaderContentOffset+0])<<8)
	if ph.contentLength > protoMaxContentLength {
		return fmt.Errorf("ffProto.ProtoHeader.Unmarshal: ContentLength[%v] > protoMaxContentLength[%v]", ph.contentLength, protoMaxContentLength)
	}

	// todo:协议头校验

	// todo:流量限制

	return nil
}

// ResetForRecv 重置, 以便再次使用
func (ph *ProtoHeader) ResetForRecv(recvExtraDataType ExtraDataType) {
	ph.recvExtraDataType = recvExtraDataType

	ph.lastRecvExtraDataNormal[0] = 0xFF
	ph.lastSendExtraDataNormal[0] = 0xFF
}

// ResetForSend 重置, 以便再次使用
func (ph *ProtoHeader) ResetForSend() {
	ph.lastRecvExtraDataNormal[0] = 0xFF
	ph.lastSendExtraDataNormal[0] = 0xFF
}

// marshalHeader 序列化协议头(bigEndian)
// 将协议头内的数据, 序列化到待发送的字节流
func (ph *ProtoHeader) marshalHeader(buf []byte, marshalExtraDataType bool) {
	// 协议号
	buf[protoHeaderIDOffset] = byte(ph.protoID >> 8)
	buf[protoHeaderIDOffset+1] = byte(ph.protoID)

	// 协议内容长度
	buf[protoHeaderContentOffset] = byte(ph.contentLength >> 8)
	buf[protoHeaderContentOffset+1] = byte(ph.contentLength)

	// 附加数据类型
	if marshalExtraDataType {
		buf[protoHeaderExtraDataTypeOffset] = ph.recvExtraDataType.Value()
	}
}

// 将lastSendExtraDataNormal序列化到协议缓冲区
func (ph *ProtoHeader) marshalSendExtraDataNormal(buf []byte) {
	ph.lastSendExtraDataNormal[0]++
	copy(buf, ph.lastSendExtraDataNormal)
}

// String 返回ProtoHeader的自我描述
func (ph *ProtoHeader) String() string {
	return fmt.Sprintf("protoID[%v:%v] contentLength[%v] recvExtraDataType[%b] lastRecvExtraDataNormal[%v] lastSendExtraDataNormal[%v]",
		ph.protoID, MessageType_name[int32(ph.protoID)], ph.contentLength, ph.recvExtraDataType, ph.lastRecvExtraDataNormal, ph.lastSendExtraDataNormal)
}

// NewProtoHeader 新建一个协议头
func NewProtoHeader() *ProtoHeader {
	header := &ProtoHeader{
		lastRecvExtraDataNormal: make([]byte, extraDataNormalLength, extraDataNormalLength),
		lastSendExtraDataNormal: make([]byte, extraDataNormalLength, extraDataNormalLength),
	}
	header.lastRecvExtraDataNormal[0] = 0
	header.lastSendExtraDataNormal[0] = 0
	return header
}

// NewProtoHeaderBuf 新建一个协议头缓冲区
func NewProtoHeaderBuf() []byte {
	return make([]byte, protoHeaderLength)
}
