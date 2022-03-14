package errdefs

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
