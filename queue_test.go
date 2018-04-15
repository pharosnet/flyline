package flyline

import (
	"testing"
	"sync"
	"runtime"
	"time"
)

func TestQueue_Single(t *testing.T)  {
	q := new(queue)
	q.add(1)
	v := q.poll()
	t.Logf("queue poll: %v", v)
}


func TestQueue_multi(t *testing.T)  {
	runtime.GOMAXPROCS(8)
	q := new(queue)
	wg := new(sync.WaitGroup)
	for i := 0 ; i < 10 ; i ++ {
		go func(q *queue, wg *sync.WaitGroup) {
			q.add(time.Now())
			wg.Add(1)
		}(q, wg)
	}
	wg.Add(2)
	for i := 0 ; i < 2 ; i ++ {
		go func(q *queue, wg *sync.WaitGroup, i int) {
			for j := 0 ; j < 5 ; j ++ {
				t.Logf("queue poll %v : %v", i, q.poll())
				wg.Done()
			}
		}(q, wg, i)
		wg.Done()
	}
	wg.Wait()
}
