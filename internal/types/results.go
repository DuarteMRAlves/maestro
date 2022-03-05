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

type LinkResult interface {
	IsError() bool
	Unwrap() Link
	Error() error
}
