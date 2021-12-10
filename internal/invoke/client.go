package invoke

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	return handleGrpcError(err)
}

func handleGrpcError(err error) error {
	if err == nil {
		return nil
	}
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.Unavailable, codes.Unimplemented:
		// Unavailable is for the case where maestro is not running. When a
		// stage is not running, it is a failed precondition.
		// Unimplemented is when the maestro server does not implement a given
		// method. When a stage does not have an implemented method, it is a
		// failed precondition.
		return errdefs.FailedPreconditionWithMsg("unary invoke: %v", st.Err())
	default:
		return errdefs.UnknownWithMsg("unary invoke: %v", st.Err())
	}
}
