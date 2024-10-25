package transport

import (
	"context"
	"io"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func fillRandomBytes(packet []byte) int {
	for i := range packet {
		packet[i] = byte(rand.Intn(256))
	}
	return len(packet)
}

var appContext = context.Background()
var transportImpl = new(TcpTransport)

var read uint64
var weite uint64

func TestTcp_Listen(t *testing.T) {
	transportImpl.Listen(appContext, ":38324", func(stream io.ReadWriter) {
		go writer(stream)
		reader(stream)
	})
}

func TestTcp_Dail(t *testing.T) {
	transportImpl.Dail(appContext, "192.168.222.66:38324", func(stream io.ReadWriter) {
		go writer(stream)
		reader(stream)
	})
}

func writer(stream io.Writer) {
	packet := make([]byte, 64*1024)
	for {
		length := fillRandomBytes(packet)
		_, _ = stream.Write(packet[:length])
		weite += uint64(length)
	}
}

func reader(stream io.Reader) {
	packet := make([]byte, 64*1024)
	for {
		num, err := stream.Read(packet)
		if err != nil {
			return
		}
		read += uint64(num)
	}
}
