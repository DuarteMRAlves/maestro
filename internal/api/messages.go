package api

// CreateOrchestrationRequest represents a message to create an orchestration.
type CreateOrchestrationRequest struct {
	// Name that should be associated with the orchestration.
	// (required, unique)
	Name OrchestrationName `yaml:"name" info:"required"`
	// Links specifies the name of the links that compose this orchestration.
	// (required, non-empty)
	Links []LinkName `yaml:"links" info:"required"`
}

// GetOrchestrationRequest represents a message to retrieve orchestrations with
// specific characteristics.
type GetOrchestrationRequest struct {
	// Name should be set to retrieve orchestrations with the given name.
	Name OrchestrationName
	// Phase should be set to retrieve orchestrations in a given phase.
	Phase OrchestrationPhase
}

// CreateStageRequest represents a message to create a stage.
type CreateStageRequest struct {
	// Name that should be associated with the stage.
	// (required, unique)
	Name StageName `yaml:"name" info:"required"`
	// Name of the grpc service that contains the rpc to execute. May be
	// omitted if the target grpc server only has one service.
	// (optional)
	Service string `yaml:"service"`
	// Name of the grpc method to execute. May be omitted if the service has
	// only a single method.
	// (optional)
	Rpc string `yaml:"rpc"`
	// Address where to connect to the grpc server. If not specified, will be
	// inferred from Host and Port as {Host}:{Port}.
	// (optional, conflicts: Host, Port)
	Address string `yaml:"address"`
	// Host where to connect to the grpc server. If not specified will be set
	// to localhost. Should not be specified if Address is defined.
	// (optional, conflicts: Address)
	Host string `yaml:"host"`
	// Port where to connect to the grpc server. If not specified will be set
	// to 8061. Should not be specified if Address is defined.
	// (optional, conflicts: Address)
	Port int32 `yaml:"port"`
	// Asset that should be used to run the stage
	// (optional)
	Asset AssetName `yaml:"asset"`
}

// GetStageRequest is a message to retrieve stages with specific characteristics.
type GetStageRequest struct {
	// Name should be specified to retrieve stages with the given name.
	// (optional)
	Name StageName
	// Phase should be specified to retrieve stages with the given phase.
	// (optional)
	Phase StagePhase
	// Service should be specified to retrieve stages with the given service.
	// (optional)
	Service string
	// Rpc should be specified to retrieve stages with the given rpc.
	// (optional)
	Rpc string
	// Address should be specified to retrieve stages with the given address.
	// (optional)
	Address string
	// Asset should be specified to retrieve stages with the given asset.
	// (optional)
	Asset AssetName
}

// CreateAssetRequest represents a message to create an asset.
type CreateAssetRequest struct {
	// Name that should be associated with the asset. Is required and should be
	// unique.
	// (required, unique)
	Name AssetName `yaml:"name" info:"required"`
	// Image specifies the container image that should be associated with this
	// asset
	// (optional)
	Image string `yaml:"image"`
}

// GetAssetRequest represents a message to get assets with specific
// characteristics.
type GetAssetRequest struct {
	// Name should be set to retrieve only assets with the given name.
	Name AssetName
	// Image should be set to retrieve only assets with the given image.
	Image string
}
