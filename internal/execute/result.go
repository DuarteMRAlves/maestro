package execute

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
