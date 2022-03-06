package execute

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
)

type linkEndpoint struct {
	stage Stage
	field domain.OptionalMessageField
}

func (e *linkEndpoint) Stage() Stage {
	return e.stage
}

func (e *linkEndpoint) Field() domain.OptionalMessageField {
	return e.field
}

func NewLinkEndpoint(
	stage Stage,
	field domain.OptionalMessageField,
) LinkEndpoint {
	return &linkEndpoint{
		stage: stage,
		field: field,
	}
}

type link struct {
	name   domain.LinkName
	source LinkEndpoint
	target LinkEndpoint
}

func (l *link) Name() domain.LinkName {
	return l.name
}

func (l *link) Source() LinkEndpoint {
	return l.source
}

func (l *link) Target() LinkEndpoint {
	return l.target
}

func NewLink(
	name domain.LinkName,
	source LinkEndpoint,
	target LinkEndpoint,
) Link {
	return &link{
		name:   name,
		source: source,
		target: target,
	}
}
