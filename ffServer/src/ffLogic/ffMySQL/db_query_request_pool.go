package ffMySQL

import (
	"ffCommon/pool"
	"fmt"
)

type mysqlQueryRequestPool struct {
	pool *pool.Pool
}

func (rpool *mysqlQueryRequestPool) apply() *mysqlQueryRequest {
	request, _ := rpool.pool.Apply().(*mysqlQueryRequest)
	return request
}

func (rpool *mysqlQueryRequestPool) back(request *mysqlQueryRequest) {
	rpool.pool.Back(request)
}

func (rpool *mysqlQueryRequestPool) String() string {
	return rpool.pool.String()
}

func (rpool *mysqlQueryRequestPool) init(initCount int) error {
	if initCount < 1 {
		return fmt.Errorf("mysqlQueryRequestPool.Init: invalid initCount[%v]", initCount)
	}

	rpool.pool = pool.New("mysqlQueryRequestPool", false, newDBQueryRequest, initCount, 50)
	return nil
}
