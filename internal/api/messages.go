package api

// CreateOrchestrationRequest represents a message to create an orchestration.
type CreateOrchestrationRequest struct {
	// Name that should be associated with the orchestration.
	// (required, unique)
	Name OrchestrationName `yaml:"name" info:"required"`
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
	// Orchestration specifies the name of the Orchestration where this stage
	// should be inserted. If not specified, the stage will be inserted into the
	// default Orchestration.
	// (optional)
	Orchestration OrchestrationName `yaml:"orchestration"`
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
	// Orchestration should be specified to retrieve stages with the given
	// orchestration.
	// (optional)
	Orchestration OrchestrationName
	// Asset should be specified to retrieve stages with the given asset.
	// (optional)
	Asset AssetName
}

// CreateLinkRequest represents a message to create a link.
type CreateLinkRequest struct {
	// Name that should be associated with the link.
	// (required, unique)
	Name LinkName `yaml:"name" info:"required"`
	// SourceStage defines the name of the stage that is the source of the link.
	// The messages returned by the rpc executed in this stage are transferred
	// through this link to the next stage.
	// (required)
	SourceStage StageName `yaml:"source_stage" info:"required"`
	// SourceField defines the field of the source message that should be sent
	// through the link. If specified, the message transferred through this link
	// is the field with the given name from the message returned by SourceStage.
	// If not specified, the entire message from SourceStage is used.
	// (optional)
	SourceField string `yaml:"source_field"`
	// TargetStage defines the name of the stage that is the target of the link.
	// The messages that are transferred through this link are used as input for
	// the rpc method in this stage.
	// (required)
	TargetStage StageName `yaml:"target_stage" info:"required"`
	// TargetField defines the field of the input of TargetStage that should be
	// filled with the messages transferred with this link. If specified, the
	// field TargetField of message that is the input of TargetStage is set to
	// the messages received through this link. If not specified, the entire
	// message is sent as input to the TargetStage.
	// (optional)
	TargetField string `yaml:"target_field"`
	// Orchestration specifies the orchestration where this link should be
	// inserted. If not specified, the link will be inserted into the "default"
	// orchestration.
	// (optional)
	Orchestration OrchestrationName `yaml:"orchestration"`
}

// GetLinkRequest is a message to retrieve links with specific characteristics.
type GetLinkRequest struct {
	// Name should be set to retrieve only links with the given name.
	// (optional)
	Name LinkName
	// SourceStage should be set to retrieve only assets with the given source
	// stage.
	// (optional)
	SourceStage StageName
	// SourceStage should be set to retrieve only assets with the given source
	// stage.
	// (optional)
	SourceField string
	// TargetStage should be set to retrieve only assets with the given target
	// stage.
	// (optional)
	TargetStage StageName
	// TargetField should be set to retrieve only assets with the given target
	// field.
	// (optional)
	TargetField string
	// Orchestration should be specified to retrieve links with the given
	// orchestration.
	// (optional)
	Orchestration OrchestrationName
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
