package link

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/types"
	"regexp"
)

var nameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$`)

type linkName string

func NewLinkName(name string) (types.LinkName, error) {
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

func NewMessageField(field string) (types.MessageField, error) {
	if len(field) == 0 {
		return nil, errdefs.InvalidArgumentWithMsg("empty field")
	}
	return messageField(field), nil
}

type presentMessageField struct{ types.MessageField }

func (p presentMessageField) Unwrap() types.MessageField { return p.MessageField }

func (p presentMessageField) Present() bool { return true }

type emptyMessageField struct{}

func (e emptyMessageField) Unwrap() types.MessageField {
	panic("Message Field not available in an empty optional")
}

func (e emptyMessageField) Present() bool { return false }

func NewPresentMessageField(m types.MessageField) types.OptionalMessageField {
	return presentMessageField{m}
}

func NewEmptyMessageField() types.OptionalMessageField {
	return emptyMessageField{}
}

type linkEndpoint struct {
	stage types.Stage
	field types.OptionalMessageField
}

func (e *linkEndpoint) Stage() types.Stage {
	return e.stage
}

func (e *linkEndpoint) Field() types.OptionalMessageField {
	return e.field
}

func NewLinkEndpoint(
	stage types.Stage,
	field types.OptionalMessageField,
) types.LinkEndpoint {
	return &linkEndpoint{
		stage: stage,
		field: field,
	}
}

type link struct {
	name   types.LinkName
	source types.LinkEndpoint
	target types.LinkEndpoint
}

func (l *link) Name() types.LinkName {
	return l.name
}

func (l *link) Source() types.LinkEndpoint {
	return l.source
}

func (l *link) Target() types.LinkEndpoint {
	return l.target
}

func NewLink(
	name types.LinkName,
	source types.LinkEndpoint,
	target types.LinkEndpoint,
) types.Link {
	return &link{
		name:   name,
		source: source,
		target: target,
	}
}
