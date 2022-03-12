package create

import (
	"github.com/DuarteMRAlves/maestro/internal"
)

type LinkResult interface {
	IsError() bool
	Unwrap() internal.Link
	Error() error
}

type someLink struct{ internal.Link }

func (s someLink) IsError() bool { return false }

func (s someLink) Unwrap() internal.Link { return s.Link }

func (s someLink) Error() error { return nil }

type errLink struct{ error }

func (e errLink) IsError() bool { return true }

func (e errLink) Unwrap() internal.Link { panic("Link not available in error result") }

func (e errLink) Error() error { return e.error }

func SomeLink(l internal.Link) LinkResult { return someLink{l} }

func ErrLink(err error) LinkResult { return errLink{err} }

func BindLink(f func(internal.Link) LinkResult) func(LinkResult) LinkResult {
	return func(result LinkResult) LinkResult {
		if result.IsError() {
			return result
		}
		return f(result.Unwrap())
	}
}
