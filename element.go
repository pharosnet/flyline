package flyline

import (
	"sync"
	"time"
)

const (
	eleStatusEmpty = int64(0)
	eleStatusFull  = int64(1)
)

// Element of buffer.
// with status about set and get
type element struct {
	value  interface{}
	status int64
	lock   *sync.Mutex
}

// set and make status to be full
func (e *element) set(v interface{}) {
	for {
		e.lock.Lock()
		if e.status == eleStatusEmpty {
			e.value = v
			e.status = eleStatusFull
			e.lock.Unlock()
			return
		}
		e.lock.Unlock()
		time.Sleep(200 * time.Microsecond)
	}
}

// get and make status to be empty
func (e *element) pop() interface{} {
	for {
		e.lock.Lock()
		if e.status == eleStatusFull {
			v := e.value
			e.status = eleStatusEmpty
			e.lock.Unlock()
			return v
		}
		e.lock.Unlock()
		time.Sleep(200 * time.Microsecond)
	}
}
