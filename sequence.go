package flyline

import (
	"sync/atomic"
	"runtime"
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

// Atomic increment
func (s *Sequence) Incr() (value int64) {
	times := 10
	for {
		times--
		nextValue := s.Get() + 1
		ok := atomic.CompareAndSwapInt64(&s.value, s.value, nextValue)
		if ok {
			value = nextValue
			break
		}
		if times <= 0 {
			times = 10
			runtime.Gosched()
		}
	}
	return
}

// Atomic decrement
func (s *Sequence) Decr() (value int64) {
	times := 10
	for {
		times--
		preValue := s.Get() - 1
		ok := atomic.CompareAndSwapInt64(&s.value, s.value, preValue)
		if ok {
			value = preValue
			break
		}
		if times <= 0 {
			times = 10
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
