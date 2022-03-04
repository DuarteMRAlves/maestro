package rpc

import (
	"context"
	"google.golang.org/grpc"
)

// UnaryClient exposes an API to invoke unary grpc methods. It requires a
// fully qualified method name (which includes the service name) and a
// connection to use.
type UnaryClient interface {
	Invoke(ctx context.Context, args interface{}, reply interface{}) error
}

type unaryClient struct {
	method string
	conn   grpc.ClientConnInterface
}

func NewUnary(method string, conn grpc.ClientConnInterface) UnaryClient {
	return &unaryClient{
		method: method,
		conn:   conn,
	}
}

// Invoke calls a unary grpc method. The context and args are passed to the
// grpc call. The reply structure should be a pre-allocated pointer and will
// be written by the call.
func (u *unaryClient) Invoke(
	ctx context.Context,
	args interface{},
	reply interface{},
) error {
	err := u.conn.Invoke(ctx, u.method, args, reply)
	return handleGrpcError(err, "unary invoke: ")
}
