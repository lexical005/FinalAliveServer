package main

type attrManager struct {
	health      int32 // 血量
	bodyDefence int32 // 身体防御
	headDefence int32 // 头部防御
}

// IsAlive 是否还活着
func (mgr *attrManager) IsAlive() bool {
	return mgr.health > 0
}

// IsHealthGreaterThan 生命值是否高于(不含等于)
func (mgr *attrManager) IsHealthGreaterThan(v int32) bool {
	return mgr.health > v
}

// IsHealthLessThan 生命值是否低于(不含等于)
func (mgr *attrManager) IsHealthLessThan(v int32) bool {
	return mgr.health < v
}

func newAttrManager() *attrManager {
	return &attrManager{}
}
