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
	// Address where to connect to the grpc server. If not specified, will be
	// inferred from Host and Port as {Host}:{Port}.
	// (optional, conflicts: Host, Port)
	Address string
	// Host where to connect to the grpc server. If not specified will be set
	// to localhost. Should not be specified if Address is defined.
	// (optional, conflicts: Address)
	Host string
	// Port where to connect to the grpc server. If not specified will be set
	// to 8061. Should not be specified if Address is defined.
	// (optional, conflicts: Address)
	Port int32
	// Asset that should be used to run the stage
	// (optional)
	Asset string
}
