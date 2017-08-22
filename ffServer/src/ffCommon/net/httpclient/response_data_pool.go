package httpclient

import "ffCommon/pool"

type responseDataPool struct {
	pool *pool.Pool
}

func (p *responseDataPool) apply() *ResponseData {
	agent, _ := p.pool.Apply().(*ResponseData)
	return agent
}

func (p *responseDataPool) back(agent *ResponseData) {
	p.pool.Back(agent)
}

func (p *responseDataPool) String() string {
	return p.pool.String()
}

func newResponseDataPool(nameOwner string, initCount int) *responseDataPool {
	funcCreator := func() interface{} {
		return newResponseData()
	}

	return &responseDataPool{
		pool: pool.New(nameOwner+".ResponseData.pool", false, funcCreator, initCount, 50),
	}
}
