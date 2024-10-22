package transport

import (
	"context"
	"io"
)

type Transport interface {
	Listen(ctx context.Context, address string, handler func(stream io.ReadWriter)) error
	Dail(ctx context.Context, address string, handler func(stream io.ReadWriter)) error
}
