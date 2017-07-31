package ffProto

var listMessageID = []MessageType{
	MessageType_ServerRegister,
	MessageType_ServerKeepAlive,
	MessageType_PrepareLoginPlatformUniqueID,
	MessageType_LoginPlatformUniqueID,
	MessageType_LoginPlatformSidToken,
	MessageType_ReLogin,
	MessageType_Kick,
	MessageType_EnterGameWorld,
	MessageType_AgentDisConnect,
	MessageType_KeepAlive,
}

var mapMessageCreator = map[MessageType]func() interface{}{
	MessageType_ServerRegister: func() interface{} {
		return &MsgServerRegister{}
	},
	MessageType_ServerKeepAlive: func() interface{} {
		return &MsgServerKeepAlive{}
	},
	MessageType_PrepareLoginPlatformUniqueID: func() interface{} {
		return &MsgPrepareLoginPlatformUniqueID{}
	},
	MessageType_LoginPlatformUniqueID: func() interface{} {
		return &MsgLoginPlatformUniqueID{}
	},
	MessageType_LoginPlatformSidToken: func() interface{} {
		return &MsgLoginPlatformSidToken{}
	},
	MessageType_ReLogin: func() interface{} {
		return &MsgReLogin{}
	},
	MessageType_Kick: func() interface{} {
		return &MsgKick{}
	},
	MessageType_EnterGameWorld: func() interface{} {
		return &MsgEnterGameWorld{}
	},
	MessageType_AgentDisConnect: func() interface{} {
		return &MsgAgentDisConnect{}
	},
	MessageType_KeepAlive: func() interface{} {
		return &MsgKeepAlive{}
	},
}
