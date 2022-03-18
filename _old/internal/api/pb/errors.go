package pb

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcErrorFromError(err error) error {
	switch {
	case errdefs.IsAlreadyExists(err):
		return status.Error(codes.AlreadyExists, err.Error())
	case errdefs.IsNotFound(err):
		return status.Error(codes.NotFound, err.Error())
	case errdefs.IsInvalidArgument(err):
		return status.Error(codes.InvalidArgument, err.Error())
	case errdefs.IsInternal(err):
		return status.Error(codes.Internal, "internal server error")
	case errdefs.IsUnknown(err):
	default:
		return status.Error(codes.Unknown, err.Error())
	}
	// Should never happen
	return status.Error(codes.Internal, err.Error())
}
