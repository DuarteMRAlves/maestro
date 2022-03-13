package invoke

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryInvoke calls a unary grpc method. The context and args are passed to the
// grpc call. The reply structure should be a pre-allocated pointer and will
// be written by the call.
type UnaryInvoke func(context.Context, interface{}, interface{}) error

// NewUnaryInvoke creates a method to invoke unary grpc methods. It requires a
// fully qualified method name (which includes the service name) and a
// connection to use.
func NewUnaryInvoke(method string, conn grpc.ClientConnInterface) UnaryInvoke {
	return func(ctx context.Context, req interface{}, rep interface{}) error {
		err := conn.Invoke(ctx, method, req, rep)
		if err == nil {
			return nil
		}
		st, _ := status.FromError(err)
		return fmt.Errorf("invoke %s: %w", method, st.Err())
	}
}
