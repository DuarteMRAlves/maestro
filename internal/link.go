package internal

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"regexp"
)

var linkNameRegExp, _ = regexp.Compile(`^[a-zA-Z0-9]+([-:_/][a-zA-Z0-9]+)*$|^$`)

type LinkName struct{ val string }

func (s LinkName) Unwrap() string { return s.val }

func (s LinkName) IsEmpty() bool { return s.val == "" }

func NewLinkName(name string) (LinkName, error) {
	if !isValidLinkName(name) {
		err := errdefs.InvalidArgumentWithMsg("invalid name '%v'", name)
		return LinkName{}, err
	}
	return LinkName{val: name}, nil
}

func isValidLinkName(name string) bool {
	return linkNameRegExp.MatchString(name)
}

type MessageField struct{ val string }

func (m MessageField) Unwrap() string { return m.val }

func (m MessageField) IsEmpty() bool { return m.val == "" }

func NewMessageField(field string) MessageField {
	return MessageField{val: field}
}

type OptionalMessageField struct {
	val     MessageField
	present bool
}

func (p OptionalMessageField) Unwrap() MessageField { return p.val }

func (p OptionalMessageField) Present() bool { return p.present }

func NewPresentMessageField(m MessageField) OptionalMessageField {
	return OptionalMessageField{val: m, present: true}
}

func NewEmptyMessageField() OptionalMessageField {
	return OptionalMessageField{}
}

type LinkEndpoint struct {
	stage StageName
	field OptionalMessageField
}

func (e LinkEndpoint) Stage() StageName {
	return e.stage
}

func (e LinkEndpoint) Field() OptionalMessageField {
	return e.field
}

func NewLinkEndpoint(
	stage StageName,
	field OptionalMessageField,
) LinkEndpoint {
	return LinkEndpoint{
		stage: stage,
		field: field,
	}
}

type Link struct {
	name          LinkName
	source        LinkEndpoint
	target        LinkEndpoint
	orchestration domain.OrchestrationName
}

func (l Link) Name() LinkName {
	return l.name
}

func (l Link) Source() LinkEndpoint {
	return l.source
}

func (l Link) Target() LinkEndpoint {
	return l.target
}

func (l Link) Orchestration() domain.OrchestrationName {
	return l.orchestration
}

func NewLink(
	name LinkName,
	source, target LinkEndpoint,
	orchestration domain.OrchestrationName,
) Link {
	return Link{
		name:          name,
		source:        source,
		target:        target,
		orchestration: orchestration,
	}
}
