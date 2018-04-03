package flyline

// buffer interface
type Buffer interface {
	// Send item into buffer.
	Send(i interface{}) error
	// Recv item from buffer, if closed eq true, then the buffer is closed and no remains.
	Recv(i interface{}) (closed bool, err error)
	// Recv original item from buffer, if closed eq true, then the buffer is closed and no remains.
	RecvOri() (i interface{}, closed bool, err error)
}
