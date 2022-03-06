package create

import (
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

func CreateLink(
	existsLink ExistsLink,
	saveLink SaveLink,
	existsStage ExistsStage,
	existsOrchestration ExistsOrchestration,
	loadOrchestration LoadOrchestration,
	saveOrchestration SaveOrchestration,
) func(LinkRequest) LinkResponse {
	return func(req LinkRequest) LinkResponse {
		res := requestToLink(req)
		res = BindLink(verifyDupLink(existsLink))(res)
		res = BindLink(verifyExistsOrchestrationLink(existsOrchestration))(res)
		res = BindLink(verifyExistsSource(existsStage))(res)
		res = BindLink(verifyExistsTarget(existsStage))(res)
		res = BindLink(addLink(loadOrchestration, saveOrchestration))(res)
		res = BindLink(saveLink)(res)
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

func verifyDupLink(existsFn ExistsLink) func(Link) LinkResult {
	return func(l Link) LinkResult {
		if existsFn(l.Name()) {
			err := errdefs.AlreadyExistsWithMsg(
				"link '%v' already exists",
				l.Name().Unwrap(),
			)
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func verifyExistsOrchestrationLink(existsFn ExistsOrchestration) func(Link) LinkResult {
	return func(l Link) LinkResult {
		if !existsFn(l.Orchestration()) {
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
	loadFn LoadOrchestration,
	saveFn SaveOrchestration,
) func(Link) LinkResult {
	return func(l Link) LinkResult {
		name := l.Orchestration()
		updateFn := addLinkToOrchestration(l)
		res := updateOrchestration(name, loadFn, updateFn, saveFn)
		if res.IsError() {
			err := errdefs.PrependMsg(
				res.Error(),
				"add link %s to orchestration: %s",
				l.Name(),
				l.Orchestration(),
			)
			return ErrLink(err)
		}
		return SomeLink(l)
	}
}

func addLinkToOrchestration(link Link) func(Orchestration) OrchestrationResult {
	return func(o Orchestration) OrchestrationResult {
		old := o.Links()
		links := make([]domain.LinkName, 0, len(old)+1)
		for _, s := range old {
			links = append(links, s)
		}
		links = append(links, link.Name())
		updated := NewOrchestration(o.Name(), o.Stages(), links)
		return SomeOrchestration(updated)
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
