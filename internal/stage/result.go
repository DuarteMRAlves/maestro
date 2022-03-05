package stage

import "github.com/DuarteMRAlves/maestro/internal/domain"

type some struct{ domain.Stage }

func (s some) IsError() bool { return false }

func (s some) Unwrap() domain.Stage { return s.Stage }

func (s some) Error() error { return nil }

type errResult struct{ error }

func (e errResult) IsError() bool { return true }

func (e errResult) Unwrap() domain.Stage { panic("Stage not available in error result") }

func (e errResult) Error() error { return e.error }

func Some(s domain.Stage) domain.StageResult { return some{s} }

func Err(err error) domain.StageResult { return errResult{err} }
