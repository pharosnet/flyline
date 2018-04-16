package flyline

import (
	"runtime"
	"sync"
	"testing"
)

func TestNewSequence(t *testing.T) {
	seq := NewSequence()
	if seq.Get() == -1 {
		t.Logf("seq %v", seq)
	} else {
		t.Errorf("seq %v", seq)
	}
}

func TestSequence_Incr(t *testing.T) {
	seq := NewSequence()
	v := seq.Get()
	n := seq.Incr()
	if n == v+1 {
		t.Logf("seq incr: %v -> %v", v, n)
	} else {
		t.Errorf("seq incr: %v -> %v", v, n)
	}
	t.Log("thread safe testing....")
	runtime.GOMAXPROCS(4)
	wg := new(sync.WaitGroup)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(seq *Sequence, wg *sync.WaitGroup) {
			defer wg.Done()
			t.Logf("seq now: %v incr : %v, get : %v", seq.Get(), seq.Incr(), seq.Get())
		}(seq, wg)
	}
	wg.Wait()
}

func TestSequence_Decr(t *testing.T) {
	seq := NewSequence()
	v := seq.Get()
	n := seq.Decr()
	if n == v-1 {
		t.Logf("seq incr: %v -> %v", v, n)
	} else {
		t.Errorf("seq incr: %v -> %v", v, n)
	}
	t.Log("thread safe testing....")
	runtime.GOMAXPROCS(4)
	wg := new(sync.WaitGroup)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(seq *Sequence, wg *sync.WaitGroup) {
			defer wg.Done()
			t.Logf("seq now: %v incr : %v, get : %v", seq.Get(), seq.Decr(), seq.Get())
		}(seq, wg)
	}
	wg.Wait()
}
