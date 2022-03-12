package create

import (
	"github.com/DuarteMRAlves/maestro/internal"
)

type LinkResult interface {
	IsError() bool
	Unwrap() internal.Link
	Error() error
}

type someLink struct{ internal.Link }

func (s someLink) IsError() bool { return false }

func (s someLink) Unwrap() internal.Link { return s.Link }

func (s someLink) Error() error { return nil }

type errLink struct{ error }

func (e errLink) IsError() bool { return true }

func (e errLink) Unwrap() internal.Link { panic("Link not available in error result") }

func (e errLink) Error() error { return e.error }

func SomeLink(l internal.Link) LinkResult { return someLink{l} }

func ErrLink(err error) LinkResult { return errLink{err} }

func BindLink(f func(internal.Link) LinkResult) func(LinkResult) LinkResult {
	return func(result LinkResult) LinkResult {
		if result.IsError() {
			return result
		}
		return f(result.Unwrap())
	}
}

type OrchestrationResult interface {
	IsError() bool
	Unwrap() internal.Orchestration
	Error() error
}

type someOrchestration struct{ internal.Orchestration }

func (s someOrchestration) IsError() bool { return false }

func (s someOrchestration) Unwrap() internal.Orchestration { return s.Orchestration }

func (s someOrchestration) Error() error { return nil }

type errOrchestration struct{ error }

func (e errOrchestration) IsError() bool { return true }

func (e errOrchestration) Unwrap() internal.Orchestration {
	panic("Orchestration not available in error result")
}

func (e errOrchestration) Error() error { return e.error }

func SomeOrchestration(o internal.Orchestration) OrchestrationResult { return someOrchestration{o} }

func ErrOrchestration(err error) OrchestrationResult { return errOrchestration{err} }

func BindOrchestration(
	f func(internal.Orchestration) OrchestrationResult,
) func(OrchestrationResult) OrchestrationResult {
	return func(result OrchestrationResult) OrchestrationResult {
		if result.IsError() {
			return result
		}
		return f(result.Unwrap())
	}
}

func ReturnOrchestration(f func(internal.Orchestration) internal.Orchestration) func(internal.Orchestration) OrchestrationResult {
	return func(o internal.Orchestration) OrchestrationResult {
		return SomeOrchestration(f(o))
	}
}
