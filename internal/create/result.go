package create

type OrchestrationResult interface {
	IsError() bool
	Unwrap() Orchestration
	Error() error
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
