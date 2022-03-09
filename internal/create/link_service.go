package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

func CreateLink(
	storage LinkStorage,
	existsStage ExistsStage,
	orchStorage OrchestrationStorage,
) func(LinkRequest) LinkResponse {
	return func(req LinkRequest) LinkResponse {
		res := requestToLink(req)
		res = BindLink(verifyDupLink(storage))(res)
		res = BindLink(verifyExistsOrchestrationLink(orchStorage))(res)
		res = BindLink(verifyExistsSource(existsStage))(res)
		res = BindLink(verifyExistsTarget(existsStage))(res)
		res = BindLink(addLink(orchStorage, orchStorage))(res)
		res = BindLink(storage.Save)(res)
		return linkToResponse(res)
	}
}

func requestToLink(req LinkRequest) LinkResult {
	name, err := domain.NewLinkName(req.Name)
	if err != nil {
		return ErrLink(err)
	}
	sourceStage, err := domain.NewStageName(req.SourceStage)
	if err != nil {
		return ErrLink(err)
	}

	sourceFieldOpt := domain.NewEmptyMessageField()
	if req.SourceField.Present() {
		sourceField, err := domain.NewMessageField(req.SourceField.Unwrap())
		if err != nil {
			return ErrLink(err)
		}
		sourceFieldOpt = domain.NewPresentMessageField(sourceField)
	}

	targetStage, err := domain.NewStageName(req.TargetStage)

	targetFieldOpt := domain.NewEmptyMessageField()
	if req.TargetField.Present() {
		targetField, err := domain.NewMessageField(req.TargetField.Unwrap())
		if err != nil {
			return ErrLink(err)
		}
		targetFieldOpt = domain.NewPresentMessageField(targetField)
	}

	orchestrationName, err := domain.NewOrchestrationName(req.Orchestration)
	if err != nil {
		return ErrLink(err)
	}

	sourceEndpoint := NewLinkEndpoint(sourceStage, sourceFieldOpt)
	targetEndpoint := NewLinkEndpoint(targetStage, targetFieldOpt)

	l := NewLink(name, sourceEndpoint, targetEndpoint, orchestrationName)

	return SomeLink(l)
}

func verifyDupLink(verifier LinkExistsVerifier) func(Link) LinkResult {
	return func(l Link) LinkResult {
		if verifier.Verify(l.Name()) {
			err := errdefs.AlreadyExistsWithMsg(
				"link '%v' already exists",
				l.Name().Unwrap(),
			)
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func verifyExistsOrchestrationLink(verifier OrchestrationExistsVerifier) func(Link) LinkResult {
	return func(l Link) LinkResult {
		if !verifier.Verify(l.Orchestration()) {
			err := errdefs.NotFoundWithMsg(
				"orchestration '%v' not found",
				l.Orchestration().Unwrap(),
			)
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func verifyExistsSource(existsFn ExistsStage) func(Link) LinkResult {
	return func(l Link) LinkResult {
		if !existsFn(l.Source().Stage()) {
			err := errdefs.NotFoundWithMsg(
				"source '%v' not found",
				l.Source().Stage().Unwrap(),
			)
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func verifyExistsTarget(existsFn ExistsStage) func(Link) LinkResult {
	return func(l Link) LinkResult {
		if !existsFn(l.Target().Stage()) {
			err := errdefs.NotFoundWithMsg(
				"target '%v' not found",
				l.Target().Stage().Unwrap(),
			)
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func addLink(
	loader OrchestrationLoader,
	saver OrchestrationSaver,
) func(Link) LinkResult {
	return func(l Link) LinkResult {
		name := l.Orchestration()
		updateFn := ReturnOrchestration(addLinkNameToOrchestration(l.Name()))
		res := updateOrchestration(name, loader, updateFn, saver)
		if res.IsError() {
			err := errdefs.PrependMsg(
				res.Error(),
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
