package api

// Pipeline specifies the schema of a pipeline to be orchestrated.
type Pipeline struct {
	Name   string
	Stages []*Stage
	Links  []*Link
}

// Stage specifies a given step of the Pipeline.
type Stage struct {
	Name    string
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
