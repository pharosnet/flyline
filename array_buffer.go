package flyline

import (
	"context"
	"runtime"
	"sync"
	"time"
)

// Note: The array capacity must be a power of two, e.g. 2, 4, 8, 16, 32, 64, etc.
func NewArrayBuffer(capacity int64) Buffer {
	b := &arrayBuffer{
		capacity: capacity,
		buffer:   newArray(capacity),
		wdSeq:    NewSequence(),
		wpSeq:    NewSequence(),
		rdSeq:    NewSequence(),
		rpSeq:    NewSequence(),
		sts:      &status{},
		mutex:    &sync.Mutex{},
	}
	b.sts.setRunning()
	return b
}

type arrayBuffer struct {
	capacity int64
	buffer   *array
	wpSeq    *Sequence
	wdSeq    *Sequence
	rpSeq    *Sequence
	rdSeq    *Sequence
	sts      *status
	mutex    *sync.Mutex
}

func (b *arrayBuffer) Send(i interface{}) (err error) {
	if b.sts.isClosed() {
		err = ERR_BUF_SEND_CLOSED
		return
	}
	next := b.wpSeq.Incr()
	times := 10
	for {
		times--
		if next-b.capacity <= b.rdSeq.Get() && next == b.wdSeq.Get()+1 {
			b.buffer.set(next, i)
			b.wdSeq.Incr()
			break
		}
		time.Sleep(ns1)
		if times <= 0 {
			runtime.Gosched()
			times = 10
		}
	}
	return
}

func (b *arrayBuffer) Recv() (value interface{}, active bool) {
	active = true
	if b.sts.isClosed() && b.Len() == int64(0) {
		active = false
		return
	}
	times := 10
	next := b.rpSeq.Incr()
	for {
		if next <= b.wdSeq.Get() && next == b.rdSeq.Get()+1 {
			value = b.buffer.get(next)
			b.rdSeq.Incr()
			break
		}
		time.Sleep(ns1)
		if times <= 0 {
			runtime.Gosched()
			times = 10
		}
	}
	return
}

func (b *arrayBuffer) Len() (length int64) {
	length = b.wpSeq.Get() - b.rdSeq.Get()
	return
}

func (b *arrayBuffer) Close() (err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.sts.isClosed() {
		err = ERR_BUF_CLOSE_CLOSED
		return
	}
	b.sts.setClosed()
	return
}

func (b *arrayBuffer) Sync(ctx context.Context) (err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.sts.isRunning() {
		err = ERR_BUF_SYNC_UNCLOSED
		return
	}
	for {
		ok := false
		select {
		case <-ctx.Done():
			ok = true
			break
		default:
			if b.Len() == int64(0) {
				ok = true
				break
			}
			time.Sleep(ms500)
		}
		if ok {
			break
		}
	}
	return
}
