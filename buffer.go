package flyline

import (
	"context"
	"errors"
)

var ERR_BUF_SEND_CLOSED error = errors.New("can not send item into the closed buffer")
var ERR_BUF_CLOSE_CLOSED error = errors.New("can not close buffer, buffer is closed")
var ERR_BUF_SYNC_UNCLOSED error = errors.New("can not sync buffer, buffer is not closed")

// buffer interface
type Buffer interface {
	// Send item into buffer.
	Send(i interface{}) (err error)
	// Recv value from buffer, if active eq false, then the buffer is closed and no remains.
	Recv() (value *Value, active bool)
	// Get remains length
	Len() (length int64)
	// Close Buffer, when closed, can not send item into buffer, but can recv remains.
	Close() (err error)
	// Sync, waiting for remains to be received. Only can be called after Close().
	Sync(ctx context.Context) (err error)
}

// Send Filter
type SendFilter interface {
	// if return true, item will be sent, if return false, item will be discarded.
	BeforeSend(i interface{}) bool
}

// Recv Filter
type RecvFilter interface {
	// it will be called before Buffer.Recv().
	AfterRecv(value *Value, closed bool)
}
