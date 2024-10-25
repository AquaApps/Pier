package tunnel

import (
	"Pier/common"
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func fillRandomBytes(packet []byte) int {
	for i := range packet {
		packet[i] = byte(rand.Intn(256))
	}
	return len(packet)
}

var total uint64 = 0

func BenchmarkTunnel(b *testing.B) {
	background := context.Background()
	for i := 1; i < 2; i++ {
		tun := new(Device)
		err := tun.Init(RandomString(10), fmt.Sprintf("10.10.10.%d/24", i), background)
		if err != nil {
			b.Fatal(err)
			return
		}

		out, in := tun.OpenChannel()
		defer tun.Destroy()
		ctx, cancelFunc := context.WithCancel(background)

		go func(out <-chan []byte, ctx context.Context) {
			for common.Opened(ctx) {
				<-out
			}
		}(out, ctx)
		go func(can context.CancelFunc) {
			timer := time.NewTimer(10 * time.Second)
			<-timer.C
			b.Logf("index %d: speed: %dkB/s", i, total/1024/10)
			can()
		}(cancelFunc)
		packet := make([]byte, 4*1024)
		for common.Opened(ctx) {
			length := fillRandomBytes(packet)
			in <- packet[:length]
			total += uint64(length)
		}
	}
}
