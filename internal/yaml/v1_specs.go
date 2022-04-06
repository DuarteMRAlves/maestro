package yaml

type v1PipelineSpec struct {
	// Name that should be associated with the pipeline.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
	// Mode specifies the execution mode for the pipeline.
	// (optional)
	Mode string `yaml:"execution_mode,omitempty"`
}

type v1StageSpec struct {
	// Name that should be associated with the stage.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
	// Address where to connect to the grpc server. If not specified, will be
	// inferred from Host and Port as {Host}:{Port}.
	// (optional)
	Address string `yaml:"address" info:"required"`
	// Name of the grpc service that contains the rpc to execute. May be
	// omitted if the target grpc server only has one service.
	// (optional)
	Service string `yaml:"service,omitempty"`
	// Name of the grpc method to execute. May be omitted if the service has
	// only a single method.
	// (optional)
	Method string `yaml:"method,omitempty"`
	// Pipeline specifies the name of the Pipeline where this stage
	// should be inserted.
	// (required)
	Pipeline string `yaml:"pipeline" info:"required"`
}

type v1LinkSpec struct {
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
	SourceField string `yaml:"source_field,omitempty"`
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
	TargetField string `yaml:"target_field,omitempty"`
	// Pipeline specifies the pipeline where this link should be
	// inserted.
	// (required)
	Pipeline string `yaml:"pipeline"`
}

type v1AssetSpec struct {
	// Name that should be associated with the asset. Is required and should be
	// unique.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
	// Image specifies the container image that should be associated with this
	// asset
	// (optional)
	Image string `yaml:"image,omitempty"`
}
