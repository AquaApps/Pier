package internal

import (
	"Pier/common"
	"context"
	"io"
)

func writer(stream io.Writer, out <-chan []byte, ctx context.Context) {
	for common.Opened(ctx) {
		_, _ = stream.Write(<-out)
	}
}

func reader(stream io.Reader, in chan<- []byte, ctx context.Context) {
	packet := make([]byte, 64*1024)
	for common.Opened(ctx) {
		num, err := stream.Read(packet)
		if err != nil {
			return
		}
		in <- packet[:num]
	}
}
