package types

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
	// Links specifies the name of the links that compose this orchestration.
	// (required, non-empty)
	Links []string `yaml:"links" info:"required"`
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
