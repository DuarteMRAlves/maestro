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
	EmptySourceField     = errors.New("empty source field")
	EmptyTargetStage     = errors.New("empty target stage")
	EmptyTargetField     = errors.New("empty target field")
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

func Link(
	storage LinkStorage,
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
		if source.Field().Present() && source.Field().Unwrap().IsEmpty() {
			return EmptySourceField
		}

		targetStage := target.Stage()
		if targetStage.IsEmpty() {
			return EmptyTargetStage
		}
		if target.Field().Present() && target.Field().Unwrap().IsEmpty() {
			return EmptyTargetField
		}

		if orchName.IsEmpty() {
			return EmptyOrchestrationName
		}

		_, err := storage.Load(name)
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

		foundTarget := false
		foundSource := false
		for _, s := range orch.Stages() {
			if s == sourceStage {
				foundSource = true
			} else if s == targetStage {
				foundTarget = true
			}
		}
		if !foundSource {
			return &StageNotInOrchestration{Orch: orchName, Stage: sourceStage}
		}
		if !foundTarget {
			return &StageNotInOrchestration{Orch: orchName, Stage: targetStage}
		}

		if sourceStage == targetStage {
			return EqualSourceAndTarget
		}

		links := orch.Links()
		links = append(links, name)
		orch = internal.NewOrchestration(orch.Name(), orch.Stages(), links)

		err = orchStorage.Save(orch)
		if err != nil {
			return fmt.Errorf("add link %s: %w", name, err)
		}

		l := internal.NewLink(name, source, target)

		return storage.Save(l)
	}
}
