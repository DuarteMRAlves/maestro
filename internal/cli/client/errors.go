package client

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorFromGrpcError restores the error sent from the maestro server into
// the generic format defined in errdefs
func ErrorFromGrpcError(err error) error {
	if err == nil {
		return nil
	}
	st, _ := status.FromError(err)
	switch st.Code() {
	case codes.AlreadyExists:
		return errdefs.AlreadyExistsWithMsg(st.Message())
	case codes.NotFound:
		return errdefs.NotFoundWithMsg(st.Message())
	case codes.InvalidArgument:
		return errdefs.InvalidArgumentWithMsg(st.Message())
	case codes.Internal:
		return errdefs.InternalWithMsg(st.Message())
	case codes.Unknown:
	default:
		return errdefs.UnknownWithMsg(st.Message())
	}
	return errdefs.UnknownWithMsg(st.Message())
}
