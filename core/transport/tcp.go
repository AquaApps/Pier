package transport

import (
	"Pier/common"
	"context"
	"io"
	"net"
	"time"
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
			if tcpConn, ok := stream.(*net.TCPConn); ok {
				tcpConn.SetNoDelay(true)
				tcpConn.SetKeepAlive(true)
				tcpConn.SetKeepAlivePeriod(30 * time.Second)
				tcpConn.SetWriteBuffer(640 * 1024)
				tcpConn.SetReadBuffer(640 * 1024)
			}
			go handler(stream)
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
		if tcpConn, ok := stream.(*net.TCPConn); ok {
			tcpConn.SetNoDelay(true)
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(30 * time.Second)
			tcpConn.SetWriteBuffer(640 * 1024)
			tcpConn.SetReadBuffer(640 * 1024)
		}
		handler(stream)
	}
	return nil
}
