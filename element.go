package flyline

import (
	"sync"
	"time"
)

const (
	ele_status_empty = int64(0)
	ele_status_full  = int64(1)
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
		if e.status == ele_status_empty {
			e.value = v
			e.status = ele_status_full
			e.lock.Unlock()
			return
		}
		e.lock.Unlock()
		time.Sleep(200 * time.Microsecond)
	}
	return
}

// get and make status to be empty
func (e *element) pop() interface{} {
	for {
		e.lock.Lock()
		if e.status == ele_status_full {
			v := e.value
			e.status = ele_status_empty
			e.lock.Unlock()
			return v
		}
		e.lock.Unlock()
		time.Sleep(200 * time.Microsecond)
	}
	return nil
}
