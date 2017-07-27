package ffProto

import (
	"fmt"
)

/*
标准Proto结构:
						Header 					 	Content  	  ExtraData
2字节协议号+2字节Content长度+1字节标识附加数据类型		协议内容	最大extraDataMaxLength字节附加数据

协议号字段: BigEndian uint16
协议内容长度字段: BigEndian uint16

1字节标识附加数据类型为ExtraDataNormal时:
	附加数据为 extraDataNormalLength 字节 自增序号


1字节标识附加数据类型为ExtraDataUUID时:
	附加数据为 extraDataUUIDLength 字节 UUID
*/

const (
	// protoHeaderIDLength 协议号用几字节表示
	protoHeaderIDLength = 2

	// protoHeaderIDOffset 协议号偏移
	protoHeaderIDOffset = 0

	// protoHeaderContentLength 协议内容长度用几字节表示
	protoHeaderContentLength = 2

	// protoHeaderContentOffset 协议内容长度偏移
	protoHeaderContentOffset = protoHeaderIDOffset + protoHeaderIDLength

	// protoHeaderExtraDataTypeLength 标识附加数据类型的长度
	protoHeaderExtraDataTypeLength = 1

	// protoHeaderExtraDataTypeOffset 标识附加数据类型的偏移
	protoHeaderExtraDataTypeOffset = protoHeaderContentOffset + protoHeaderContentLength

	// protoHeaderLength 协议头长度
	protoHeaderLength = protoHeaderExtraDataTypeOffset + protoHeaderExtraDataTypeLength

	// protoMaxContentLength 协议内容最大长度限制
	protoMaxContentLength = protoMaxLength - protoHeaderLength - extraDataMaxLength

	// protoMaxLength 协议最大长度限制（协议头 + 协议体 + 附加数据）
	protoMaxLength = 4096
)

// useState 协议使用状态
type useState byte

const (
	// useStateNone 空状态
	useStateNone useState = 0

	// useStateRecv 接收，禁止调用强制回收
	useStateRecv useState = 1

	// useStateSend 发送，禁止调用强制回收
	useStateSend useState = 2

	// useStateCacheWaitDispatch 缓存以待分发，禁止调用强制回收
	useStateCacheWaitDispatch useState = 3

	// useStateCacheWaitSend 缓存以待异步查询结果出来后再发送，在异步查询结果出来前，如果需要销毁，允许执行强制回收
	useStateCacheWaitSend useState = 4
)

// ExtraDataType 附加数量类型
type ExtraDataType byte

const (
	// ExtraDataTypeNormal 客户端与服务端之间交互的标准附加数据, 1字节
	ExtraDataTypeNormal ExtraDataType = 0

	// ExtraDataTypeUUID 服务端与服务端之间交互的附加数据, 8字节UUID
	ExtraDataTypeUUID ExtraDataType = 1

	// 附加数据类型数量
	extraDataTypeCount = 2

	// 附加数据为ExtraDataNormal时的长度
	extraDataNormalLength = 1

	// 附加数据为ExtraDataUUID时的长度
	extraDataUUIDLength = 8

	// 附加数据最大长度
	extraDataMaxLength = 8
)

var errCheckSerial = fmt.Errorf("ffProto check serial number failed")

var extraDataLengthConfig = [extraDataTypeCount]int{
	extraDataNormalLength,
	extraDataUUIDLength,
}

func getExtraDataLength(extraDataType ExtraDataType) int {
	if extraDataType < extraDataTypeCount {
		return extraDataLengthConfig[extraDataType]
	}
	return 0
}
