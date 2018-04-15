package flyline

import (
	"errors"
	"context"
)

// buffer interface
type Buffer interface {
	// Send item into buffer.
	Send(i interface{}) error
	// Recv value from buffer, if closed eq true, then the buffer is closed and no remains.
	Recv() (value *Value, closed bool, err error)
	// Close Buffer, when closed, can not send item into buffer, but can recv remains.
	Close() (err error)
	// Sync, waiting for remains to be received. Only can be called after Close().
	Sync(ctx context.Context) (err error)
}

var ERR_BUF_SEND_CLOSED error = errors.New("can not send item into the closed buffer")
var ERR_BUF_RECV_CLOSED error = errors.New("can not recv item from the closed buffer")
var ERR_BUF_CLOSE_CLOSED error = errors.New("can not close buffer, buffer is closed")
var ERR_BUF_SYNC_UNCLOSED error = errors.New("can not sync buffer, buffer is not closed")

