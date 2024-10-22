package transport

import (
	"Pier/common"
	"context"
	"io"
	"net"
)

type TcpTransport struct {
}

func (t *TcpTransport) Listen(ctx context.Context, address string, handler func(stream io.ReadWriter)) error {
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	for common.Opened(ctx) {
		if stream, err := listener.Accept(); err != nil {
			continue
		} else {
			handler(stream)
		}
	}
	return nil
}

func (t *TcpTransport) Dail(ctx context.Context, address string, handler func(stream io.ReadWriter)) error {
	for common.Opened(ctx) {
		stream, err := net.Dial("tcp", address)
		if err != nil {
			return err
		}
		handler(stream)
	}
	return nil
}
