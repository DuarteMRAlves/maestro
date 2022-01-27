package api

// Orchestration represents a graph composed by stages and links that should be
// executed.
type Orchestration struct {
	// Name that should be associated with the orchestration.
	// (required, unique)
	Name OrchestrationName `yaml:"name" info:"required"`
	// Phase specifies the current phase of the Orchestration. This field should
	// not be defined in a yaml file as it is a value defined by the current
	// state of the system.
	// (optional)
	Phase OrchestrationPhase `yaml:"-"`
	// Stages specifies the name of the stages that compose this orchestration.
	// (required, non-empty)
	Stages []StageName `yaml:"stages" info:"required"`
	// Links specifies the name of the links that compose this orchestration.
	// (required, non-empty)
	Links []LinkName `yaml:"links" info:"required"`
}

// OrchestrationName is a name that uniquely identifies an Orchestration.
type OrchestrationName string

// OrchestrationPhase is an enum that describes the current status of the
// orchestration.
type OrchestrationPhase string

const (
	// OrchestrationPending means the orchestration was accepted by the system
	// and is waiting to be executed.
	OrchestrationPending OrchestrationPhase = "Pending"
	// OrchestrationRunning means the orchestration is currently running
	OrchestrationRunning OrchestrationPhase = "Running"
	// OrchestrationSucceeded means all the stages in the orchestration
	// voluntarily terminated without any error.
	OrchestrationSucceeded OrchestrationPhase = "Succeeded"
	// OrchestrationFailed means the orchestration terminated with at least one
	// failure in some stage.
	OrchestrationFailed OrchestrationPhase = "Failed"
)

// Stage represents a node of the pipeline where a specific grpc method is
// executed.
type Stage struct {
	// Name that should be associated with the stage.
	// (unique)
	Name StageName
	// Phase specifies the current phase of the Stage. This field should not be
	// specified in a yaml file as it is a value defined by the current state
	// of the system.
	Phase StagePhase
	// Name of the grpc service that contains the rpc to execute. May be
	// omitted if the target grpc server only has one service.
	Service string
	// Name of the grpc method to execute. May be omitted if the service has
	// only a single method.
	Rpc string
	// Address where to connect to the grpc server.
	Address string
	// Orchestration that is associated with this stage
	Orchestration OrchestrationName
	// Asset that should be used to run the stage.
	Asset AssetName
}

// StageName is a name that uniquely identifies a Stage
type StageName string

// StagePhase is an enum that describes the current status of a stage
type StagePhase string

const (
	// StagePending means the stage was accepted by the system and is awaiting
	// to be executed. In this phase, the stage can be linked to other stages
	// in the orchestration.
	StagePending StagePhase = "Pending"

	// StageRunning means the stage is currently running.
	StageRunning StagePhase = "Running"

	// StageSucceeded means the stage voluntarily terminated with no errors.
	StageSucceeded StagePhase = "Succeeded"

	// StageFailed means the stage exited terminated with a failure.
	StageFailed StagePhase = "Failed"
)

// Link represents a connection between two stages of the orchestration, where
// a specific grpc message is transferred from one stage to the next,
type Link struct {
	// Name that should be associated with the link.
	// (unique)
	Name LinkName
	// SourceStage defines the name of the stage that is the source of the link.
	// The messages returned by the rpc executed in this stage are transferred
	// through this link to the next stage.
	SourceStage StageName
	// SourceField defines the field of the source message that should be sent
	// through the link. If specified, the message transferred through this link
	// is the field with the given name from the message returned by SourceStage.
	// If not specified, the entire message from SourceStage is used.
	// (optional)
	SourceField string
	// TargetStage defines the name of the stage that is the target of the link.
	// The messages that are transferred through this link are used as input for
	// the rpc method in this stage.
	TargetStage StageName
	// TargetField defines the field of the input of TargetStage that should be
	// filled with the messages transferred with this link. If specified, the
	// field TargetField of message that is the input of TargetStage is set to
	// the messages received through this link. If not specified, the entire
	// message is sent as input to the TargetStage.
	// (optional)
	TargetField string
}

// LinkName is a name that uniquely identifies a Link
type LinkName string

// Asset represents an image with a grpc server that can be deployed.
type Asset struct {
	// Name that should be associated with the asset. Is required and should be
	// unique.
	// (required, unique)
	Name AssetName `yaml:"name" info:"required"`
	// Image specifies the container image that should be associated with this
	// asset
	// (optional)
	Image string `yaml:"image"`
}

// AssetName is a name that uniquely identifies an Asset
type AssetName string
