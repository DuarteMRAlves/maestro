package errdefs

import "fmt"

func NotFoundWithError(err error) error {
	if err == nil || IsNotFound(err) {
		return err
	}
	return notFound{err}
}

func NotFoundWithMsg(format string, a ...interface{}) error {
	return NotFoundWithError(fmt.Errorf(format, a...))
}

func AlreadyExistsWithError(err error) error {
	if err == nil || IsAlreadyExists(err) {
		return err
	}
	return alreadyExists{err}
}

func AlreadyExistsWithMsg(format string, a ...interface{}) error {
	return AlreadyExistsWithError(fmt.Errorf(format, a...))
}

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

func UnknownWithError(err error) error {
	if err == nil || IsUnknown(err) {
		return err
	}
	return unknown{err}
}

func UnknownWithMsg(format string, a ...interface{}) error {
	return UnknownWithError(fmt.Errorf(format, a...))
}
