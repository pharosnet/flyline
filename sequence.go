package flyline

import (
	"runtime"
	"sync/atomic"
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
	incrTimes := 0
	nextValue := s.Get() + 1
	ok := atomic.CompareAndSwapInt64(&s.value, s.value, nextValue)
	for !ok {
		incrTimes++
		if incrTimes < 5 {
			time.Sleep(100 * time.Microsecond)
			continue
		}
		runtime.Gosched()
	}
	value = nextValue
	return
}

// Atomic get Sequence value.
func (s *Sequence) Get() (value int64) {
	value = atomic.LoadInt64(&s.value)
	return
}
