package flyline

// buffer interface
type Buffer interface {
	Send(i interface{}) error
	Recv(i interface{}) (bool, error)
	RecvOri() (i interface{}, closed bool, err error)
}
