package tcpsession

import (
	"ffCommon/net/base"
	"ffProto"
	"fmt"
)

type sessionNetEventData struct {
	eventType base.NetEventType

	session     *tcpSession
	manualClose bool
	proto       *ffProto.Proto
}

// Back 回收
func (s *sessionNetEventData) Back() {

	if s.eventType == base.NetEventProto { // 回收proto
		s.proto.BackAfterDispatch()
		s.proto = nil
	} else if s.eventType == base.NetEventOff { // 回收tcpsession
		s.session.back()
		s.session = nil
	}

	s.eventType = base.NetEventInvalid

	// 回收自身
	eventDataPool.back(s)
}

// Session Session
func (s *sessionNetEventData) Session() base.Session {
	return s.session
}

// NetEventType 获取事件类型
func (s *sessionNetEventData) NetEventType() base.NetEventType {
	return s.eventType
}

// ManualClose 当NetEvent为NetEventOff时有效, 返回是不是主动关闭引发的Session断开
func (s *sessionNetEventData) ManualClose() bool {
	return s.manualClose
}

// Proto 当NetEvent为NetEventProto时有效, 返回事件携带的协议
func (s *sessionNetEventData) Proto() *ffProto.Proto {
	return s.proto
}

func (s *sessionNetEventData) String() string {
	return fmt.Sprintf(`eventType[%v] uuidSession[%v]`, s.eventType, s.session.uuid)
}

func newSessionNetEventData() *sessionNetEventData {
	return &sessionNetEventData{}
}

func newSessionNetEventOn(session *tcpSession) base.NetEventData {
	data := eventDataPool.apply()
	data.session, data.eventType = session, base.NetEventOn
	return data
}

func newSessionNetEventOff(session *tcpSession, manualClose bool) base.NetEventData {
	data := eventDataPool.apply()
	data.session, data.eventType, data.manualClose = session, base.NetEventOff, manualClose
	return data
}

func newSessionNetEventProto(session *tcpSession, proto *ffProto.Proto) base.NetEventData {
	data := eventDataPool.apply()
	data.session, data.eventType, data.proto = session, base.NetEventProto, proto
	proto.SetCacheWaitDispatch()
	return data
}
