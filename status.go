package flyline

import "sync/atomic"

const (
	status_running = int64(1)
	status_closed  = int64(0)
)

// status: running, closeds
type status struct {
	v int64
}

func (s *status) setRunning() {
	atomic.StoreInt64(&s.v, status_running)
}

func (s *status) isRunning() bool {
	return status_running == atomic.LoadInt64(&s.v)
}

func (s *status) setClosed() {
	atomic.StoreInt64(&s.v, status_closed)
}

func (s *status) isClosed() bool {
	return status_closed == atomic.LoadInt64(&s.v)
}
