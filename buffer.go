package flyline

import (
	"errors"
)

// buffer interface
type Buffer interface {
	// Send item into buffer.
	Send(i interface{}) error
	// Recv item from buffer, if closed eq true, then the buffer is closed and no remains.
	Recv(i interface{}) (closed bool, err error)
	// Recv original item from buffer, if closed eq true, then the buffer is closed and no remains.
	RecvOri() (i interface{}, closed bool, err error)
}

var ERR_BUF_SEND_CLOSED error = errors.New("cant not send item into the closed buffer")
var ERR_BUF_RECV_CLOSED error = errors.New("cant not recv item from the closed buffer")
