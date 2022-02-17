package util

import (
	"fmt"
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
	case codes.FailedPrecondition:
		return errdefs.FailedPreconditionWithMsg(st.Message())
	case codes.Unavailable:
		return errdefs.UnavailableWithMsg(st.Message())
	case codes.Internal:
		return errdefs.InternalWithMsg(st.Message())
	case codes.Unknown:
		return errdefs.UnknownWithMsg(st.Message())
	default:
		return errdefs.UnknownWithMsg(st.Message())
	}
}

func DisplayMsgFromError(err error) string {
	switch {
	case errdefs.IsAlreadyExists(err):
		return fmt.Sprintf("already exists: %v", err.Error())
	case errdefs.IsNotFound(err):
		return fmt.Sprintf("not found: %v", err.Error())
	case errdefs.IsInvalidArgument(err):
		return fmt.Sprintf("invalid argument: %v", err.Error())
	case errdefs.IsFailedPrecondition(err):
		return fmt.Sprintf("failed precondition: %v", err.Error())
	case errdefs.IsUnavailable(err):
		return fmt.Sprintf("unavailable: %v", err.Error())
	case errdefs.IsInternal(err):
		return fmt.Sprintf("internal error: %v", err.Error())
	case errdefs.IsUnknown(err):
		return fmt.Sprintf("unknwon error: %v", err.Error())
	default:
		return fmt.Sprintf("unknwon error: %v", err.Error())
	}
}
