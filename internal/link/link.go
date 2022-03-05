package link

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type linkName string

func NewLinkName(name string) (domain.LinkName, error) {
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
	return nameRegExp.MatchString(name)
}

type messageField string

func (m messageField) Unwrap() string {
	return string(m)
}

func NewMessageField(field string) (domain.MessageField, error) {
	if len(field) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty field")
	}
	return messageField(field), nil
}

type presentMessageField struct{ domain.MessageField }

func (p presentMessageField) Unwrap() domain.MessageField { return p.MessageField }

func (p presentMessageField) Present() bool { return true }

type emptyMessageField struct{}

func (e emptyMessageField) Unwrap() domain.MessageField {
	panic("Message Field not available in an empty optional")
}

func (e emptyMessageField) Present() bool { return false }

func NewPresentMessageField(m domain.MessageField) domain.OptionalMessageField {
	return presentMessageField{m}
}

func NewEmptyMessageField() domain.OptionalMessageField {
	return emptyMessageField{}
}

type linkEndpoint struct {
	stage domain.Stage
	field domain.OptionalMessageField
}

func (e *linkEndpoint) Stage() domain.Stage {
	return e.stage
}

func (e *linkEndpoint) Field() domain.OptionalMessageField {
	return e.field
}

func NewLinkEndpoint(
	stage domain.Stage,
	field domain.OptionalMessageField,
) domain.LinkEndpoint {
	return &linkEndpoint{
		stage: stage,
		field: field,
	}
}

type link struct {
	name   domain.LinkName
	source domain.LinkEndpoint
	target domain.LinkEndpoint
}

func (l *link) Name() domain.LinkName {
	return l.name
}

func (l *link) Source() domain.LinkEndpoint {
	return l.source
}

func (l *link) Target() domain.LinkEndpoint {
	return l.target
}

func NewLink(
	name domain.LinkName,
	source domain.LinkEndpoint,
	target domain.LinkEndpoint,
) domain.Link {
	return &link{
		name:   name,
		source: source,
		target: target,
	}
}
