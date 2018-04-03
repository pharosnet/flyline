package flyline

import (
	"reflect"
	"runtime"
	"sync"
	"time"
)

const (
	status_running = iota
	status_close
	status_sync
)

// Array Buffer New Function.
func NewArrayBuffer(size int64) Buffer {
	buf := new(arrayBuffer)
	buf.status = status_running
	buf.size = size
	buf.array = make([]*element, size)
	for i := int64(0); i < size; i++ {
		buf.array[i] = &element{
			status: ele_status_empty,
			lock:   new(sync.Mutex),
		}
	}
	buf.sendSeq = NewSequence()
	buf.recvSeq = NewSequence()
	return buf
}

// Array Buffer implements Buffer.
type arrayBuffer struct {
	status  int
	array   []*element
	size    int64
	sendSeq *Sequence
	recvSeq *Sequence
}

// get element index by sequence value.
func (b *arrayBuffer) getElementIndexBySequence(seq int64) int64 {
	return seq % b.size
}

// Send item into buffer, thread safe.
func (b *arrayBuffer) Send(i interface{}) error {
	if b.status == status_sync || b.status == status_close {
		return ERR_BUF_SEND_CLOSED
	}
	trySendTimes := 10
	for {
		trySendTimes--
		if b.sendSeq.Get()-b.size < b.recvSeq.Get() {
			nextSendSeq := b.sendSeq.Incr()
			b.array[b.getElementIndexBySequence(nextSendSeq)].set(i)
			break
		}
		time.Sleep(100 * time.Microsecond)
		if trySendTimes < 0 {
			runtime.Gosched()
			trySendTimes = 10
		}
	}
	return nil
}

// Recv item from buffer by type mapping, thread safe.
// TODO: MORE SAFE AT TYPE JUDGING.
func (b *arrayBuffer) Recv(i interface{}) (closed bool, err error) {
	var oriRecved interface{}
	oriRecved, closed, err = b.RecvOri()
	if err != nil {
		return
	}
	oriRecvedValue := reflect.ValueOf(oriRecved)
	recvedvalue := reflect.ValueOf(i)
	if recvedvalue.Kind() == reflect.Ptr {
		if oriRecvedValue.Kind() == reflect.Ptr {
			recvedvalue.Elem().Set(oriRecvedValue.Elem())
			return
		}
		recvedvalue.Elem().Set(oriRecvedValue)
		return
	} else if recvedvalue.Kind() == reflect.Slice || recvedvalue.Kind() == reflect.Map {
		recvedvalue.Set(oriRecvedValue)
		return
	}
	return
}

// Recv original item from buffer, thread safe.
func (b *arrayBuffer) RecvOri() (i interface{}, closed bool, err error) {
	if b.status == status_close {
		closed = true
		err = ERR_BUF_RECV_CLOSED
		return
	}
	tryRecvTimes := 10
	awaitRecvSeq := b.recvSeq.Incr()
	for {
		tryRecvTimes--
		if awaitRecvSeq <= b.sendSeq.Get() {
			i = b.array[b.getElementIndexBySequence(awaitRecvSeq)].pop()
			if b.status != status_running {
				closed = true
			}
			return
		}
		time.Sleep(100 * time.Microsecond)
		if tryRecvTimes < 0 {
			runtime.Gosched()
			tryRecvTimes = 10
		}
	}
	return
}
