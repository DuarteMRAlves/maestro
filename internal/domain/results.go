package domain

type AssetResult interface {
	IsError() bool
	Unwrap() Asset
	Error() error
}
