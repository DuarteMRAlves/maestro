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

type LinkResponse struct {
	Err domain.OptionalError
}

var (
	EmptyLinkName    = fmt.Errorf("empty link name")
	EmptySourceField = fmt.Errorf("empty source field")
	EmptyTargetField = fmt.Errorf("empty target field")
)

func CreateLink(
	storage LinkStorage,
	stageLoader StageLoader,
	orchStorage OrchestrationStorage,
) func(LinkRequest) LinkResponse {
	return func(req LinkRequest) LinkResponse {
		name, err := internal.NewLinkName(req.Name)
		if err != nil {
			return LinkResponse{Err: domain.NewPresentError(err)}
		}
		if name.IsEmpty() {
			return LinkResponse{Err: domain.NewPresentError(EmptyLinkName)}
		}

		sourceStage, err := internal.NewStageName(req.SourceStage)
		if err != nil {
			return LinkResponse{Err: domain.NewPresentError(err)}
		}
		sourceFieldOpt := internal.NewEmptyMessageField()
		if req.SourceField.Present() {
			sourceField := internal.NewMessageField(req.SourceField.Unwrap())
			if sourceField.IsEmpty() {
				err := EmptySourceField
				return LinkResponse{Err: domain.NewPresentError(err)}
			}
			sourceFieldOpt = internal.NewPresentMessageField(sourceField)
		}

		targetStage, err := internal.NewStageName(req.TargetStage)
		if err != nil {
			return LinkResponse{Err: domain.NewPresentError(err)}
		}
		targetFieldOpt := internal.NewEmptyMessageField()
		if req.TargetField.Present() {
			targetField := internal.NewMessageField(req.TargetField.Unwrap())
			if targetField.IsEmpty() {
				err := EmptyTargetField
				return LinkResponse{Err: domain.NewPresentError(err)}
			}
			targetFieldOpt = internal.NewPresentMessageField(targetField)
		}

		orchestrationName, err := internal.NewOrchestrationName(req.Orchestration)
		if err != nil {
			return LinkResponse{Err: domain.NewPresentError(err)}
		}

		_, err = storage.Load(name)
		if err == nil {
			err := &internal.AlreadyExists{Type: "link", Ident: name.Unwrap()}
			return LinkResponse{Err: domain.NewPresentError(err)}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return LinkResponse{Err: domain.NewPresentError(err)}
		}

		_, err = orchStorage.Load(orchestrationName)
		if err != nil {
			return LinkResponse{Err: domain.NewPresentError(err)}
		}

		_, err = stageLoader.Load(sourceStage)
		if err != nil {
			return LinkResponse{Err: domain.NewPresentError(err)}
		}

		_, err = stageLoader.Load(targetStage)
		if err != nil {
			return LinkResponse{Err: domain.NewPresentError(err)}
		}

		updateFn := addLinkNameToOrchestration(name)
		err = updateOrchestration(
			orchestrationName,
			orchStorage,
			updateFn,
			orchStorage,
		)
		if err != nil {
			err := errdefs.PrependMsg(
				err,
				"add link %s to orchestration %s",
				name,
				orchestrationName,
			)
			return LinkResponse{Err: domain.NewPresentError(err)}
		}

		sourceEndpoint := internal.NewLinkEndpoint(sourceStage, sourceFieldOpt)
		targetEndpoint := internal.NewLinkEndpoint(targetStage, targetFieldOpt)

		l := internal.NewLink(
			name,
			sourceEndpoint,
			targetEndpoint,
			orchestrationName,
		)

		err = storage.Save(l)
		errOpt := domain.NewEmptyError()
		if err != nil {
			errOpt = domain.NewPresentError(err)
		}
		return LinkResponse{Err: errOpt}
	}
}
