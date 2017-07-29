package uuid

import "sync"

// uuidGeneratorSafe support multi goroutine call
type uuidGeneratorSafe struct {
	Generator

	muLock sync.Mutex
}

// Gen 生成 UUID
func (g *uuidGeneratorSafe) Gen() UUID {
	g.muLock.Lock()
	defer g.muLock.Unlock()

	return g.Generator.Gen()
}
