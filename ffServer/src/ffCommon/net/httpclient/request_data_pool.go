package httpclient

import "ffCommon/pool"

type postRequestPool struct {
	pool *pool.Pool
}

func (p *postRequestPool) apply() *PostRequest {
	agent, _ := p.pool.Apply().(*PostRequest)
	return agent
}

func (p *postRequestPool) back(agent *PostRequest) {
	p.pool.Back(agent)
}

func (p *postRequestPool) String() string {
	return p.pool.String()
}

func newPostRequestPool(nameOwner string, initCount int) *postRequestPool {
	funcCreator := func() interface{} {
		return newPostRequest()
	}

	return &postRequestPool{
		pool: pool.New(nameOwner+".PostRequest.pool", false, funcCreator, initCount, 50),
	}
}
