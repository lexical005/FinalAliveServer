// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ffProto.proto

/*
Package ffProto is a generated protocol buffer package.

It is generated from these files:
	ffProto.proto

It has these top-level messages:
	StDeviceInfo
	StAccountData
	MsgServerRegister
	MsgServerKeepAlive
	MsgPrepareLoginPlatformUniqueId
	MsgLoginPlatformUniqueId
	MsgLoginPlatformSidToken
	MsgReLogin
	MsgKick
	MsgEnterGameWorld
	MsgAgentDisConnect
	MsgKeepAlive
*/
package ffProto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MessageType int32

const (
	MessageType_ServerRegister               MessageType = 0
	MessageType_ServerKeepAlive              MessageType = 1
	MessageType_PrepareLoginPlatformUniqueId MessageType = 2
	MessageType_LoginPlatformUniqueId        MessageType = 3
	MessageType_LoginPlatformSidToken        MessageType = 4
	MessageType_ReLogin                      MessageType = 5
	MessageType_Kick                         MessageType = 6
	MessageType_EnterGameWorld               MessageType = 7
	MessageType_AgentDisConnect              MessageType = 8
	MessageType_KeepAlive                    MessageType = 9
)

var MessageType_name = map[int32]string{
	0: "ServerRegister",
	1: "ServerKeepAlive",
	2: "PrepareLoginPlatformUniqueId",
	3: "LoginPlatformUniqueId",
	4: "LoginPlatformSidToken",
	5: "ReLogin",
	6: "Kick",
	7: "EnterGameWorld",
	8: "AgentDisConnect",
	9: "KeepAlive",
}
var MessageType_value = map[string]int32{
	"ServerRegister":               0,
	"ServerKeepAlive":              1,
	"PrepareLoginPlatformUniqueId": 2,
	"LoginPlatformUniqueId":        3,
	"LoginPlatformSidToken":        4,
	"ReLogin":                      5,
	"Kick":                         6,
	"EnterGameWorld":               7,
	"AgentDisConnect":              8,
	"KeepAlive":                    9,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}
func (MessageType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// Generated by the msggen.py message compiler.
type StDeviceInfo struct {
	DeviceGUID            string `protobuf:"bytes,1,opt,name=deviceGUID" json:"deviceGUID,omitempty"`
	DeviceType            string `protobuf:"bytes,2,opt,name=deviceType" json:"deviceType,omitempty"`
	DeviceModel           string `protobuf:"bytes,3,opt,name=deviceModel" json:"deviceModel,omitempty"`
	OperatingSystem       string `protobuf:"bytes,4,opt,name=operatingSystem" json:"operatingSystem,omitempty"`
	OperatingSystemFamily string `protobuf:"bytes,5,opt,name=operatingSystemFamily" json:"operatingSystemFamily,omitempty"`
	ProcessorType         string `protobuf:"bytes,6,opt,name=processorType" json:"processorType,omitempty"`
	ProcessorFrequency    int32  `protobuf:"varint,7,opt,name=processorFrequency" json:"processorFrequency,omitempty"`
	ProcessorCount        int32  `protobuf:"varint,8,opt,name=processorCount" json:"processorCount,omitempty"`
	SystemMemorySize      int32  `protobuf:"varint,9,opt,name=systemMemorySize" json:"systemMemorySize,omitempty"`
	GraphicsMemorySize    int32  `protobuf:"varint,10,opt,name=graphicsMemorySize" json:"graphicsMemorySize,omitempty"`
}

func (m *StDeviceInfo) Reset()                    { *m = StDeviceInfo{} }
func (m *StDeviceInfo) String() string            { return proto.CompactTextString(m) }
func (*StDeviceInfo) ProtoMessage()               {}
func (*StDeviceInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *StDeviceInfo) GetDeviceGUID() string {
	if m != nil {
		return m.DeviceGUID
	}
	return ""
}

func (m *StDeviceInfo) GetDeviceType() string {
	if m != nil {
		return m.DeviceType
	}
	return ""
}

func (m *StDeviceInfo) GetDeviceModel() string {
	if m != nil {
		return m.DeviceModel
	}
	return ""
}

func (m *StDeviceInfo) GetOperatingSystem() string {
	if m != nil {
		return m.OperatingSystem
	}
	return ""
}

func (m *StDeviceInfo) GetOperatingSystemFamily() string {
	if m != nil {
		return m.OperatingSystemFamily
	}
	return ""
}

func (m *StDeviceInfo) GetProcessorType() string {
	if m != nil {
		return m.ProcessorType
	}
	return ""
}

func (m *StDeviceInfo) GetProcessorFrequency() int32 {
	if m != nil {
		return m.ProcessorFrequency
	}
	return 0
}

func (m *StDeviceInfo) GetProcessorCount() int32 {
	if m != nil {
		return m.ProcessorCount
	}
	return 0
}

func (m *StDeviceInfo) GetSystemMemorySize() int32 {
	if m != nil {
		return m.SystemMemorySize
	}
	return 0
}

func (m *StDeviceInfo) GetGraphicsMemorySize() int32 {
	if m != nil {
		return m.GraphicsMemorySize
	}
	return 0
}

type StAccountData struct {
	IsFresh    bool    `protobuf:"varint,1,opt,name=isFresh" json:"isFresh,omitempty"`
	ServerTime int32   `protobuf:"varint,2,opt,name=serverTime" json:"serverTime,omitempty"`
	ServerZone int32   `protobuf:"varint,3,opt,name=serverZone" json:"serverZone,omitempty"`
	Name       string  `protobuf:"bytes,4,opt,name=name" json:"name,omitempty"`
	BaseKeys   []int32 `protobuf:"varint,5,rep,packed,name=baseKeys" json:"baseKeys,omitempty"`
	BaseDatas  []int32 `protobuf:"varint,6,rep,packed,name=baseDatas" json:"baseDatas,omitempty"`
}

func (m *StAccountData) Reset()                    { *m = StAccountData{} }
func (m *StAccountData) String() string            { return proto.CompactTextString(m) }
func (*StAccountData) ProtoMessage()               {}
func (*StAccountData) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *StAccountData) GetIsFresh() bool {
	if m != nil {
		return m.IsFresh
	}
	return false
}

func (m *StAccountData) GetServerTime() int32 {
	if m != nil {
		return m.ServerTime
	}
	return 0
}

func (m *StAccountData) GetServerZone() int32 {
	if m != nil {
		return m.ServerZone
	}
	return 0
}

func (m *StAccountData) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *StAccountData) GetBaseKeys() []int32 {
	if m != nil {
		return m.BaseKeys
	}
	return nil
}

func (m *StAccountData) GetBaseDatas() []int32 {
	if m != nil {
		return m.BaseDatas
	}
	return nil
}

type MsgServerRegister struct {
	ServerType string `protobuf:"bytes,1,opt,name=serverType" json:"serverType,omitempty"`
	ServerID   int32  `protobuf:"varint,2,opt,name=serverID" json:"serverID,omitempty"`
}

func (m *MsgServerRegister) Reset()                    { *m = MsgServerRegister{} }
func (m *MsgServerRegister) String() string            { return proto.CompactTextString(m) }
func (*MsgServerRegister) ProtoMessage()               {}
func (*MsgServerRegister) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *MsgServerRegister) GetServerType() string {
	if m != nil {
		return m.ServerType
	}
	return ""
}

func (m *MsgServerRegister) GetServerID() int32 {
	if m != nil {
		return m.ServerID
	}
	return 0
}

type MsgServerKeepAlive struct {
}

func (m *MsgServerKeepAlive) Reset()                    { *m = MsgServerKeepAlive{} }
func (m *MsgServerKeepAlive) String() string            { return proto.CompactTextString(m) }
func (*MsgServerKeepAlive) ProtoMessage()               {}
func (*MsgServerKeepAlive) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type MsgPrepareLoginPlatformUniqueId struct {
	SubChannel        string `protobuf:"bytes,1,opt,name=subChannel" json:"subChannel,omitempty"`
	UUIDPlatformBound string `protobuf:"bytes,2,opt,name=UUIDPlatformBound" json:"UUIDPlatformBound,omitempty"`
	UUIDPlatformLogin string `protobuf:"bytes,3,opt,name=UUIDPlatformLogin" json:"UUIDPlatformLogin,omitempty"`
	RandomSalt        string `protobuf:"bytes,4,opt,name=randomSalt" json:"randomSalt,omitempty"`
	Timestamp         int32  `protobuf:"varint,5,opt,name=timestamp" json:"timestamp,omitempty"`
	Status            int32  `protobuf:"varint,6,opt,name=status" json:"status,omitempty"`
	Result            int32  `protobuf:"varint,7,opt,name=result" json:"result,omitempty"`
}

func (m *MsgPrepareLoginPlatformUniqueId) Reset()                    { *m = MsgPrepareLoginPlatformUniqueId{} }
func (m *MsgPrepareLoginPlatformUniqueId) String() string            { return proto.CompactTextString(m) }
func (*MsgPrepareLoginPlatformUniqueId) ProtoMessage()               {}
func (*MsgPrepareLoginPlatformUniqueId) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *MsgPrepareLoginPlatformUniqueId) GetSubChannel() string {
	if m != nil {
		return m.SubChannel
	}
	return ""
}

func (m *MsgPrepareLoginPlatformUniqueId) GetUUIDPlatformBound() string {
	if m != nil {
		return m.UUIDPlatformBound
	}
	return ""
}

func (m *MsgPrepareLoginPlatformUniqueId) GetUUIDPlatformLogin() string {
	if m != nil {
		return m.UUIDPlatformLogin
	}
	return ""
}

func (m *MsgPrepareLoginPlatformUniqueId) GetRandomSalt() string {
	if m != nil {
		return m.RandomSalt
	}
	return ""
}

func (m *MsgPrepareLoginPlatformUniqueId) GetTimestamp() int32 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *MsgPrepareLoginPlatformUniqueId) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *MsgPrepareLoginPlatformUniqueId) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

type MsgLoginPlatformUniqueId struct {
	TokenCustom string        `protobuf:"bytes,1,opt,name=tokenCustom" json:"tokenCustom,omitempty"`
	DeviceInfo  *StDeviceInfo `protobuf:"bytes,2,opt,name=deviceInfo" json:"deviceInfo,omitempty"`
	UUIDLogin   uint64        `protobuf:"varint,3,opt,name=UUIDLogin" json:"UUIDLogin,omitempty"`
	Result      int32         `protobuf:"varint,4,opt,name=result" json:"result,omitempty"`
}

func (m *MsgLoginPlatformUniqueId) Reset()                    { *m = MsgLoginPlatformUniqueId{} }
func (m *MsgLoginPlatformUniqueId) String() string            { return proto.CompactTextString(m) }
func (*MsgLoginPlatformUniqueId) ProtoMessage()               {}
func (*MsgLoginPlatformUniqueId) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *MsgLoginPlatformUniqueId) GetTokenCustom() string {
	if m != nil {
		return m.TokenCustom
	}
	return ""
}

func (m *MsgLoginPlatformUniqueId) GetDeviceInfo() *StDeviceInfo {
	if m != nil {
		return m.DeviceInfo
	}
	return nil
}

func (m *MsgLoginPlatformUniqueId) GetUUIDLogin() uint64 {
	if m != nil {
		return m.UUIDLogin
	}
	return 0
}

func (m *MsgLoginPlatformUniqueId) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

type MsgLoginPlatformSidToken struct {
	TokenPlatform string        `protobuf:"bytes,1,opt,name=tokenPlatform" json:"tokenPlatform,omitempty"`
	DeviceInfo    *StDeviceInfo `protobuf:"bytes,2,opt,name=deviceInfo" json:"deviceInfo,omitempty"`
	UUIDLogin     uint64        `protobuf:"varint,3,opt,name=UUIDLogin" json:"UUIDLogin,omitempty"`
	Result        int32         `protobuf:"varint,4,opt,name=result" json:"result,omitempty"`
}

func (m *MsgLoginPlatformSidToken) Reset()                    { *m = MsgLoginPlatformSidToken{} }
func (m *MsgLoginPlatformSidToken) String() string            { return proto.CompactTextString(m) }
func (*MsgLoginPlatformSidToken) ProtoMessage()               {}
func (*MsgLoginPlatformSidToken) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *MsgLoginPlatformSidToken) GetTokenPlatform() string {
	if m != nil {
		return m.TokenPlatform
	}
	return ""
}

func (m *MsgLoginPlatformSidToken) GetDeviceInfo() *StDeviceInfo {
	if m != nil {
		return m.DeviceInfo
	}
	return nil
}

func (m *MsgLoginPlatformSidToken) GetUUIDLogin() uint64 {
	if m != nil {
		return m.UUIDLogin
	}
	return 0
}

func (m *MsgLoginPlatformSidToken) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

type MsgReLogin struct {
	CheckData  string        `protobuf:"bytes,1,opt,name=checkData" json:"checkData,omitempty"`
	DeviceInfo *StDeviceInfo `protobuf:"bytes,2,opt,name=deviceInfo" json:"deviceInfo,omitempty"`
	Result     int32         `protobuf:"varint,3,opt,name=result" json:"result,omitempty"`
}

func (m *MsgReLogin) Reset()                    { *m = MsgReLogin{} }
func (m *MsgReLogin) String() string            { return proto.CompactTextString(m) }
func (*MsgReLogin) ProtoMessage()               {}
func (*MsgReLogin) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *MsgReLogin) GetCheckData() string {
	if m != nil {
		return m.CheckData
	}
	return ""
}

func (m *MsgReLogin) GetDeviceInfo() *StDeviceInfo {
	if m != nil {
		return m.DeviceInfo
	}
	return nil
}

func (m *MsgReLogin) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

type MsgKick struct {
	Result int32 `protobuf:"varint,1,opt,name=result" json:"result,omitempty"`
}

func (m *MsgKick) Reset()                    { *m = MsgKick{} }
func (m *MsgKick) String() string            { return proto.CompactTextString(m) }
func (*MsgKick) ProtoMessage()               {}
func (*MsgKick) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *MsgKick) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

type MsgEnterGameWorld struct {
	ServerID  int32  `protobuf:"varint,1,opt,name=serverID" json:"serverID,omitempty"`
	UUIDLogin uint64 `protobuf:"varint,2,opt,name=UUIDLogin" json:"UUIDLogin,omitempty"`
	Result    int32  `protobuf:"varint,3,opt,name=result" json:"result,omitempty"`
}

func (m *MsgEnterGameWorld) Reset()                    { *m = MsgEnterGameWorld{} }
func (m *MsgEnterGameWorld) String() string            { return proto.CompactTextString(m) }
func (*MsgEnterGameWorld) ProtoMessage()               {}
func (*MsgEnterGameWorld) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *MsgEnterGameWorld) GetServerID() int32 {
	if m != nil {
		return m.ServerID
	}
	return 0
}

func (m *MsgEnterGameWorld) GetUUIDLogin() uint64 {
	if m != nil {
		return m.UUIDLogin
	}
	return 0
}

func (m *MsgEnterGameWorld) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

type MsgAgentDisConnect struct {
}

func (m *MsgAgentDisConnect) Reset()                    { *m = MsgAgentDisConnect{} }
func (m *MsgAgentDisConnect) String() string            { return proto.CompactTextString(m) }
func (*MsgAgentDisConnect) ProtoMessage()               {}
func (*MsgAgentDisConnect) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type MsgKeepAlive struct {
	Number int32 `protobuf:"varint,1,opt,name=number" json:"number,omitempty"`
}

func (m *MsgKeepAlive) Reset()                    { *m = MsgKeepAlive{} }
func (m *MsgKeepAlive) String() string            { return proto.CompactTextString(m) }
func (*MsgKeepAlive) ProtoMessage()               {}
func (*MsgKeepAlive) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *MsgKeepAlive) GetNumber() int32 {
	if m != nil {
		return m.Number
	}
	return 0
}

func init() {
	proto.RegisterType((*StDeviceInfo)(nil), "StDeviceInfo")
	proto.RegisterType((*StAccountData)(nil), "StAccountData")
	proto.RegisterType((*MsgServerRegister)(nil), "MsgServerRegister")
	proto.RegisterType((*MsgServerKeepAlive)(nil), "MsgServerKeepAlive")
	proto.RegisterType((*MsgPrepareLoginPlatformUniqueId)(nil), "MsgPrepareLoginPlatformUniqueId")
	proto.RegisterType((*MsgLoginPlatformUniqueId)(nil), "MsgLoginPlatformUniqueId")
	proto.RegisterType((*MsgLoginPlatformSidToken)(nil), "MsgLoginPlatformSidToken")
	proto.RegisterType((*MsgReLogin)(nil), "MsgReLogin")
	proto.RegisterType((*MsgKick)(nil), "MsgKick")
	proto.RegisterType((*MsgEnterGameWorld)(nil), "MsgEnterGameWorld")
	proto.RegisterType((*MsgAgentDisConnect)(nil), "MsgAgentDisConnect")
	proto.RegisterType((*MsgKeepAlive)(nil), "MsgKeepAlive")
	proto.RegisterEnum("MessageType", MessageType_name, MessageType_value)
}

func init() { proto.RegisterFile("ffProto.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 762 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x55, 0xdb, 0x6e, 0xeb, 0x44,
	0x14, 0xc5, 0x49, 0x9c, 0xcb, 0xce, 0xc9, 0x39, 0x3e, 0x03, 0x07, 0x19, 0x54, 0x41, 0x88, 0x50,
	0x15, 0x55, 0xd0, 0x07, 0xe0, 0x07, 0x4a, 0x42, 0xab, 0xa8, 0x58, 0x54, 0x4e, 0x23, 0x24, 0xde,
	0x26, 0xce, 0x8e, 0x63, 0xd5, 0x9e, 0x71, 0x67, 0xc6, 0x95, 0xc2, 0x2f, 0xf0, 0x15, 0x3c, 0xf0,
	0xc2, 0x0f, 0xf0, 0x35, 0xfc, 0x0b, 0x9a, 0xb1, 0xe3, 0x4b, 0x12, 0x90, 0x78, 0x39, 0x2f, 0x51,
	0xf6, 0x5a, 0xcb, 0xde, 0x6b, 0xf6, 0x65, 0x0c, 0xa3, 0xed, 0xf6, 0x41, 0x70, 0xc5, 0xaf, 0x53,
	0xfd, 0x3b, 0xf9, 0xb3, 0x0d, 0xaf, 0x96, 0x6a, 0x8e, 0x2f, 0x51, 0x80, 0x0b, 0xb6, 0xe5, 0xe4,
	0x33, 0x80, 0x8d, 0x89, 0xee, 0x56, 0x8b, 0xb9, 0x6b, 0x8d, 0xad, 0xe9, 0xc0, 0xaf, 0x21, 0x15,
	0xff, 0xb8, 0x4f, 0xd1, 0x6d, 0xd5, 0x79, 0x8d, 0x90, 0x31, 0x0c, 0xf3, 0xc8, 0xe3, 0x1b, 0x8c,
	0xdd, 0xb6, 0x11, 0xd4, 0x21, 0x32, 0x85, 0x37, 0x3c, 0x45, 0x41, 0x55, 0xc4, 0xc2, 0xe5, 0x5e,
	0x2a, 0x4c, 0xdc, 0x8e, 0x51, 0x1d, 0xc3, 0xe4, 0x3b, 0x78, 0x77, 0x04, 0xdd, 0xd2, 0x24, 0x8a,
	0xf7, 0xae, 0x6d, 0xf4, 0xe7, 0x49, 0xf2, 0x25, 0x8c, 0x52, 0xc1, 0x03, 0x94, 0x92, 0x0b, 0x63,
	0xb2, 0x6b, 0xd4, 0x4d, 0x90, 0x5c, 0x03, 0x29, 0x81, 0x5b, 0x81, 0xcf, 0x19, 0xb2, 0x60, 0xef,
	0xf6, 0xc6, 0xd6, 0xd4, 0xf6, 0xcf, 0x30, 0xe4, 0x12, 0x5e, 0x97, 0xe8, 0x8c, 0x67, 0x4c, 0xb9,
	0x7d, 0xa3, 0x3d, 0x42, 0xc9, 0x15, 0x38, 0xd2, 0xb8, 0xf1, 0x30, 0xe1, 0x62, 0xbf, 0x8c, 0x7e,
	0x45, 0x77, 0x60, 0x94, 0x27, 0xb8, 0xf6, 0x10, 0x0a, 0x9a, 0xee, 0xa2, 0x40, 0xd6, 0xd4, 0x90,
	0x7b, 0x38, 0x65, 0x26, 0x7f, 0x59, 0x30, 0x5a, 0xaa, 0x9b, 0x20, 0xd0, 0x99, 0xe6, 0x54, 0x51,
	0xe2, 0x42, 0x2f, 0x92, 0xb7, 0x02, 0xe5, 0xce, 0xb4, 0xaa, 0xef, 0x1f, 0x42, 0xdd, 0x27, 0x89,
	0xe2, 0x05, 0xc5, 0x63, 0x94, 0xe4, 0x7d, 0xb2, 0xfd, 0x1a, 0x52, 0xf1, 0xbf, 0x70, 0x86, 0xa6,
	0x4d, 0x25, 0xaf, 0x11, 0x42, 0xa0, 0xc3, 0x68, 0x82, 0x45, 0x6b, 0xcc, 0x7f, 0xf2, 0x29, 0xf4,
	0xd7, 0x54, 0xe2, 0x3d, 0xee, 0xa5, 0x6b, 0x8f, 0xdb, 0x53, 0xdb, 0x2f, 0x63, 0x72, 0x01, 0x03,
	0xfd, 0x5f, 0xbb, 0x92, 0x6e, 0xd7, 0x90, 0x15, 0x30, 0xf9, 0x09, 0xde, 0x7a, 0x32, 0x5c, 0x9a,
	0xd7, 0xfb, 0x18, 0x46, 0x52, 0xa1, 0xa8, 0x59, 0xd4, 0x5d, 0x2a, 0x46, 0xad, 0x42, 0x74, 0xba,
	0x3c, 0x5a, 0xcc, 0x8b, 0x03, 0x94, 0xf1, 0xe4, 0x23, 0x20, 0xe5, 0x0b, 0xef, 0x11, 0xd3, 0x9b,
	0x38, 0x7a, 0xc1, 0xc9, 0x6f, 0x2d, 0xf8, 0xdc, 0x93, 0xe1, 0x83, 0xc0, 0x94, 0x0a, 0xfc, 0x91,
	0x87, 0x11, 0x7b, 0x88, 0xa9, 0xda, 0x72, 0x91, 0xac, 0x58, 0xf4, 0x9c, 0xe1, 0x62, 0x63, 0xb2,
	0x66, 0xeb, 0xd9, 0x8e, 0x32, 0x86, 0x71, 0x99, 0xb5, 0x44, 0xc8, 0x57, 0xf0, 0x76, 0xb5, 0x5a,
	0xcc, 0x0f, 0xcf, 0x7d, 0xcf, 0x33, 0xb6, 0x29, 0xe6, 0xfc, 0x94, 0x38, 0x56, 0x9b, 0x94, 0xc5,
	0xd0, 0x9f, 0x12, 0x3a, 0xb7, 0xa0, 0x6c, 0xc3, 0x93, 0x25, 0x8d, 0x55, 0x51, 0xda, 0x1a, 0xa2,
	0x8b, 0xa8, 0xa2, 0x04, 0xa5, 0xa2, 0x49, 0x6a, 0x86, 0xdc, 0xf6, 0x2b, 0x80, 0x7c, 0x0c, 0x5d,
	0xa9, 0xa8, 0xca, 0xa4, 0x99, 0x68, 0xdb, 0x2f, 0x22, 0x8d, 0x0b, 0x94, 0x59, 0xac, 0x8a, 0xf1,
	0x2d, 0xa2, 0xc9, 0xef, 0x16, 0xb8, 0x9e, 0x0c, 0xcf, 0x97, 0x61, 0x0c, 0x43, 0xc5, 0x9f, 0x90,
	0xcd, 0x32, 0xa9, 0x78, 0x52, 0xd4, 0xa1, 0x0e, 0x91, 0xaf, 0x0f, 0x9b, 0xae, 0xef, 0x05, 0x53,
	0x81, 0xe1, 0x37, 0xa3, 0xeb, 0xfa, 0x65, 0xe1, 0xd7, 0x04, 0xda, 0xbb, 0x3e, 0x70, 0x55, 0x81,
	0x8e, 0x5f, 0x01, 0x35, 0x8f, 0x9d, 0x86, 0xc7, 0x3f, 0xce, 0x78, 0x5c, 0x46, 0x9b, 0x47, 0xed,
	0x43, 0x6f, 0xb2, 0x31, 0x74, 0x20, 0x0a, 0x97, 0x4d, 0xf0, 0xfd, 0xf8, 0x7c, 0x06, 0xf0, 0x64,
	0xe8, 0xe7, 0x33, 0xa5, 0xdf, 0x11, 0xec, 0x30, 0x78, 0xd2, 0xc3, 0x5d, 0x98, 0xaa, 0x80, 0xff,
	0x6b, 0xa8, 0x4a, 0xd9, 0x6e, 0xa4, 0xfc, 0x02, 0x7a, 0x9e, 0x0c, 0xef, 0xa3, 0xe0, 0xa9, 0x26,
	0xb1, 0x1a, 0x12, 0x34, 0x6b, 0xf5, 0x03, 0x53, 0x28, 0xee, 0x68, 0x82, 0x3f, 0x73, 0x11, 0x6f,
	0x1a, 0x6b, 0x63, 0x35, 0xd7, 0xa6, 0x79, 0xf8, 0xd6, 0xbf, 0x1f, 0xbe, 0xe9, 0x24, 0x5f, 0xb6,
	0x9b, 0x10, 0x99, 0x9a, 0x47, 0x72, 0xc6, 0x19, 0xc3, 0x40, 0x4d, 0x2e, 0xe1, 0x95, 0xf6, 0x77,
	0x58, 0x3e, 0xfd, 0x34, 0xcb, 0x92, 0x35, 0x8a, 0x83, 0xc9, 0x3c, 0xba, 0xfa, 0xdb, 0x82, 0xa1,
	0x87, 0x52, 0xd2, 0x30, 0xff, 0x42, 0x10, 0x78, 0xdd, 0xbc, 0x08, 0x9c, 0x0f, 0xc8, 0x87, 0xf0,
	0xe6, 0x68, 0x97, 0x1d, 0x8b, 0x8c, 0xe1, 0xe2, 0xbf, 0x36, 0xd9, 0x69, 0x91, 0x4f, 0xe0, 0xdd,
	0x79, 0xaa, 0x7d, 0x42, 0x1d, 0x86, 0xca, 0xe9, 0x90, 0x21, 0xf4, 0x8a, 0x46, 0x3a, 0x36, 0xe9,
	0x43, 0x47, 0x97, 0xd8, 0xe9, 0x6a, 0x5f, 0xcd, 0x4a, 0x3a, 0x3d, 0xed, 0xeb, 0xe8, 0xd8, 0x4e,
	0x9f, 0x8c, 0x60, 0x50, 0xd9, 0x1c, 0xac, 0xbb, 0xe6, 0x4b, 0xfa, 0xed, 0x3f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0x78, 0x00, 0xf9, 0x79, 0x5a, 0x07, 0x00, 0x00,
}
