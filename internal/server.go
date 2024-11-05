package internal

import (
	"Pier/core/transport"
	"context"
	pun "github.com/AquaApps/Pun"

	//pun "github.com/AquaApps/Pun"
	"io"
	"log"
)

func PierServer(Addr string, dev *pun.Device, appContext context.Context, transportImpl transport.Transport) error {
	out, in := dev.OpenStream()
	first := false
	return transportImpl.Listen(appContext, Addr, func(stream io.ReadWriter) {
		if first == false {
			first = true
			log.Println("Accept a connection.")
			connectCtx, cancelFunc := context.WithCancel(appContext)
			defer cancelFunc()
			go writer(stream, out, connectCtx)
			reader(stream, in, connectCtx)
		} else {
			log.Println("Accept a extra connection.")
			connectCtx, cancelFunc := context.WithCancel(appContext)
			defer cancelFunc()
			extraStream, err := dev.OpenExtraStream()
			if err != nil {
				log.Println("Extra connection error.", err)
				return
			}
			defer extraStream.Close()

			go writer(stream, extraStream.OutputStream, connectCtx)
			reader(stream, extraStream.InputStream, connectCtx)
		}

	})
}
