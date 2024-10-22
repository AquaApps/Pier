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
		//_, _ = io.CopyBuffer(stream, out, buffer)
		//_, _ = io.Copy(stream, out)
		n, _ := out.Read(buffer)
		_, _ = stream.Write(buffer[:n])
	}
}

func reader(stream io.Reader, in io.Writer, ctx context.Context) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	buffer := make([]byte, 4*1024)
	for common.Opened(ctx) {
		n, err := stream.Read(buffer)
		if err != nil {
			log.Println("read", err)
			return
		}
		_, err = in.Write(buffer[:n])
		if err != nil {
			log.Println("write", err)
			return
		}
		//_, err := io.Copy(in, stream)
		////_, err := io.CopyBuffer(in, stream, buffer)
		if err != nil {
			log.Println(err)
			return
		}
	}
}
