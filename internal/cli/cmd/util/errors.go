package util

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

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