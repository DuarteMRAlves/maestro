package link

import "github.com/DuarteMRAlves/maestro/internal/domain"

type some struct{ domain.Link }

func (s some) IsError() bool { return false }

func (s some) Unwrap() domain.Link { return s.Link }

func (s some) Error() error { return nil }

type errResult struct{ error }

func (e errResult) IsError() bool { return true }

func (e errResult) Unwrap() domain.Link { panic("Stage not available in error result") }

func (e errResult) Error() error { return e.error }

func Some(l domain.Link) domain.LinkResult { return some{l} }

func Err(err error) domain.LinkResult { return errResult{err} }
