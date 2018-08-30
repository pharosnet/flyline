# flyline
Fly multi-producer / multi-consumer channels for Go

A flyline buffer works similar to a standard Go channel with the following features:

- It has a sync method. That means, you can wait for all after close buffer.
- It is lock-free and thread-safe.
- By default, it does not have a cap, you can use send filter to implement it.

# Getting Started

### Installing

To start using flyline, install Go and run `go get`:

```sh
$ go get github.com/pharosnet/flyline
```

This will retrieve the library. 

### Usage

It's easily to use, like standard Go channel.

#### ArrayBuffer (it's better than go channel)

```go
// buffer.go
package main

import (
	"github.com/pharosnet/flyline"
	"time"
	"context"
)

func main() {
	buf := flyline.NewArrayBuffer(1024*16)
	buf.Send(time.Now())
	value, ok := buf.Recv()
	recvTime := time.Time{}
	flyline.ValueScan(value, &recvTime)
	println(recvTime)
	// close buffer
	buf.Close()
	// change ctx to timeout, if sync with timeout. 
	buf.Sync(context.Background())
}
```

```sh
$ go run buffer.go 
```



#### QueueBuffer

```go
// buffer.go
package main

import (
	"github.com/pharosnet/flyline"
	"time"
	"context"
)

func main() {
	buf := flyline.NewQueueBuffer()
	buf.Send(time.Now())
	value, ok := buf.Recv()
	recvTime := time.Time{}
	flyline.ValueScan(value, &recvTime)
	println(recvTime)
	// close buffer
	buf.Close()
	// change ctx to timeout, if sync with timeout. 
	buf.Sync(context.Background())
}
```

```sh
$ go run buffer.go 
```

Benchmarks
----------------------------
Each of the following benchmark tests sends an incrementing sequence message from one goroutine to another. The receiving goroutine asserts that the message is received is the expected incrementing sequence value. Any failures cause a panic. Unless otherwise noted, all tests were run using `GOMAXPROCS=2`.

##### MacBook Pro 13" Retina, Mid 2015

* CPU: `Intel Core i5 @ 2.70 Ghz`
* Operation System: `OS X 10.13.4`
* Go Runtime: `Go 1.10.0`
* Go Architecture: `amd64`

Scenario | Per Operation Time
-------- | ------------------
Channels: Buffered, Non-blocking, GOMAXPROCS=1| 19.7 ns/op
Channels: Buffered, Non-blocking, GOMAXPROCS=2| 21.0 ns/op
Channels: Buffered, Non-blocking, GOMAXPROCS=3, Contended Write | 110.0 ns/op
ArrayBuffer: Buffered, Non-blocking, GOMAXPROCS=1| 11.9 ns/op
ArrayBuffer: Buffered, Non-blocking, GOMAXPROCS=2| 7.14 ns/op
ArrayBuffer: Buffered, Non-blocking, GOMAXPROCS=3, Contended Write | 11.9 ns/op

## Contact

Ryougi Nevermore [@ryougi](https://github.com/RyougiNevermore)

## License

`flyline` source code is available under the GPL-3 [License](/LICENSE).
