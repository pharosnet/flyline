package flyline

// buffer interface
type Buffer interface {
	Send(i interface{}) error
	Recv(i interface{}) error
	RecvOri() (i interface{}, err error)
}
