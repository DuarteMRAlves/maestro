package domain

type OptionalError interface {
	Unwrap() error
	Present() bool
}

type presentError struct{ error }

func (e presentError) Unwrap() error { return e.error }

func (e presentError) Present() bool { return true }

type emptyError struct{}

func (e emptyError) Unwrap() error {
	panic("Error not available in empty optional")
}

func (e emptyError) Present() bool { return false }

func NewPresentError(data error) OptionalError { return presentError{data} }

func NewEmptyError() OptionalError { return emptyError{} }
