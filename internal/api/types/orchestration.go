package types

// Orchestration represents a graph composed by stages and links that should be
// executed.
type Orchestration struct {
	// Name that should be associated with the orchestration.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
	// Links specifies the name of the links that compose this orchestration.
	// (required, non-empty)
	Links []string `yaml:"links" info:"required"`
}
