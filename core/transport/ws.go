package transport

import (
	"Pier/common"
	"context"
	"github.com/gobwas/ws"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

const wsTimeout = 30

type WsTransport struct {
}

func (w *WsTransport) Listen(ctx context.Context, address string, handler func(stream io.ReadWriter)) error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}
		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(wsTimeout) * time.Second))
		handler(conn)
	})

	return http.ListenAndServe(address, nil)
}

func (w *WsTransport) Dail(ctx context.Context, address string, handler func(stream io.ReadWriter)) error {
	for common.Opened(ctx) {
		connCtx, connCancel := context.WithCancel(ctx)
		defer connCancel()
		conn := connectServer(address, connCtx)
		if conn == nil {
			connCancel()
			time.Sleep(3 * time.Second)
			continue
		}
		handler(conn)
		_ = conn.Close()
	}

	return nil
}

func connectServer(address string, ctx context.Context) net.Conn {
	scheme := "ws"

	u := url.URL{Scheme: scheme, Host: address, Path: "/ws"}
	header := make(http.Header)
	header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")

	dialer := ws.Dialer{
		Header:  ws.HandshakeHeaderHTTP(header),
		Timeout: time.Duration(wsTimeout) * time.Second,
		NetDial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(network, address)
		},
	}
	c, _, _, err := dialer.Dial(ctx, u.String())
	if err != nil {
		log.Printf("[client] failed to dial websocket %s %v", u.String(), err)
		return nil
	}
	return c
}
