package flyline

import (
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

type node struct {
	value interface{}
	prev  *node
	next  *node
}

type queue struct {
	wg      sync.WaitGroup
	in      *node
	inAvail *node
	out     *node
	outTail *node
	outFree *node
}

func (q *queue) newNode() *node {
	if q.inAvail != nil {
		n := q.inAvail
		q.inAvail = q.inAvail.next
		return n
	}
	times := 10
	for {
		q.inAvail = (*node)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.outFree)),
		))
		if q.inAvail == nil {
			return &node{}
		}
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.outFree)),
			unsafe.Pointer(q.inAvail), nil) {
			return q.newNode()
		}
		times--
		if times <= 0 {
			runtime.Gosched()
			times = 10
		}
	}
}

func (q *queue) peek() *node {
	return (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&q.in))))
}

func (q *queue) swapNode(old *node, new *node) bool {
	return atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.in)),
		unsafe.Pointer(old),
		unsafe.Pointer(new),
	)
}

func (q *queue) free(recv *node) {
	times := 10
	for {
		freed := (*node)(atomic.LoadPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.outFree)),
		))
		q.outTail.next = freed
		if atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.outFree)),
			unsafe.Pointer(freed), unsafe.Pointer(recv)) {
			return
		}
		times--
		if times <= 0 {
			runtime.Gosched()
			times = 10
		}
	}
}

var emptyNode = &node{}

func (q *queue) add(value interface{}) {
	n := q.newNode()
	n.value = value
	hasRemains := false
	times := 10
	for {
		n.next = q.peek()
		if n.next == emptyNode {
			if q.swapNode(n.next, n.next.next) {
				hasRemains = true
			}
		} else {
			if q.swapNode(n.next, n) {
				break
			}
		}
		times--
		if times <= 0 {
			runtime.Gosched()
			times = 10
		}
	}
	if hasRemains {
		q.wg.Done()
	}
}

func (q *queue) poll() interface{} {
	if q.out != nil {
		v := q.out.value
		if q.out.prev == nil {
			q.free(q.out)
			q.out = nil
		} else {
			q.out = q.out.prev
			q.out.next.prev = nil
		}
		return v
	}
	var n *node
	times := 10
	for {
		n = q.peek()
		if n == nil {
			q.wg.Add(1)
			if q.swapNode(n, emptyNode) {
				q.wg.Wait()
			} else {
				q.wg.Done()
			}
		} else if q.swapNode(n, nil) {
			break
		}
		times--
		if times <= 0 {
			runtime.Gosched()
			times = 10
		}
	}
	for n.next != nil {
		n.next.prev = n
		n = n.next
	}
	q.out = n
	q.outTail = n
	return q.poll()
}
