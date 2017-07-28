package tcpclient

import (
	"ffCommon/net/base"
	"ffCommon/pool"
)

type clientNetEventDataPool struct {
	pool *pool.Pool
}

func (s *clientNetEventDataPool) apply() *clientNetEventData {
	eventData, _ := s.pool.Apply().(*clientNetEventData)
	return eventData
}

func (s *clientNetEventDataPool) back(eventData base.ClientNetEventData) {
	s.pool.Back(eventData)
}

func (s *clientNetEventDataPool) String() string {
	return s.pool.String()
}
