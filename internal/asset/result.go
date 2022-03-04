package asset

import "github.com/DuarteMRAlves/maestro/internal/types"

type success struct {
	types.Asset
}

func (s success) IsError() bool { return false }

func (s success) Unwrap() types.Asset { return s.Asset }

func (s success) Error() error { return nil }

type errResult struct {
	error
}

func (e errResult) IsError() bool { return true }

func (e errResult) Unwrap() types.Asset { panic("Asset not available in error") }

func (e errResult) Error() error { return e.error }

func NewResult(a types.Asset) types.AssetResult { return success{a} }

func NewErrResult(err error) types.AssetResult { return errResult{err} }

func Bind(f func(types.Asset) types.AssetResult) func(types.AssetResult) types.AssetResult {
	return func(resultAsset types.AssetResult) types.AssetResult {
		if resultAsset.IsError() {
			return resultAsset
		}
		return f(resultAsset.Unwrap())
	}
}

func BindNoErr(f func(types.Asset) types.Asset) func(types.AssetResult) types.AssetResult {
	return Bind(
		func(d types.Asset) types.AssetResult {
			return NewResult(f(d))
		},
	)
}

func Compose(funcs ...func(types.AssetResult) types.AssetResult) func(types.AssetResult) types.AssetResult {
	return func(r types.AssetResult) types.AssetResult {
		for _, f := range funcs {
			r = f(r)
		}
		return r
	}
}
