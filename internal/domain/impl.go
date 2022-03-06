package domain

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var linkNameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type linkName string

func (s linkName) LinkName() {}

func (s linkName) Unwrap() string {
	return string(s)
}

func NewLinkName(name string) (LinkName, error) {
	if len(name) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty name")
	}
	if !isValidLinkName(name) {
		return nil, errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
	}
	return linkName(name), nil
}

func isValidLinkName(name string) bool {
	return linkNameRegExp.MatchString(name)
}

type messageField string

func (m messageField) MessageField() {}

func (m messageField) Unwrap() string {
	return string(m)
}

func NewMessageField(field string) (MessageField, error) {
	if len(field) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty field")
	}
	return messageField(field), nil
}

type presentMessageField struct{ MessageField }

func (p presentMessageField) Unwrap() MessageField { return p.MessageField }

func (p presentMessageField) Present() bool { return true }

type emptyMessageField struct{}

func (e emptyMessageField) Unwrap() MessageField {
	panic("Message Field not available in an empty optional")
}

func (e emptyMessageField) Present() bool { return false }

func NewPresentMessageField(m MessageField) OptionalMessageField {
	return presentMessageField{m}
}

func NewEmptyMessageField() OptionalMessageField {
	return emptyMessageField{}
}

var orchestrationNameReqExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type orchestrationName string

func (o orchestrationName) OrchestrationName() {}

func (o orchestrationName) Unwrap() string {
	return string(o)
}

func NewOrchestrationName(name string) (OrchestrationName, error) {
	if isValidOrchestrationName(name) {
		return orchestrationName(name), nil
	}
	return nil, errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
}

func isValidOrchestrationName(name string) bool {
	return orchestrationNameReqExp.MatchString(name)
}
