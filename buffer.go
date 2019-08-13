package flyline

import (
	"context"
	"errors"
)

var ErrBufSendClosed = errors.New("can not send item into the closed buffer")
var ErrBufCloseClosed = errors.New("can not close buffer, buffer is closed")
var ErrBufSyncUnclosed = errors.New("can not sync buffer, buffer is not closed")

// buffer interface
type Buffer interface {
	// Send item into buffer.
	Send(i interface{}) (err error)
	// Recv value from buffer, if active eq false, then the buffer is closed and no remains.
	Recv() (value interface{}, active bool)
	// Get remains length
	Len() (length int64)
	// Close Buffer, when closed, can not send item into buffer, but can recv remains.
	Close() (err error)
	// Sync, waiting for remains to be received. Only can be called after Close().
	Sync(ctx context.Context) (err error)
}
