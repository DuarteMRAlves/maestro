package api

// InternalAPI is an interface that collects all the available commands
// for the maestro server. All calls on external APIs should be redirected
// through this API that collects all functionality.
type InternalAPI interface {
	CreateAsset(config *CreateAssetRequest) error
	GetAsset(query *GetAssetRequest) ([]*Asset, error)

	CreateStage(config *Stage) error
	GetStage(query *Stage) ([]*Stage, error)

	CreateLink(config *Link) error
	GetLink(query *Link) ([]*Link, error)

	CreateOrchestration(*CreateOrchestrationRequest) error
	GetOrchestration(*GetOrchestrationRequest) (
		[]*Orchestration,
		error,
	)
}
