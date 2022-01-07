package types

// Orchestration represents a graph composed by stages and links that should be
// executed.
type Orchestration struct {
	// Name that should be associated with the orchestration.
	// (required, unique)
	Name OrchestrationName `yaml:"name" info:"required"`
	// Links specifies the name of the links that compose this orchestration.
	// (required, non-empty)
	Links []string `yaml:"links" info:"required"`
}

// OrchestrationName is a name that uniquely identifies an Orchestration.
type OrchestrationName string
