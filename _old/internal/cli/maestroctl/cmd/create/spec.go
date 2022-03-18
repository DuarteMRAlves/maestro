package create

type OrchestrationSpec struct {
	// Name that should be associated with the orchestration.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
}

type StageSpec struct {
	// Name that should be associated with the stage.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
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
	// Host where to connect to the grpc server. Should not be specified if
	// Address is defined.
	// (optional, conflicts: Address)
	Host string `yaml:"host"`
	// Port where to connect to the grpc server. Should not be specified if
	// Address is defined.
	// (optional, conflicts: Address)
	Port int32 `yaml:"port"`
	// Orchestration specifies the name of the Orchestration where this stage
	// should be inserted.
	// (required)
	Orchestration string `yaml:"orchestration" info:"required"`
	// Asset that should be used to run the stage
	// (optional)
	Asset string `yaml:"asset"`
}

type LinkSpec struct {
	// Name that should be associated with the link.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
	// SourceStage defines the name of the stage that is the source of the link.
	// The messages returned by the rpc executed in this stage are transferred
	// through this link to the next stage.
	// (required)
	SourceStage string `yaml:"source_stage" info:"required"`
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
	TargetStage string `yaml:"target_stage" info:"required"`
	// TargetField defines the field of the input of TargetStage that should be
	// filled with the messages transferred with this link. If specified, the
	// field TargetField of message that is the input of TargetStage is set to
	// the messages received through this link. If not specified, the entire
	// message is sent as input to the TargetStage.
	// (optional)
	TargetField string `yaml:"target_field"`
	// Orchestration specifies the orchestration where this link should be
	// inserted.
	// (required)
	Orchestration string `yaml:"orchestration"`
}

type AssetSpec struct {
	// Name that should be associated with the asset. Is required and should be
	// unique.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
	// Image specifies the container image that should be associated with this
	// asset
	// (optional)
	Image string `yaml:"image"`
}
