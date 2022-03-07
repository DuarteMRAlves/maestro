package invoke

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleGrpcError(err error, prependMsg string) error {
	if err == nil {
		return nil
	}
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.Unavailable, codes.Unimplemented:
		// Unavailable is for the case where maestro is not running. When a
		// stage is not running, it is a failed precondition.
		// Unimplemented is when the maestro server does not implement a given
		// method. When a stage does not have an implemented method, in this
		// case reflection, it is a failed precondition.
		return errdefs.FailedPreconditionWithMsg("%v%v", prependMsg, st.Err())
	default:
		return errdefs.UnknownWithMsg("%s%s", prependMsg, st.Err())
	}
}

func isGrpcErr(err error) bool {
	_, ok := status.FromError(err)
	return ok
}

func isElementNotFoundErr(err error) bool {
	return grpcreflect.IsElementNotFoundError(err)
}

func isProtocolError(err error) bool {
	_, ok := err.(*grpcreflect.ProtocolError)
	return ok
}
