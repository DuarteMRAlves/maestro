package link

import "github.com/DuarteMRAlves/maestro/internal/types"

type some struct{ types.Link }

func (s some) IsError() bool { return false }

func (s some) Unwrap() types.Link { return s.Link }

func (s some) Error() error { return nil }

type errResult struct{ error }

func (e errResult) IsError() bool { return true }

func (e errResult) Unwrap() types.Link { panic("Stage not available in error result") }

func (e errResult) Error() error { return e.error }

func Some(l types.Link) types.LinkResult { return some{l} }

func Err(err error) types.LinkResult { return errResult{err} }
