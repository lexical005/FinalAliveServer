package ffMySQL

import (
	"ffCommon/pool"
	"fmt"
)

type mysqlQueryResultPool struct {
	pool *pool.Pool
}

func (rpool *mysqlQueryResultPool) apply() *mysqlQueryReuslt {
	request, _ := rpool.pool.Apply().(*mysqlQueryReuslt)
	return request
}

func (rpool *mysqlQueryResultPool) back(request *mysqlQueryReuslt) {
	rpool.pool.Back(request)
}

func (rpool *mysqlQueryResultPool) String() string {
	return rpool.pool.String()
}

func (rpool *mysqlQueryResultPool) init(initCount int) error {
	if initCount < 1 {
		return fmt.Errorf("mysqlQueryResultPool.Init: invalid initCount[%v]", initCount)
	}

	rpool.pool = pool.New("mysqlQueryResultPool", false, newDBQueryResult, initCount, 50)
	return nil
}
