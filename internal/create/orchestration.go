package create

import (
	"errors"
	"github.com/DuarteMRAlves/maestro/internal"
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

var EmptyOrchestrationName = errors.New("empty orchestration name")

func Orchestration(storage OrchestrationStorage) func(internal.OrchestrationName) error {
	return func(name internal.OrchestrationName) error {
		if name.IsEmpty() {
			return EmptyOrchestrationName
		}

		_, err := storage.Load(name)
		if err == nil {
			return &internal.AlreadyExists{
				Type:  "orchestration",
				Ident: name.Unwrap(),
			}
		}
		var notFound *internal.NotFound
		if !errors.As(err, &notFound) {
			return err
		}

		o := internal.NewOrchestration(name, nil, nil)
		return storage.Save(o)
	}
}
