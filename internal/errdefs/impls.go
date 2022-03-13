package errdefs

type invalidArgument struct {
	error
}

func (e invalidArgument) InvalidArgument() { /* Do nothing */ }

func (e invalidArgument) Cause() error {
	return e.error
}

type unavailable struct {
	error
}

func (e unavailable) Unavailable() { /* Do nothing */ }

func (e unavailable) Cause() error {
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
