package internal

import (
	"Pier/common"
	"context"
	"io"
	"log"
	"runtime"
)

func writer(stream io.Writer, out io.Reader, ctx context.Context) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	buffer := make([]byte, 4*1024)
	for common.Opened(ctx) {
		_, _ = io.CopyBuffer(stream, out, buffer)
	}
}

func reader(stream io.Reader, in io.Writer, ctx context.Context) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	buffer := make([]byte, 4*1024)
	for common.Opened(ctx) {
		_, err := io.CopyBuffer(in, stream, buffer)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
