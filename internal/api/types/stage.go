package types

// Stage represents a node of the pipeline where a specific rpc method is
// executed.
type Stage struct {
	// Name that should be associated with the stage. Is required and should be
	// unique.
	Name string
	// Name of the grpc service that contains the method to execute. May be
	// omitted if the target grpc server only has one service.
	// (optional)
	Service string
	// Name of the grpc method to execute. May be omitted if the service has
	// only a single method.
	// (optional)
	Method string
	// Address where to connect to the grpc server.
	Address string
	// Asset that should be used to run the stage
	// (optional)
	Asset string
}
