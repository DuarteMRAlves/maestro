package flow

import "github.com/DuarteMRAlves/maestro/internal/link"

// Output receives the output flow.State for a given stage and sends it to the
// next stages.
type Output interface {
	Out() chan<- *State
}

// OutputCfg represents the several output flows for a stage
type OutputCfg struct {
	typ         OutputType
	connections map[string]*link.Link
}

// OutputType defines the type of output the stage.Stage associated with this
// OutputCfg is expecting.
type OutputType string

const (
	// OutputInfer means the output type is not specified and should be inferred
	// from the received connections.
	OutputInfer OutputType = "Infer"
	// OutputSingle means the stage only sends its entire output to another
	// stage.
	OutputSingle OutputType = "Single"
	// OutputSink means the stage outputs the empty message and its output
	// should be discarded as it does not connect to any other stage.
	OutputSink OutputType = "Sink"
	// OutputSplit means the stage output should be split into multiple messages
	// sent to different stages.
	OutputSplit OutputType = "Split"
	// OutputDuplicate means the stage output should be duplicated and sent as
	// input to multiple stages.
	OutputDuplicate OutputType = "Duplicate"
)

func NewOutputCfg() *OutputCfg {
	return &OutputCfg{
		typ:         OutputInfer,
		connections: map[string]*link.Link{},
	}
}

func (o *OutputCfg) ToFlow() Output {
	return nil
}
