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
		num, _ := out.Read(buffer)
		_, _ = stream.Write(buffer[:num])
	}
}

func reader(stream io.Reader, in io.Writer, ctx context.Context) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	buffer := make([]byte, 4*1024)
	for common.Opened(ctx) {
		_, _ = stream.Read(buffer)

		_, err := in.Write(buffer)
		if err != nil {
			log.Println("read", err)
			return
		}
		//n, _ := device._f.Write(fun(<-device._inputStream))
		//device.incrReadBytes(n)
		//n, err := stream.Read(buffer)
		//if err != nil {
		//	log.Println("read", err)
		//	return
		//}
		//_, err = in.Write(buffer[:n])
		//if err != nil {
		//	log.Println("write", err)
		//	return
		//}
		////_, err := io.Copy(in, stream)
		//////_, err := io.CopyBuffer(in, stream, buffer)
		//if err != nil {
		//	log.Println(err)
		//	return
		//}
	}
}
