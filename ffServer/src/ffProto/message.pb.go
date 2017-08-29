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
	MessageType_LeaveTeam,
	MessageType_EnterMatchServer,
	MessageType_LeaveMatchServer,
	MessageType_StartMatch,
	MessageType_StopMatch,
	MessageType_MatchResult,
	MessageType_BattleStartSync,
	MessageType_BattleMember,
	MessageType_BattleMemberLeave,
	MessageType_BattleProp,
	MessageType_BattlePickProp,
	MessageType_BattleDropProp,
	MessageType_BattleAddProp,
	MessageType_BattleRemoveProp,
	MessageType_BattleRunAway,
	MessageType_BattleSettle,
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
	MessageType_LeaveTeam: func() interface{} {
		return &MsgLeaveTeam{}
	},
	MessageType_EnterMatchServer: func() interface{} {
		return &MsgEnterMatchServer{}
	},
	MessageType_LeaveMatchServer: func() interface{} {
		return &MsgLeaveMatchServer{}
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
	MessageType_BattleStartSync: func() interface{} {
		return &MsgBattleStartSync{}
	},
	MessageType_BattleMember: func() interface{} {
		return &MsgBattleMember{}
	},
	MessageType_BattleMemberLeave: func() interface{} {
		return &MsgBattleMemberLeave{}
	},
	MessageType_BattleProp: func() interface{} {
		return &MsgBattleProp{}
	},
	MessageType_BattlePickProp: func() interface{} {
		return &MsgBattlePickProp{}
	},
	MessageType_BattleDropProp: func() interface{} {
		return &MsgBattleDropProp{}
	},
	MessageType_BattleAddProp: func() interface{} {
		return &MsgBattleAddProp{}
	},
	MessageType_BattleRemoveProp: func() interface{} {
		return &MsgBattleRemoveProp{}
	},
	MessageType_BattleRunAway: func() interface{} {
		return &MsgBattleRunAway{}
	},
	MessageType_BattleSettle: func() interface{} {
		return &MsgBattleSettle{}
	},
}
