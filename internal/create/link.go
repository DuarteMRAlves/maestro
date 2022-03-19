package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
)

type LinkSaver interface {
	Save(internal.Link) error
}

type LinkLoader interface {
	Load(internal.LinkName) (internal.Link, error)
}

type LinkStorage interface {
	LinkSaver
	LinkLoader
}

var (
	EmptyLinkName        = errors.New("empty link name")
	EmptySourceStage     = errors.New("empty source stage")
	EmptyTargetStage     = errors.New("empty target stage")
	EqualSourceAndTarget = errors.New("equal source and target stages")
)

type StageNotInOrchestration struct {
	Orch  internal.OrchestrationName
	Stage internal.StageName
}

func (err *StageNotInOrchestration) Error() string {
	format := "stage '%s' not found in orchestration '%s'."
	return fmt.Sprintf(format, err.Stage, err.Orch)
}

type IncompatibleLinks struct {
	A, B string
}

func (err *IncompatibleLinks) Error() string {
	return fmt.Sprintf("incompatible links: %s, %s", err.A, err.B)
}

func Link(
	linkStorage LinkStorage,
	stageLoader StageLoader,
	orchStorage OrchestrationStorage,
) func(
	internal.LinkName,
	internal.LinkEndpoint,
	internal.LinkEndpoint,
	internal.OrchestrationName,
) error {
	return func(
		name internal.LinkName,
		source, target internal.LinkEndpoint,
		orchName internal.OrchestrationName,
	) error {
		if name.IsEmpty() {
			return EmptyLinkName
		}

		sourceStage := source.Stage()
		if sourceStage.IsEmpty() {
			return EmptySourceStage
		}

		targetStage := target.Stage()
		if targetStage.IsEmpty() {
			return EmptyTargetStage
		}

		if orchName.IsEmpty() {
			return EmptyOrchestrationName
		}

		_, err := linkStorage.Load(name)
		if err == nil {
			return &internal.AlreadyExists{Type: "link", Ident: name.Unwrap()}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return err
		}

		orch, err := orchStorage.Load(orchName)
		if err != nil {
			return fmt.Errorf("add link %s: %w", name, err)
		}

		_, err = stageLoader.Load(sourceStage)
		if err != nil {
			return err
		}

		_, err = stageLoader.Load(targetStage)
		if err != nil {
			return err
		}

		if !existsStage(sourceStage, orch.Stages()...) {
			return &StageNotInOrchestration{Orch: orchName, Stage: sourceStage}
		}
		if !existsStage(targetStage, orch.Stages()...) {
			return &StageNotInOrchestration{Orch: orchName, Stage: targetStage}
		}

		if sourceStage == targetStage {
			return EqualSourceAndTarget
		}

		targetLinks, err := linksForTarget(linkStorage, targetStage, orch.Links()...)
		if err != nil {
			return fmt.Errorf("create link %s: %w", name, err)
		}

		for _, l := range targetLinks {
			// 1. Target receives entire message from this link but another exists.
			// 2. Target already receives entire message from existing link.
			// 3. Target receives same field from both links.
			if target.Field().IsEmpty() ||
				l.Target().Field().IsEmpty() ||
				target.Field().Unwrap() == l.Target().Field().Unwrap() {
				return &IncompatibleLinks{A: name.Unwrap(), B: l.Name().Unwrap()}
			}
		}

		links := orch.Links()
		links = append(links, name)
		orch = internal.NewOrchestration(orch.Name(), orch.Stages(), links)

		err = orchStorage.Save(orch)
		if err != nil {
			return fmt.Errorf("add link %s: %w", name, err)
		}

		l := internal.NewLink(name, source, target)

		return linkStorage.Save(l)
	}
}

func existsStage(name internal.StageName, ss ...internal.StageName) bool {
	for _, s := range ss {
		if s == name {
			return true
		}
	}
	return false
}

// Retrieves the links that have a specific target stage.
func linksForTarget(
	linkLoader LinkLoader, target internal.StageName, links ...internal.LinkName,
) ([]internal.Link, error) {
	var ret []internal.Link
	for _, n := range links {
		l, err := linkLoader.Load(n)
		if err != nil {
			return nil, err
		}
		if l.Target().Stage() == target {
			ret = append(ret, l)
		}
	}
	return ret, nil
}
