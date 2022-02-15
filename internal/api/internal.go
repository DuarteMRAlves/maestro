package api

// InternalAPI is an interface that collects all the available commands
// for the maestro server. All calls on external APIs should be redirected
// through this API that collects all functionality.
type InternalAPI interface {
	CreateAsset(*CreateAssetRequest) error
	GetAsset(*GetAssetRequest) ([]*Asset, error)

	CreateStage(*CreateStageRequest) error
	GetStage(*GetStageRequest) ([]*Stage, error)

	CreateLink(*CreateLinkRequest) error
	GetLink(*GetLinkRequest) ([]*Link, error)

	CreateOrchestration(*CreateOrchestrationRequest) error
	GetOrchestration(*GetOrchestrationRequest) ([]*Orchestration, error)

	StartExecution(*StartExecutionRequest) error
	AttachExecution(*AttachExecutionRequest) ([]*Event, <-chan *Event, error)
}
