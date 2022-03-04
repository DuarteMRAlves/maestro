package stage

import "github.com/DuarteMRAlves/maestro/internal/types"

type some struct{ types.Stage }

func (s some) IsError() bool { return false }

func (s some) Unwrap() types.Stage { return s.Stage }

func (s some) Error() error { return nil }

type errResult struct{ error }

func (e errResult) IsError() bool { return true }

func (e errResult) Unwrap() types.Stage { panic("Stage not available in error result") }

func (e errResult) Error() error { return e.error }

func Some(s types.Stage) types.StageResult { return some{s} }

func Err(err error) types.StageResult { return errResult{err} }
