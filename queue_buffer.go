package flyline

import (
	"context"
	"sync"
	"time"
)

func NewQueueBuffer() Buffer {
	b := new(queueBuffer)
	b.seq = NewSequence()
	b.sts = new(status)
	b.queue = new(queue)
	b.mutex = new(sync.Mutex)
	b.sts.setRunning()
	return b
}

// Queue Buffer implements Buffer.
type queueBuffer struct {
	sts         *status
	queue       *queue
	seq         *Sequence
	mutex       *sync.Mutex
}

func (b *queueBuffer) Send(i interface{}) (err error) {
	if b.sts.isClosed() {
		err = ERR_BUF_SEND_CLOSED
		return
	}
	b.seq.Incr()
	b.queue.add(i)
	return
}

func (b *queueBuffer) Recv() (value interface{}, active bool) {
	active = true
	if b.sts.isClosed() && b.Len() == int64(0) {
		active = false
		return
	}
	value = b.queue.poll()
	b.seq.Decr()
	return
}

func (b *queueBuffer) Len() (length int64) {
	return b.seq.Get() + 1
}

func (b *queueBuffer) Close() (err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	if b.sts.isClosed() {
		err = ERR_BUF_CLOSE_CLOSED
		return
	}
	b.sts.setClosed()
	return
}

func (b *queueBuffer) Sync(ctx context.Context) (err error) {
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
			time.Sleep(500 * time.Microsecond)
		}
		if ok {
			break
		}
	}
	return
}
