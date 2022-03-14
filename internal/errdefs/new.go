package errdefs

import (
	"fmt"
)

func InvalidArgumentWithError(err error) error {
	if err == nil || IsInvalidArgument(err) {
		return err
	}
	return invalidArgument{err}
}

func InvalidArgumentWithMsg(format string, a ...interface{}) error {
	return InvalidArgumentWithError(fmt.Errorf(format, a...))
}

func InternalWithError(err error) error {
	if err == nil || IsInternal(err) {
		return err
	}
	return internal{err}
}

func InternalWithMsg(format string, a ...interface{}) error {
	return InternalWithError(fmt.Errorf(format, a...))
}

// PrependMsg returns a new error of the same type of the received error that
// prepends a message to the received error message
func PrependMsg(err error, format string, a ...interface{}) error {
	format = fmt.Sprintf("%s: %%w", format)
	a = append(a, err)

	causeErr := getImplementer(err)
	switch causeErr.(type) {
	case InvalidArgument:
		return InvalidArgumentWithMsg(format, a...)
	case Internal:
		return InternalWithMsg(format, a...)
	}
	return fmt.Errorf(format, a...)
}
