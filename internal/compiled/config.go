package compiled

// Pipeline specifies the schema of a pipeline to be orchestrated.
type PipelineConfig struct {
	Name   string
	Mode   ExecutionMode
	Stages []*StageConfig
	Links  []*LinkConfig
}

// Stage specifies a given step of the Pipeline.
type StageConfig struct {
	Name     string
	MethodID MethodID
}

// Link defines a connection between two Stage objects in a Pipeline.
type LinkConfig struct {
	Name        string
	SourceStage string
	SourceField string
	TargetStage string
	TargetField string
}
