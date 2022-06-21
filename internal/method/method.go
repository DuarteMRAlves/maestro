package method

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal/message"
)

// Conn represents an rpc connection with a remote method that can
// be used to call the method.
type Conn interface {
	Call(ctx context.Context, req message.Instance) (message.Instance, error)
	Close() error
}

// Dialer creates a connection with the method to execute
type Dialer interface {
	Dial() (Conn, error)
}

type DialFunc func() (Conn, error)

func (fn DialFunc) Dial() (Conn, error) { return fn() }

// Desc describes a method.
type Desc interface {
	Dialer
	Input() message.Type
	Output() message.Type
}

type Resolver interface {
	Resolve(ctx context.Context, address string) (Desc, error)
}

type ResolveFunc func(ctx context.Context, address string) (Desc, error)

func (fn ResolveFunc) Resolve(ctx context.Context, address string) (Desc, error) {
	return fn(ctx, address)
}
