package pool

import (
	"fmt"
	"sync"
)

const (
	// 最低初始大小
	minInitCount = 2

	// 扩展时最低增长大小
	minIncreaseCount = 2

	// 扩展时最低增长百分比
	minIncreaseRate = 10
)

// ElementCreator 缓存元素的创建器
type ElementCreator func() interface{}

// Pool 缓存元素的管理器
type Pool struct {
	name    string         // 缓存管理器的名称，用以区分彼此
	isFree  bool           //自由管理模式
	creator ElementCreator // 创建工厂

	pool []interface{} // 内部缓存

	increaseRate int //缓存使用完毕后，增长比例（百分比）

	mu sync.Mutex
}

// Apply 获取
func (m *Pool) Apply() (e interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	l := len(m.pool)

	// 扩展
	if l == 0 {
		capExtand := cap(m.pool) * m.increaseRate / 100
		if capExtand < minIncreaseCount {
			capExtand = minIncreaseCount
		}
		m.pool = make([]interface{}, 0, cap(m.pool)+capExtand)
		for i := 0; i < capExtand; i++ {
			m.pool = append(m.pool, m.creator())
		}
		l = len(m.pool)
	}

	// 返回尾部元素
	e = m.pool[l-1]
	m.pool = m.pool[:l-1]

	return e
}

// Back 归还
func (m *Pool) Back(e interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.pool) == cap(m.pool) {
		if !m.isFree {
			// 简单检查重复归还
			panic(fmt.Sprintf("pool.Pool.Back:[%s] duplicate element back or new element back", m.name))
		}

		// 扩展
		capExtand := cap(m.pool) * m.increaseRate / 100
		if capExtand < minIncreaseCount {
			capExtand = minIncreaseCount
		}
		pool := make([]interface{}, 0, cap(m.pool)+capExtand)
		m.pool = append(pool, m.pool...)
	}

	m.pool = append(m.pool, e)
}

// String 返回Pool的自我描述
func (m *Pool) String() string {
	return fmt.Sprintf("name[%s] isFree[%v] len(pool)[%v] cap(pool)[%v]", m.name, m.isFree, len(m.pool), cap(m.pool))
}

// New 新建一个缓存管理
// name: 缓存管理器的名称
// isFree：true, 允许管理外界创建的元素
// creator：缓存元素的创建器
// initCount：初始缓存多少个元素，最低 minInitCount
// increaseRate：缓存的元素都被使用后，增长比例（百分比），最低 minIncreaseRate%，普通用途，建议配置为50%
func New(
	name string,
	isFree bool,
	creator ElementCreator,
	initCount int,
	increaseRate int) *Pool {
	if initCount < minInitCount {
		initCount = minInitCount
	}

	if increaseRate < minIncreaseRate {
		increaseRate = minIncreaseRate
	}

	pool := make([]interface{}, initCount, initCount)
	for i := 0; i < initCount; i++ {
		pool[i] = creator()
	}

	return &Pool{
		name:    name,
		isFree:  isFree,
		creator: creator,

		pool: pool,

		increaseRate: increaseRate,
	}
}
