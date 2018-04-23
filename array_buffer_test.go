package flyline

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestNewArrayBuffer(t *testing.T) {
	buf := NewArrayBuffer(4)
	t.Logf("new buffer: %v", buf)
	sendErr := buf.Send(time.Now())
	if sendErr != nil {
		t.Errorf("send failed, %v", sendErr)
		t.FailNow()
	}
	t.Logf("send ok, len = %v", buf.Len())
	v, ok := buf.Recv()
	recvTime := time.Time{}
	ValueScan(v, &recvTime)
	t.Logf("recv: [%v:%v] %v", ok, buf.Len(), recvTime)
	buf.Close()
	buf.Sync(context.Background())
}

func TestNewArrayBuffer_Sample(t *testing.T) {
	runtime.GOMAXPROCS(8)
	buf := NewArrayBuffer(8)
	t.Logf("new buffer: %v", buf)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	// send
	go func(buf Buffer, wg *sync.WaitGroup) {
		for i := 0; i < 10000; i++ {
			sendErr := buf.Send(time.Now())
			if sendErr != nil {
				t.Errorf("send failed, %v", sendErr)
				t.FailNow()
			}
			t.Logf("send ok, len = %v", buf.Len())
		}
		buf.Close()
		t.Logf("buf closed.")
		buf.Sync(context.Background())
		t.Logf("buf sync.")
		wg.Done()
	}(buf, wg)

	// recv
	go func(buf Buffer) {
		for {
			v, ok := buf.Recv()
			if !ok {
				t.Logf("recve: %v", ok)
				break
			}
			tt := time.Time{}
			ValueScan(v, &tt)
			t.Logf("recv: %v, %v", ok, tt)
		}
	}(buf)

	wg.Wait()
	time.Sleep(1 * time.Second)
}
