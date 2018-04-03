package flyline

import (
	"sync/atomic"
	"runtime"
	"time"
)

// Sequence New Function, value starts from -1.
func NewSequence() (seq *Sequence) {
	seq = &Sequence{lhs: [7]int64{}, value: int64(-1), rhs: [7]int64{}}
	return
}

// sequence, atomic operators.
type Sequence struct {
	lhs   [7]int64
	value int64
	rhs   [7]int64
}

// Atomic increment, if 5 times failed, then call runtime.Gosched().
func (s *Sequence) Incr() (value int64) {
	tryIncrTimes := 10
	for {
		tryIncrTimes--
		nextValue := s.Get() + 1
		ok := atomic.CompareAndSwapInt64(&s.value, s.value, nextValue)
		if ok {
			value = nextValue
			break
		}
		time.Sleep(100 * time.Microsecond)
		if tryIncrTimes < 0 {
			tryIncrTimes = 10
			runtime.Gosched()
		}
	}
	return
}

// Atomic get Sequence value.
func (s *Sequence) Get() (value int64) {
	value = atomic.LoadInt64(&s.value)
	return
}
