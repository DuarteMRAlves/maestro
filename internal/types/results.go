package types

type AssetResult interface {
	IsError() bool
	Unwrap() Asset
	Error() error
}
