package errdefs

type internal struct {
	error
}

func (e internal) Internal() { /* Do nothing */ }

func (e internal) Cause() error {
	return e.error
}
