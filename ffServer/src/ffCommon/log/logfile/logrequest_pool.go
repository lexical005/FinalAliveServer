package logfile

import p "ffCommon/pool"

type logRequestPool struct {
	pool *p.Pool
}

func (lrq *logRequestPool) apply() (l *logRequest) {
	l, _ = lrq.pool.Apply().(*logRequest)
	return l
}

func (lrq *logRequestPool) back(l *logRequest) {
	lrq.pool.Back(l)
}

func (lrq *logRequestPool) String() string {
	return lrq.pool.String()
}
