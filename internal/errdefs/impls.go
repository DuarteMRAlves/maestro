package errdefs

type notFound struct {
	error
}

func (e notFound) NotFound() { /* Do nothing */ }

func (e notFound) Cause() error {
	return e.error
}

type alreadyExists struct {
	error
}

func (e alreadyExists) AlreadyExists() { /* Do nothing */ }

func (e alreadyExists) Cause() error {
	return e.error
}

type invalidArgument struct {
	error
}

func (e invalidArgument) InvalidArgument() { /* Do nothing */ }

func (e invalidArgument) Cause() error {
	return e.error
}

type internal struct {
	error
}

func (e internal) Internal() { /* Do nothing */ }

func (e internal) Cause() error {
	return e.error
}

type unknown struct {
	error
}

func (e unknown) Unknown() { /* Do nothing */ }

func (e unknown) Cause() error {
	return e.error
}
