package ffProto

import (
	"ffCommon/log/log"
	"fmt"

	"github.com/golang/protobuf/proto"
)

// Proto proto
type Proto struct {
	pb *proto.Buffer //

	msg            proto.Message
	msgNeedMarshal bool // 协议内容需要序列化

	protoID    MessageType // 协议id
	useState   useState    // 协议当前状态
	limitState limitState  // 协议被限定使用范围
	buf        []byte      // 协议使用的缓冲区

	extraDataType ExtraDataType // 附加数据类型
	extraData     []byte        // 附加数据
}

func (p *Proto) setBuf(buf []byte) {
	p.buf = buf
	p.pb.SetBuf(nil)
}

func (p *Proto) back() {
	log.RunLogger.Printf("ffProto.Proto[%p].back: %v", p, p)

	if p.msg != nil {
		p.useState, p.limitState = useStateNone, limitStateInvalid

		backMessage(p.protoID, p.msg)
		p.msg = nil
		p.msgNeedMarshal = false

		backBuffer(p.buf)
		p.setBuf(nil)

		backProto(p)
	}
}

// BackAfterRecv 在接收完协议后, 尝试回收(比如解析失败时)
func (p *Proto) BackAfterRecv() {
	log.RunLogger.Printf("ffProto.Proto[%p].BackAfterRecv: %v", p, p)
	if p.useState == useStateRecv && p.msg != nil {
		p.back()
	}
}

// BackAfterSend 在底层发送协议后, 尝试回收
func (p *Proto) BackAfterSend() {
	log.RunLogger.Printf("ffProto.Proto[%p].BackAfterSend: %v", p, p)
	if p.useState == useStateSend && p.msg != nil {
		p.back()
	}
}

// BackAfterDispatch 缓存分发后, 不再需要时, 尝试回收
func (p *Proto) BackAfterDispatch() {
	log.RunLogger.Printf("ffProto.Proto[%p].BackAfterDispatch: %v", p, p)
	if p.limitState == limitStateRecv && p.useState == useStateBackAfterDispatch && p.msg != nil {
		p.back()
	}
}

// BackForce 强制回收, 慎用!!
func (p *Proto) BackForce() {
	log.RunLogger.Printf("ffProto.Proto[%p].BackForce: %v", p, p)
	p.back()
}

func (p *Proto) resetForSend(protoID MessageType) {
	p.useState, p.limitState = useStateSend, limitStateSend
	p.msgNeedMarshal = true
	p.protoID = protoID
	p.buf = p.buf[0:protoHeaderLength] // len(p.buf) = protoHeaderLength

	p.extraDataType = ExtraDataTypeNormal // 默认附加数据类型

	p.msg = applyMessage(p.protoID)
}

func (p *Proto) resetForRecv(header *ProtoHeader, bufLengthLimit int) {
	p.useState, p.limitState = useStateRecv, limitStateRecv
	p.msgNeedMarshal = false
	p.protoID = MessageType(header.protoID)

	p.buf = p.buf[0:bufLengthLimit] // len(p.buf) = protoHeaderLength+header.contentLength+extraDataLength

	p.pb.SetBuf(p.buf[protoHeaderLength : protoHeaderLength+header.contentLength]) // len(p.pb.buf) = header.contentLength

	header.marshalHeader(p.buf, true)

	p.extraDataType = header.recvExtraDataType

	p.msg = applyMessage(p.protoID)
}

// Message returns the Message
func (p *Proto) Message() proto.Message {
	return p.msg
}

// ProtoID returns the protoID
func (p *Proto) ProtoID() MessageType {
	return p.protoID
}

// Cap return cap of buffer
func (p *Proto) Cap() int {
	return cap(p.buf)
}

// BytesForSend returns the contents for send.
func (p *Proto) BytesForSend() []byte {
	return p.buf
}

// BytesForRecv 协议体和附加数据部分的buffer
func (p *Proto) BytesForRecv() []byte {
	return p.buf[protoHeaderLength:]
}

// OnRecvAllBytes 本协议的所有字节流接收完毕
func (p *Proto) OnRecvAllBytes(header *ProtoHeader) error {
	extraDataLen := p.extraDataType.BufferLength()
	extraDataOffset := protoHeaderLength + header.contentLength

	// 校验并记录附加数据
	if p.extraDataType == ExtraDataTypeUUID {

		copy(p.extraData[:extraDataLen], p.buf[extraDataOffset:])

	} else if p.extraDataType == ExtraDataTypeNormal {

		if header.lastRecvExtraDataNormal[0]+1 != p.buf[extraDataOffset+0] {
			return errCheckSerial
		}
		header.lastRecvExtraDataNormal[0] = p.buf[extraDataOffset+0]

	}

	// 清除附加数据
	if extraDataLen > 0 {
		p.buf = p.buf[0:extraDataOffset]
	}

	return nil
}

// Unmarshal parses the protocol buffer representation in the
// Buffer and places the decoded result in pb.  If the struct
// underlying pb does not match the data in the buffer, the results can be
// unpredictable.
// 接收完全字节流后, 需要获取协议内容Message时, 调用此方法
// 禁止外界保存返回的proto.Message
// 此方法在Proto使用期间内，只应该被调用一次
// 一旦调用此接口, 则认为外界需要修改协议内容, 则在转发协议时, 需要重新序列化协议内容到字节流
func (p *Proto) Unmarshal() error {
	log.RunLogger.Printf("ffProto.Proto[%p].Unmarshal: %v", p, p)

	p.msgNeedMarshal = true
	return p.pb.Unmarshal(p.msg)
}

// Marshal takes the protocol buffer and encodes it into the wire format, writing the result to the Buffer.
// 此方法在Proto使用期间内，只应该被调用一次
// 设置协议内容Message完毕后, 发送前夕, 调用此接口, 以生成待发送的字节流
func (p *Proto) Marshal(header *ProtoHeader) (err error) {
	log.RunLogger.Printf("ffProto.Proto[%p].Marshal: %v", p, p)

	var contentBuf []byte

	// 协议内容需要重新序列号称字节流
	if p.msgNeedMarshal {
		p.msgNeedMarshal = false
		p.pb.SetBuf(p.buf[protoHeaderLength:protoHeaderLength]) // len(p.pb.buf) = 0
		if err = p.pb.Marshal(p.msg); err != nil {
			return fmt.Errorf("Proto.Marshal: protoID[%s] msg is nil[%v] Marshal error[%v]",
				MessageType_name[int32(p.protoID)], p.msg == nil, err)
		}

		contentBuf = p.pb.Bytes()
	} else {
		contentBuf = p.buf[protoHeaderLength:]
	}

	contentLen := len(contentBuf)
	extraDataLen := p.extraDataType.BufferLength()
	bufLengthLimit := protoHeaderLength + contentLen + extraDataLen

	if bufLengthLimit > cap(p.buf) {
		// 协议体的缓冲区，进行了重新分配
		if cap(contentBuf) >= bufLengthLimit {
			copy(contentBuf[protoHeaderLength:protoHeaderLength+contentLen], contentBuf[0:contentLen])
			backBuffer(p.buf[:0]) // 缓存之前的buf
			p.buf = contentBuf[0:bufLengthLimit]
		} else {
			bufCapLimit := upperProtoBufferLength(bufLengthLimit)
			buf := make([]byte, bufLengthLimit, bufCapLimit)
			buf[protoHeaderExtraDataTypeOffset] = p.buf[protoHeaderExtraDataTypeOffset]
			copy(buf[protoHeaderLength:protoHeaderLength+contentLen], contentBuf[0:contentLen])
			backBuffer(p.buf[:0]) // 缓存之前的buf
			p.buf = buf
		}
	} else {
		p.buf = p.buf[0:bufLengthLimit]
		copy(p.buf[protoHeaderLength:bufLengthLimit], contentBuf)
	}

	// 写入协议头
	header.protoID = int32(p.protoID)
	header.contentLength = contentLen
	header.marshalHeader(p.buf, false)

	// 写入附加数据
	p.buf[protoHeaderExtraDataTypeOffset] = byte(p.extraDataType)
	if p.extraDataType == ExtraDataTypeUUID {
		copy(p.buf[bufLengthLimit-extraDataLen:], p.extraData[:extraDataLen])
	} else if p.extraDataType == ExtraDataTypeNormal {
		header.marshalSendExtraDataNormal(p.buf[bufLengthLimit-extraDataLen:])
	}

	return err
}

// ExtraData 附加数据
func (p *Proto) ExtraData() (extraData uint64) {
	return (uint64(p.extraData[0]) << 56) |
		(uint64(p.extraData[1]) << 48) |
		(uint64(p.extraData[2]) << 40) |
		(uint64(p.extraData[3]) << 32) |
		(uint64(p.extraData[4]) << 24) |
		(uint64(p.extraData[5]) << 16) |
		(uint64(p.extraData[6]) << 8) |
		(uint64(p.extraData[7]))
}

// SetExtraDataUUID 发送协议前, 必须设置附加数据类型及数据, 为发送而申请的协议, 其默认附加数据类型是ExtraDataTypeNormal
//	extraData: 附加数据
func (p *Proto) SetExtraDataUUID(extraData uint64) {
	log.RunLogger.Printf("ffProto.Proto[%p].SetExtraDataUUID: %v", p, p)

	if p.limitState == limitStateSend {
		p.useState, p.limitState = useStateSend, limitStateInvalid
		p.extraDataType = ExtraDataTypeUUID

		p.extraData[0] = byte(extraData >> 56)
		p.extraData[1] = byte(extraData >> 48)
		p.extraData[2] = byte(extraData >> 40)
		p.extraData[3] = byte(extraData >> 32)
		p.extraData[4] = byte(extraData >> 24)
		p.extraData[5] = byte(extraData >> 16)
		p.extraData[6] = byte(extraData >> 8)
		p.extraData[7] = byte(extraData)
	} else {
		log.FatalLogger.Printf("ffProto.Proto[%p].SetExtraDataUUID invalid limitState: %v", p, p)
	}
}

// SetExtraDataNormal 发送协议前, 必须设置附加数据类型及数据, 为发送而申请的协议, 其默认附加数据类型是ExtraDataTypeNormal
//	extraData: 附加数据
func (p *Proto) SetExtraDataNormal() {
	log.RunLogger.Printf("ffProto.Proto[%p].SetExtraDataNormal: %v", p, p)

	if p.limitState == limitStateSend {
		p.useState, p.limitState = useStateSend, limitStateInvalid
		p.extraDataType = ExtraDataTypeNormal
	} else {
		log.FatalLogger.Printf("ffProto.Proto[%p].SetExtraDataNormal invalid limitState: %v", p, p)
	}
}

// SetCacheWaitDispatch 协议被缓存以待分发
func (p *Proto) SetCacheWaitDispatch() {
	log.RunLogger.Printf("ffProto.Proto[%p].SetCacheWaitDispatch: %v", p, p)

	if p.limitState == limitStateRecv {
		p.useState = useStateCacheWaitDispatch
	} else {
		log.FatalLogger.Printf("ffProto.Proto[%p].SetCacheWaitDispatch invalid limitState: %v", p, p)
	}
}

// SetCacheDispatched 协议被分发处理了
func (p *Proto) SetCacheDispatched() {
	log.RunLogger.Printf("ffProto.Proto[%p].SetCacheDispatched: %v", p, p)

	if p.limitState == limitStateRecv && p.useState == useStateCacheWaitDispatch {
		p.useState = useStateBackAfterDispatch
	} else {
		log.FatalLogger.Printf("ffProto.Proto[%p].SetCacheDispatched invalid limitState or useState: %v", p, p)
	}
}

// SetCacheWaitSend 协议被缓存以待异步查询结果出来后再发送，在异步查询结果出来前，如果需要销毁，允许执行强制回收
// 设置成此状态的前提：
//	1. 逻辑处理过程中，涉及到异步查询等待结果
//	2. 协议最终会被返回给客户端或者转发给其他服务端
// 在异步查询结果出来前，如果需要销毁，允许执行强制回收。
func (p *Proto) SetCacheWaitSend() {
	log.RunLogger.Printf("ffProto.Proto[%p].SetCacheWaitSend: %v", p, p)

	if p.limitState == limitStateSend {
		p.useState = useStateCacheWaitSend
	} else {
		log.FatalLogger.Printf("ffProto.Proto[%p].SetCacheWaitSend invalid limitState: %v", p, p)
	}
}

// ChangeLimitStateRecvToSend 将limitState从limitStateRecv转到limitStateSend, 你要明白, 此操作意味着什么!
func (p *Proto) ChangeLimitStateRecvToSend() {
	log.RunLogger.Printf("ffProto.Proto[%p].ChangeLimitStateRecvToSend: %v", p, p)

	if p.limitState == limitStateRecv {
		p.limitState = limitStateSend
	} else {
		log.FatalLogger.Printf("ffProto.Proto[%p].ChangeLimitStateRecvToSend invalid limitState: %v", p, p)
	}
}

// String 返回Proto的自我描述
func (p *Proto) String() string {
	return fmt.Sprintf("%p protoID[%v] useState[%v] limitState[%v] msg[%v] extraDataType[%v] extraData[%v] buf[%v:%v:%v:%p]",
		p, p.protoID, p.useState, p.limitState, p.msg, p.extraDataType, p.extraData, len(p.buf), cap(p.buf), p.buf, p.buf)
}
