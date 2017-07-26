package session

import (
	"ffCommon/net/base"
	"ffCommon/pool"
	"ffCommon/uuid"
)

type sessionPool struct {
	pool *pool.Pool

	uuidGenerator *uuid.Generator
}

func (sp *sessionPool) apply() base.Session {
	s, _ := sp.pool.Apply().(*Session)
	s.uuid = sp.uuidGenerator.Gen()
	return s
}

func (sp *sessionPool) back(s base.Session) {
	sp.pool.Back(s)
}

func (sp *sessionPool) String() string {
	return sp.pool.String()
}
