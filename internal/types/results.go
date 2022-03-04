package types

type AssetResult interface {
	IsError() bool
	Unwrap() Asset
	Error() error
}

type StageResult interface {
	IsError() bool
	Unwrap() Stage
	Error() error
}
