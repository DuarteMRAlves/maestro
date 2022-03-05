package domain

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

type OrchestrationResult interface {
	IsError() bool
	Unwrap() Orchestration
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

func Bind(f func(Asset) AssetResult) func(AssetResult) AssetResult {
	return func(resultAsset AssetResult) AssetResult {
		if resultAsset.IsError() {
			return resultAsset
		}
		return f(resultAsset.Unwrap())
	}
}

type someStage struct{ Stage }

func (s someStage) IsError() bool { return false }

func (s someStage) Unwrap() Stage { return s.Stage }

func (s someStage) Error() error { return nil }

type errStage struct{ error }

func (e errStage) IsError() bool { return true }

func (e errStage) Unwrap() Stage { panic("Stage not available in error result") }

func (e errStage) Error() error { return e.error }

func SomeStage(s Stage) StageResult { return someStage{s} }

func ErrStage(err error) StageResult { return errStage{err} }

type someLink struct{ Link }

func (s someLink) IsError() bool { return false }

func (s someLink) Unwrap() Link { return s.Link }

func (s someLink) Error() error { return nil }

type errLink struct{ error }

func (e errLink) IsError() bool { return true }

func (e errLink) Unwrap() Link { panic("Link not available in error result") }

func (e errLink) Error() error { return e.error }

func SomeLink(l Link) LinkResult { return someLink{l} }

func ErrLink(err error) LinkResult { return errLink{err} }

type someOrchestration struct{ Orchestration }

func (s someOrchestration) IsError() bool { return false }

func (s someOrchestration) Unwrap() Orchestration { return s.Orchestration }

func (s someOrchestration) Error() error { return nil }

type errOrchestration struct{ error }

func (e errOrchestration) IsError() bool { return true }

func (e errOrchestration) Unwrap() Orchestration {
	panic("Orchestration not available in error result")
}

func (e errOrchestration) Error() error { return e.error }

func SomeOrchestration(o Orchestration) OrchestrationResult { return someOrchestration{o} }

func ErrOrchestration(err error) OrchestrationResult { return errOrchestration{err} }
