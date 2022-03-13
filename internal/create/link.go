package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

type LinkSaver interface {
	Save(internal.Link) LinkResult
}

type LinkLoader interface {
	Load(internal.LinkName) LinkResult
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
	EmptySourceField = fmt.Errorf("empty source field")
	EmptyTargetField = fmt.Errorf("empty target field")
)

func CreateLink(
	storage LinkStorage,
	stageLoader StageLoader,
	orchStorage OrchestrationStorage,
) func(LinkRequest) LinkResponse {
	return func(req LinkRequest) LinkResponse {
		res := requestToLink(req)
		res = BindLink(verifyDupLink(storage))(res)
		res = BindLink(verifyExistsOrchestrationLink(orchStorage))(res)
		res = BindLink(verifyExistsSource(stageLoader))(res)
		res = BindLink(verifyExistsTarget(stageLoader))(res)
		res = BindLink(addLink(orchStorage, orchStorage))(res)
		res = BindLink(storage.Save)(res)
		return linkToResponse(res)
	}
}

func requestToLink(req LinkRequest) LinkResult {
	name, err := internal.NewLinkName(req.Name)
	if err != nil {
		return ErrLink(err)
	}
	sourceStage, err := internal.NewStageName(req.SourceStage)
	if err != nil {
		return ErrLink(err)
	}

	sourceFieldOpt := internal.NewEmptyMessageField()
	if req.SourceField.Present() {
		sourceField := internal.NewMessageField(req.SourceField.Unwrap())
		if sourceField.IsEmpty() {
			return ErrLink(EmptySourceField)
		}
		sourceFieldOpt = internal.NewPresentMessageField(sourceField)
	}

	targetStage, err := internal.NewStageName(req.TargetStage)

	targetFieldOpt := internal.NewEmptyMessageField()
	if req.TargetField.Present() {
		targetField := internal.NewMessageField(req.TargetField.Unwrap())
		if targetField.IsEmpty() {
			return ErrLink(EmptyTargetField)
		}
		targetFieldOpt = internal.NewPresentMessageField(targetField)
	}

	orchestrationName, err := internal.NewOrchestrationName(req.Orchestration)
	if err != nil {
		return ErrLink(err)
	}

	sourceEndpoint := internal.NewLinkEndpoint(sourceStage, sourceFieldOpt)
	targetEndpoint := internal.NewLinkEndpoint(targetStage, targetFieldOpt)

	l := internal.NewLink(
		name,
		sourceEndpoint,
		targetEndpoint,
		orchestrationName,
	)

	return SomeLink(l)
}

func verifyDupLink(loader LinkLoader) func(internal.Link) LinkResult {
	return func(l internal.Link) LinkResult {
		res := loader.Load(l.Name())
		if res.IsError() {
			var notFound *internal.NotFound
			err := res.Error()
			if errors.As(err, &notFound) {
				return SomeLink(l)
			}
			return ErrLink(err)
		}
		err := &internal.AlreadyExists{Type: "link", Ident: l.Name().Unwrap()}
		return ErrLink(err)
	}
}

func verifyExistsOrchestrationLink(orchLoader OrchestrationLoader) func(internal.Link) LinkResult {
	return func(l internal.Link) LinkResult {
		_, err := orchLoader.Load(l.Orchestration())
		if err != nil {
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func verifyExistsSource(stageLoader StageLoader) func(internal.Link) LinkResult {
	return func(l internal.Link) LinkResult {
		_, err := stageLoader.Load(l.Source().Stage())
		if err != nil {
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func verifyExistsTarget(stageLoader StageLoader) func(internal.Link) LinkResult {
	return func(l internal.Link) LinkResult {
		_, err := stageLoader.Load(l.Target().Stage())
		if err != nil {
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func addLink(
	loader OrchestrationLoader,
	saver OrchestrationSaver,
) func(internal.Link) LinkResult {
	return func(l internal.Link) LinkResult {
		name := l.Orchestration()
		updateFn := addLinkNameToOrchestration(l.Name())
		err := updateOrchestration(name, loader, updateFn, saver)
		if err != nil {
			err := errdefs.PrependMsg(
				err,
				"add link %s to orchestration %s",
				l.Name(),
				l.Orchestration(),
			)
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func linkToResponse(res LinkResult) LinkResponse {
	var errOpt domain.OptionalError
	if res.IsError() {
		errOpt = domain.NewPresentError(res.Error())
	} else {
		errOpt = domain.NewEmptyError()
	}
	return LinkResponse{Err: errOpt}
}
