package ffMySQL

import (
	"sync"
)

type mysqlConn struct {
	valid bool

	chClose        chan struct{}
	chQueryRequest chan *mysqlQueryRequest

	wgClose sync.WaitGroup
}

func (conn *mysqlConn) close() {
	conn.valid = false

	close(conn.chClose)

	conn.wgClose.Wait()

	close(conn.chQueryRequest)

	// 归还
	for req := range conn.chQueryRequest {
		req.back()
	}
}

func (conn *mysqlConn) addQuery(req *mysqlQueryRequest) {
	if conn.valid {
		conn.chQueryRequest <- req
	}
}

// dbQueryLoop 数据库查询循环
func (conn *mysqlConn) queryLoop(params ...interface{}) {
	conn.wgClose.Add(1)

deadLoop:
	for {
		select {
		case req := <-conn.chQueryRequest:
			req.doQuery()
		case <-conn.chClose:
			break deadLoop
		}
	}
}
func (conn *mysqlConn) queryLoopEnd(isPanic bool) {
	conn.wgClose.Done()
}
