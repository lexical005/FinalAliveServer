package tcpserver

import (
	"ffCommon/net/base"
	"ffCommon/pool"
)

type serverNetEventDataPool struct {
	pool *pool.Pool
}

func (s *serverNetEventDataPool) apply() *serverNetEventData {
	eventData, _ := s.pool.Apply().(*serverNetEventData)
	return eventData
}

func (s *serverNetEventDataPool) back(eventData base.ServerNetEventData) {
	s.pool.Back(eventData)
}

func (s *serverNetEventDataPool) String() string {
	return s.pool.String()
}
