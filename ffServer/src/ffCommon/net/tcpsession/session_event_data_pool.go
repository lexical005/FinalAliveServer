package tcpsession

import (
	"ffCommon/net/base"
	"ffCommon/pool"
)

type sessionNetEventDataPool struct {
	pool *pool.Pool
}

func (s *sessionNetEventDataPool) apply() *sessionNetEventData {
	eventData, _ := s.pool.Apply().(*sessionNetEventData)
	return eventData
}

func (s *sessionNetEventDataPool) back(eventData base.SessionNetEventData) {
	s.pool.Back(eventData)
}

func (s *sessionNetEventDataPool) String() string {
	return s.pool.String()
}
