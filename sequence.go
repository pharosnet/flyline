package flyline

// sequence, atomic operators.
type sequence struct {
	padding []int64
	value int64
}

func (s *sequence) Incr() (value int64, err error) {

	return
}

func (s *sequence) Get() (value int64, err error) {

	return
}


