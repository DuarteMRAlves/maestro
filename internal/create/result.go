package create

type StageResult interface {
	IsError() bool
	Unwrap() Stage
	Error() error
}

type OrchestrationResult interface {
	IsError() bool
	Unwrap() Orchestration
	Error() error
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

func BindStage(f func(Stage) StageResult) func(StageResult) StageResult {
	return func(result StageResult) StageResult {
		if result.IsError() {
			return result
		}
		return f(result.Unwrap())
	}
}

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

func BindOrchestration(
	f func(Orchestration) OrchestrationResult,
) func(OrchestrationResult) OrchestrationResult {
	return func(result OrchestrationResult) OrchestrationResult {
		if result.IsError() {
			return result
		}
		return f(result.Unwrap())
	}
}
