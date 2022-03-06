package domain

type AssetResult interface {
	IsError() bool
	Unwrap() Asset
	Error() error
}

type someAsset struct{ Asset }

func (s someAsset) IsError() bool { return false }

func (s someAsset) Unwrap() Asset { return s.Asset }

func (s someAsset) Error() error { return nil }

type errAsset struct{ error }

func (e errAsset) IsError() bool { return true }

func (e errAsset) Unwrap() Asset { panic("Asset not available in error") }

func (e errAsset) Error() error { return e.error }

func SomeAsset(a Asset) AssetResult { return someAsset{a} }

func ErrAsset(err error) AssetResult { return errAsset{err} }

func BindAsset(f func(Asset) AssetResult) func(AssetResult) AssetResult {
	return func(resultAsset AssetResult) AssetResult {
		if resultAsset.IsError() {
			return resultAsset
		}
		return f(resultAsset.Unwrap())
	}
}
