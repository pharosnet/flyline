package flyline

import (
	"runtime"
	"sync/atomic"
	"time"
)

const (
	padding = 7
)

// Sequence New Function, value starts from -1.
func NewSequence() (seq *Sequence) {
	seq = &Sequence{value: int64(-1)}
	return
}

// sequence, atomic operators.
type Sequence struct {
	value int64
	rhs   [padding]int64
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
		time.Sleep(1 * time.Nanosecond)
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
		time.Sleep(1 * time.Nanosecond)
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
