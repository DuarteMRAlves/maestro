package domain

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var linkNameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type linkName string

func NewLinkName(name string) (LinkName, error) {
	if len(name) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty name")
	}
	if !isValidLinkName(name) {
		return nil, errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
	}
	return linkName(name), nil
}

func (s linkName) Unwrap() string {
	return string(s)
}

func isValidLinkName(name string) bool {
	return linkNameRegExp.MatchString(name)
}

type messageField string

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

type linkEndpoint struct {
	stage Stage
	field OptionalMessageField
}

func (e *linkEndpoint) Stage() Stage {
	return e.stage
}

func (e *linkEndpoint) Field() OptionalMessageField {
	return e.field
}

func NewLinkEndpoint(
	stage Stage,
	field OptionalMessageField,
) LinkEndpoint {
	return &linkEndpoint{
		stage: stage,
		field: field,
	}
}

type link struct {
	name   LinkName
	source LinkEndpoint
	target LinkEndpoint
}

func (l *link) Name() LinkName {
	return l.name
}

func (l *link) Source() LinkEndpoint {
	return l.source
}

func (l *link) Target() LinkEndpoint {
	return l.target
}

func NewLink(
	name LinkName,
	source LinkEndpoint,
	target LinkEndpoint,
) Link {
	return &link{
		name:   name,
		source: source,
		target: target,
	}
}
