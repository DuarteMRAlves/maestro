package types

// Link represents a connection between two stages of the orchestration, where
// a specific grpc message is transferred from one stage to the next,
type Link struct {
	// Name that should be associated with the link.
	// (required, unique)
	Name string `yaml:"name" info:"required"`
	// SourceStage defines the name of the stage that is the source of the link.
	// The messages returned by the rpc executed in this stage are transferred
	// through this link to the next stage.
	// (required)
	SourceStage StageName `yaml:"source_stage" info:"required"`
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
	TargetStage StageName `yaml:"target_stage" info:"required"`
	// TargetField defines the field of the input of TargetStage that should be
	// filled with the messages transferred with this link. If specified, the
	// field TargetField of message that is the input of TargetStage is set to
	// the messages received through this link. If not specified, the entire
	// message is sent as input to the TargetStage.
	// (optional)
	TargetField string `yaml:"target_field"`
}
