package ffProto

var listMessageID = []MessageType{
	MessageType_KeepAlive,
	MessageType_ServerRegister,
	MessageType_PrepareLoginPlatformUniqueId,
	MessageType_LoginPlatformUniqueId,
	MessageType_LoginPlatformSidToken,
	MessageType_ReLogin,
	MessageType_Kick,
	MessageType_EnterGameWorld,
	MessageType_AgentDisConnect,
	MessageType_EnterTeam,
	MessageType_InviteJoinTeam,
	MessageType_AnswerJoinTeam,
	MessageType_LeaveJoinTeam,
	MessageType_StartMatch,
	MessageType_StopMatch,
	MessageType_MatchResult,
}

var mapMessageCreator = map[MessageType]func() interface{}{
	MessageType_KeepAlive: func() interface{} {
		return &MsgKeepAlive{}
	},
	MessageType_ServerRegister: func() interface{} {
		return &MsgServerRegister{}
	},
	MessageType_PrepareLoginPlatformUniqueId: func() interface{} {
		return &MsgPrepareLoginPlatformUniqueId{}
	},
	MessageType_LoginPlatformUniqueId: func() interface{} {
		return &MsgLoginPlatformUniqueId{}
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
	MessageType_EnterTeam: func() interface{} {
		return &MsgEnterTeam{}
	},
	MessageType_InviteJoinTeam: func() interface{} {
		return &MsgInviteJoinTeam{}
	},
	MessageType_AnswerJoinTeam: func() interface{} {
		return &MsgAnswerJoinTeam{}
	},
	MessageType_LeaveJoinTeam: func() interface{} {
		return &MsgLeaveJoinTeam{}
	},
	MessageType_StartMatch: func() interface{} {
		return &MsgStartMatch{}
	},
	MessageType_StopMatch: func() interface{} {
		return &MsgStopMatch{}
	},
	MessageType_MatchResult: func() interface{} {
		return &MsgMatchResult{}
	},
}
