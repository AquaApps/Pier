package internal

import (
	"context"
	"io"
	"sync"
)

var bufferPool = &sync.Pool{
	New: func() interface{} {
		return make([]byte, 1600)
	},
}

func init() {
	for i := 0; i < 50; i++ {
		bytes := make([]byte, 1600)
		bufferPool.Put(bytes)
	}
}

func writer(stream io.Writer, out <-chan []byte, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case packet := <-out:
			num, _ := stream.Write(packet)
			if num == 1600 {
				bufferPool.Put(packet)
			}
		}
	}
}

func reader(stream io.Reader, in chan<- []byte, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			packet := bufferPool.Get().([]byte)
			num, err := stream.Read(packet)
			if err != nil {
				continue
			}
			in <- packet[:num]
		}
	}
}
