package domain

type OptionalString interface {
	Unwrap() string
	Present() bool
}

type presentString struct{ string }

func (s presentString) Unwrap() string { return s.string }

func (s presentString) Present() bool { return true }

type emptyString struct{}

func (s emptyString) Unwrap() string {
	panic("String not available in empty optional")
}

func (s emptyString) Present() bool { return false }

func NewPresentString(data string) OptionalString {
	return presentString{data}
}

func NewEmptyString() OptionalString { return emptyString{} }
