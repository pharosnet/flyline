package flyline

import (
	"sync"
	"testing"
	"time"
	"fmt"
)

func TestFlyline_QueueBuffer(t *testing.T) {
	N := 1000
	var start time.Time
	var dur time.Duration

	start = time.Now()
	benchmarkFlylineQueueBuffer(N, true)
	dur = time.Since(start)
	fmt.Printf("flyline-queuqbuffer: %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

	start = time.Now()
	benchmarkGoChan(N, 100, true)
	dur = time.Since(start)
	fmt.Printf("go-channel(100):  %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

	start = time.Now()
	benchmarkGoChan(N, 10, true)
	dur = time.Since(start)
	fmt.Printf("go-channel(10):   %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

	start = time.Now()
	benchmarkGoChan(N, 0, true)
	dur = time.Since(start)
	fmt.Printf("go-channel(0):    %d ops in %s (%d/sec)\n", N, dur, int(float64(N)/dur.Seconds()))

}

func benchmarkFlylineQueueBuffer(N int, validate bool) {
	buf := NewQueueBuffer()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < N; i++ {
			v, ok := buf.Recv()
			//fmt.Println(ok, v)
			if !ok {
				break
			}
			vInt := int64(-1)
			ValueScan(v, &vInt)
			if validate {
				if vInt != int64(i) {
					panic("out of order")
				}
			}
		}
		wg.Done()
	}()
	for i := 0; i < N; i++ {
		buf.Send(int64(i))
	}
	buf.Close()
	wg.Wait()
	//buf.Sync(context.Background())
}

func benchmarkGoChan(N, buffered int, validate bool) {
	ch := make(chan int64, buffered)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for i := 0; i < N; i++ {
			v := <-ch
			if validate {
				if v != int64(i) {
					panic("out of order")
				}
			}
		}
		wg.Done()
	}()
	for i := 0; i < N; i++ {
		ch <- int64(i)
	}
	close(ch)
	wg.Wait()
}

func BenchmarkFlyline_QueueBuffer(b *testing.B) {
	b.ReportAllocs()
	benchmarkFlylineQueueBuffer(b.N, false)
}

func BenchmarkGoChan100(b *testing.B) {
	b.ReportAllocs()
	benchmarkGoChan(b.N, 100, false)
}

func BenchmarkGoChan10(b *testing.B) {
	b.ReportAllocs()
	benchmarkGoChan(b.N, 10, false)
}

func BenchmarkGoChanUnbuffered(b *testing.B) {
	b.ReportAllocs()
	benchmarkGoChan(b.N, 0, false)
}