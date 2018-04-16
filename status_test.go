package flyline

import "testing"

func TestStatus_Sample(t *testing.T) {
	s := new(status)
	s.setRunning()
	t.Logf("status running: %v", s.isRunning())
	s.setClosed()
	t.Logf("status close: %v", s.isClosed())
}
