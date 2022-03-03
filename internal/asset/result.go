package asset

import "github.com/DuarteMRAlves/maestro/internal/domain"

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

func NewResult(a domain.Asset) domain.AssetResult { return success{a} }

func NewErrResult(err error) domain.AssetResult { return errResult{err} }

func Bind(f func(domain.Asset) domain.AssetResult) func(domain.AssetResult) domain.AssetResult {
	return func(resultAsset domain.AssetResult) domain.AssetResult {
		if resultAsset.IsError() {
			return resultAsset
		}
		return f(resultAsset.Unwrap())
	}
}

func BindNoErr(f func(domain.Asset) domain.Asset) func(domain.AssetResult) domain.AssetResult {
	return Bind(
		func(d domain.Asset) domain.AssetResult {
			return NewResult(f(d))
		},
	)
}

func Compose(funcs ...func(domain.AssetResult) domain.AssetResult) func(domain.AssetResult) domain.AssetResult {
	return func(r domain.AssetResult) domain.AssetResult {
		for _, f := range funcs {
			r = f(r)
		}
		return r
	}
}
