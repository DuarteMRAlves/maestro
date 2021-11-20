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

func UnknownWithError(err error) error {
	if err == nil || IsUnknown(err) {
		return err
	}
	return unknown{err}
}

func UnknownWithMsg(format string, a ...interface{}) error {
	return UnknownWithError(fmt.Errorf(format, a...))
}
