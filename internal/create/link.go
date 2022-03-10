package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
)

type LinkSaver interface {
	Save(Link) LinkResult
}

type LinkLoader interface {
	Load(domain.LinkName) LinkResult
}

type LinkStorage interface {
	LinkSaver
	LinkLoader
}

type LinkEndpoint interface {
	Stage() domain.StageName
	Field() domain.OptionalMessageField
}

type Link interface {
	Name() domain.LinkName
	Source() LinkEndpoint
	Target() LinkEndpoint
	Orchestration() domain.OrchestrationName
}

type LinkRequest struct {
	Name string

	SourceStage string
	SourceField domain.OptionalString
	TargetStage string
	TargetField domain.OptionalString

	Orchestration string
}

type LinkResponse struct {
	Err domain.OptionalError
}

// Implementation of interfaces

type linkEndpoint struct {
	stage domain.StageName
	field domain.OptionalMessageField
}

func (e linkEndpoint) Stage() domain.StageName {
	return e.stage
}

func (e linkEndpoint) Field() domain.OptionalMessageField {
	return e.field
}

func NewLinkEndpoint(
	stage domain.StageName,
	field domain.OptionalMessageField,
) LinkEndpoint {
	return linkEndpoint{
		stage: stage,
		field: field,
	}
}

type link struct {
	name          domain.LinkName
	source        LinkEndpoint
	target        LinkEndpoint
	orchestration domain.OrchestrationName
}

func (l link) Name() domain.LinkName {
	return l.name
}

func (l link) Source() LinkEndpoint {
	return l.source
}

func (l link) Target() LinkEndpoint {
	return l.target
}

func (l link) Orchestration() domain.OrchestrationName {
	return l.orchestration
}

func NewLink(
	name domain.LinkName,
	source, target LinkEndpoint,
	orchestration domain.OrchestrationName,
) Link {
	return link{
		name:          name,
		source:        source,
		target:        target,
		orchestration: orchestration,
	}
}
