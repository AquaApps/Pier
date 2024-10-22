package internal

import (
	"Pier/core/transport"
	"Pier/core/tunnel"
	"context"
	"io"
	"log"
)

func PierServer(Addr string, dev *tunnel.Device, appContext context.Context, transportImpl transport.Transport) error {
	out, in := dev.OpenChannel()
	return transportImpl.Listen(appContext, Addr, func(stream io.ReadWriter) {
		log.Println("Accept a connection.")
		connectCtx, cancelFunc := context.WithCancel(appContext)
		defer cancelFunc()
		go writer(stream, out, connectCtx)
		reader(stream, in, connectCtx)
	})
}
