package flyline

import "sync/atomic"

const (
	statusRunning = int64(1)
	statusClosed  = int64(0)
)

// status: running, closed
type status struct {
	v   int64
	rhs [padding]int64
}

func (s *status) setRunning() {
	atomic.StoreInt64(&s.v, statusRunning)
}

func (s *status) isRunning() bool {
	return statusRunning == atomic.LoadInt64(&s.v)
}

func (s *status) setClosed() {
	atomic.StoreInt64(&s.v, statusClosed)
}

func (s *status) isClosed() bool {
	return statusClosed == atomic.LoadInt64(&s.v)
}
