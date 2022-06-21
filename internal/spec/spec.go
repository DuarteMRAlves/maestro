package spec

// Pipeline specifies the schema of a pipeline to be orchestrated.
type Pipeline struct {
	Name   string
	Mode   ExecutionMode
	Stages []*Stage
	Links  []*Link
}

type ExecutionMode int

const (
	OfflineExecution ExecutionMode = iota
	OnlineExecution
)

func (e ExecutionMode) String() string {
	switch e {
	case OfflineExecution:
		return "OfflineExecution"
	case OnlineExecution:
		return "OnlineExecution"
	default:
		return "UnknownExecutionMode"
	}
}

// Stage specifies a given step of the Pipeline.
type Stage struct {
	Name          string
	MethodContext MethodContext
}

// MethodContext specifies the method to be executed in a given
// Stage.
type MethodContext struct {
	Address string
	Service string
	Method  string
}

// Link defines a connection between two Stage objects in a Pipeline.
type Link struct {
	Name        string
	SourceStage string
	SourceField string
	TargetStage string
	TargetField string
	// Number of empty messages to fill the link with.
	NumEmptyMessages uint
}
