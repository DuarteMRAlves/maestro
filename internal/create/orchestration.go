package create

import (
	"errors"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

type OrchestrationSaver interface {
	Save(internal.Orchestration) error
}

type OrchestrationLoader interface {
	Load(internal.OrchestrationName) (internal.Orchestration, error)
}

type OrchestrationStorage interface {
	OrchestrationSaver
	OrchestrationLoader
}

type OrchestrationRequest struct {
	Name string
}

type OrchestrationResponse struct {
	Err domain.OptionalError
}

var EmptyOrchestrationName = fmt.Errorf("empty orchestration name")

func Create(storage OrchestrationStorage) func(OrchestrationRequest) OrchestrationResponse {
	return func(req OrchestrationRequest) OrchestrationResponse {
		name, err := internal.NewOrchestrationName(req.Name)
		if err != nil {
			return OrchestrationResponse{Err: domain.NewPresentError(err)}
		}
		if name.IsEmpty() {
			err := EmptyOrchestrationName
			return OrchestrationResponse{Err: domain.NewPresentError(err)}
		}

		_, err = storage.Load(name)
		if err == nil {
			err := errdefs.AlreadyExistsWithMsg(
				"orchestration '%v' already exists",
				name.Unwrap(),
			)
			return OrchestrationResponse{Err: domain.NewPresentError(err)}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return OrchestrationResponse{Err: domain.NewPresentError(err)}
		}

		o := internal.NewOrchestration(
			name,
			[]internal.StageName{},
			[]internal.LinkName{},
		)
		err = storage.Save(o)
		errOpt := domain.NewEmptyError()
		if err != nil {
			errOpt = domain.NewPresentError(err)
		}
		return OrchestrationResponse{Err: errOpt}
	}
}

func updateOrchestration(
	name internal.OrchestrationName,
	loader OrchestrationLoader,
	updateFn func(internal.Orchestration) internal.Orchestration,
	saver OrchestrationSaver,
) error {
	orch, err := loader.Load(name)
	if err != nil {
		return err
	}
	orch = updateFn(orch)
	return saver.Save(orch)
}

func addStageNameToOrchestration(
	s internal.StageName,
) func(internal.Orchestration) internal.Orchestration {
	return func(o internal.Orchestration) internal.Orchestration {
		old := o.Stages()
		stages := make([]internal.StageName, 0, len(old)+1)
		for _, name := range old {
			stages = append(stages, name)
		}
		stages = append(stages, s)
		return internal.NewOrchestration(o.Name(), stages, o.Links())
	}
}

func addLinkNameToOrchestration(l internal.LinkName) func(internal.Orchestration) internal.Orchestration {
	return func(o internal.Orchestration) internal.Orchestration {
		old := o.Links()
		links := make([]internal.LinkName, 0, len(old)+1)
		for _, name := range old {
			links = append(links, name)
		}
		links = append(links, l)
		return internal.NewOrchestration(o.Name(), o.Stages(), links)
	}
}
