package storage

import "github.com/DuarteMRAlves/maestro/internal/api"

const (
	defaultOrchestrationName = api.OrchestrationName("default")
	defaultStageHost         = "localhost"
	defaultStagePort         = 8061
)

func defaultOrchestration() *api.Orchestration {
	return &api.Orchestration{
		Name:   defaultOrchestrationName,
		Phase:  api.OrchestrationPending,
		Stages: []api.StageName{},
		Links:  []api.LinkName{},
	}
}
