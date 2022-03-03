package asset

import "github.com/DuarteMRAlves/maestro/internal/domain"

type Result interface {
	IsError() bool
	Unwrap() domain.Asset
	Error() error
}

type success struct {
	domain.Asset
}

func (s success) IsError() bool { return false }

func (s success) Unwrap() domain.Asset { return s.Asset }

func (s success) Error() error { return nil }

type errResult struct {
	error
}

func (e errResult) IsError() bool { return true }

func (e errResult) Unwrap() domain.Asset { panic("Asset not available in error") }

func (e errResult) Error() error { return e.error }

func NewResult(a domain.Asset) Result { return success{a} }

func NewErrResult(err error) Result { return errResult{err} }

func Bind(f func(domain.Asset) Result) func(Result) Result {
	return func(resultAsset Result) Result {
		if resultAsset.IsError() {
			return resultAsset
		}
		return f(resultAsset.Unwrap())
	}
}

func BindNoErr(f func(domain.Asset) domain.Asset) func(Result) Result {
	return Bind(
		func(d domain.Asset) Result {
			return NewResult(f(d))
		},
	)
}

func Compose(funcs ...func(Result) Result) func(Result) Result {
	return func(r Result) Result {
		for _, f := range funcs {
			r = f(r)
		}
		return r
	}
}
