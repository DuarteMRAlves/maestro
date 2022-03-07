package invoke

type DynamicMessageResult interface {
	IsError() bool
	Unwrap() DynamicMessage
	Error() error
}

type someDynamicMessage struct{ DynamicMessage }

func (s someDynamicMessage) IsError() bool { return false }

func (s someDynamicMessage) Unwrap() DynamicMessage { return s.DynamicMessage }

func (s someDynamicMessage) Error() error { return nil }

type errDynamicMessage struct{ error }

func (e errDynamicMessage) IsError() bool { return true }

func (e errDynamicMessage) Unwrap() DynamicMessage {
	panic("DynamicMessage not available in error result")
}

func (e errDynamicMessage) Error() error { return e.error }

func SomeDynamicMessage(s DynamicMessage) DynamicMessageResult { return someDynamicMessage{s} }

func ErrDynamicMessage(err error) DynamicMessageResult { return errDynamicMessage{err} }
