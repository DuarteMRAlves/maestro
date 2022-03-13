package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
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

type LinkRequest struct {
	Name string

	SourceStage string
	SourceField domain.OptionalString
	TargetStage string
	TargetField domain.OptionalString

	Orchestration string
}

var (
	EmptyLinkName        = fmt.Errorf("empty link name")
	EmptySourceField     = fmt.Errorf("empty source field")
	EmptyTargetField     = fmt.Errorf("empty target field")
	EqualSourceAndTarget = fmt.Errorf("equal source and target stages")
)

type StageNotInOrchestration struct {
	Orch  internal.OrchestrationName
	Stage internal.StageName
}

func (err *StageNotInOrchestration) Error() string {
	format := "stage '%s' not found in orchestration '%s'."
	return fmt.Sprintf(format, err.Stage, err.Orch)
}

func CreateLink(
	storage LinkStorage,
	stageLoader StageLoader,
	orchStorage OrchestrationStorage,
) func(LinkRequest) error {
	return func(req LinkRequest) error {
		name, err := internal.NewLinkName(req.Name)
		if err != nil {
			return err
		}
		if name.IsEmpty() {
			return EmptyLinkName
		}

		sourceStage, err := internal.NewStageName(req.SourceStage)
		if err != nil {
			return err
		}
		sourceFieldOpt := internal.NewEmptyMessageField()
		if req.SourceField.Present() {
			sourceField := internal.NewMessageField(req.SourceField.Unwrap())
			if sourceField.IsEmpty() {
				return EmptySourceField
			}
			sourceFieldOpt = internal.NewPresentMessageField(sourceField)
		}

		targetStage, err := internal.NewStageName(req.TargetStage)
		if err != nil {
			return err
		}
		targetFieldOpt := internal.NewEmptyMessageField()
		if req.TargetField.Present() {
			targetField := internal.NewMessageField(req.TargetField.Unwrap())
			if targetField.IsEmpty() {
				return EmptyTargetField
			}
			targetFieldOpt = internal.NewPresentMessageField(targetField)
		}

		orchName, err := internal.NewOrchestrationName(req.Orchestration)
		if err != nil {
			return err
		}

		_, err = storage.Load(name)
		if err == nil {
			return &internal.AlreadyExists{Type: "link", Ident: name.Unwrap()}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return err
		}

		orch, err := orchStorage.Load(orchName)
		if err != nil {
			return err
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

		updateFn := addLinkNameToOrchestration(name)
		err = updateOrchestration(orchName, orchStorage, updateFn, orchStorage)
		if err != nil {
			format := "add link %s to orchestration %s"
			return errdefs.PrependMsg(err, format, name, orchName)
		}

		sourceEndpoint := internal.NewLinkEndpoint(sourceStage, sourceFieldOpt)
		targetEndpoint := internal.NewLinkEndpoint(targetStage, targetFieldOpt)

		l := internal.NewLink(name, sourceEndpoint, targetEndpoint, orchName)

		return storage.Save(l)
	}
}
