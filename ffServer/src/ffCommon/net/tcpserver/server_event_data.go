package tcpserver

import (
	"ffCommon/net/base"
	"ffCommon/uuid"
	"ffProto"
	"fmt"
)

type serverNetEventData struct {
	eventType base.NetEventType

	data   base.SessionNetEventData
	server *tcpServer
}

// Back 回收
func (s *serverNetEventData) Back() {
	eventType := s.NetEventType()
	if eventType == base.NetEventOff {
		s.server.onSessionClosed(s.SessionUUID())
	} else if eventType == base.NetEventEnd {
		s.server.back()
	}
	s.server = nil

	// 回收data
	s.data.Back()

	// 回收自身
	eventDataPool.back(s)
}

// SessionUUID 事件关联的session的UUID, 当NetEventType为NetEventEnd时无效
func (s *serverNetEventData) SessionUUID() uuid.UUID {
	return s.data.Session().UUID()
}

// Server 事件关联的server
func (s *serverNetEventData) Server() base.Server {
	return s.server
}

// NetEventType 获取事件类型
func (s *serverNetEventData) NetEventType() base.NetEventType {
	if s.eventType == base.NetEventInvalid {
		return s.data.NetEventType()
	}
	return s.eventType
}

// ManualClose 当NetEvent为NetEventOff时有效, 返回是不是主动关闭引发的Session断开
func (s *serverNetEventData) ManualClose() bool {
	return s.data.ManualClose()
}

// Proto 当NetEvent为NetEventProto时有效, 返回事件携带的协议
func (s *serverNetEventData) Proto() *ffProto.Proto {
	return s.data.Proto()
}

func (s *serverNetEventData) String() string {
	return fmt.Sprintf(`uuidServer[%v] dataSession[%v]`,
		s.server.uuid, s.data)
}

func newServerNetEventData() *serverNetEventData {
	return &serverNetEventData{}
}

func newServerNetEventDataFromSessionNetEventData(server *tcpServer, dataSession base.SessionNetEventData) *serverNetEventData {
	dataServer := eventDataPool.apply()
	dataServer.server, dataServer.data, dataServer.eventType = server, dataSession, base.NetEventInvalid
	return dataServer
}

func newServerNetEventDataEnd(server *tcpServer) *serverNetEventData {
	dataServer := eventDataPool.apply()
	dataServer.server, dataServer.data, dataServer.eventType = server, nil, base.NetEventEnd
	return dataServer
}
