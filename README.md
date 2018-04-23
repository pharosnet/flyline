# flyline
Fly multi-producer / multi-consumer channels for Go

A flyline buffer works similar to a standard Go channel with the following features:

- It has a sync method. That means, you can wait for all after close buffer.
- It is lock-free and thread-safe.
- By default, it does not have a cap, you can use send filter to implement it.

# Getting Started

### Installing

To start using fastlane, install Go and run `go get`:

```sh
$ go get github.com/pharosnet/flyline
```

This will retrieve the library. 

### Usage

It's easily to use, like standard Go channel.

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

## Contact

Ryougi Nevermore [@ryougi](https://github.com/RyougiNevermore)

## License

`flyline` source code is available under the GPL-3 [License](/LICENSE).
