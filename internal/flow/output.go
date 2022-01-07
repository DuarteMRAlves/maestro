package flow

import "github.com/DuarteMRAlves/maestro/internal/link"

// Output represents the several output flows for a stage
type Output struct {
	typ         OutputType
	connections map[string]*link.Link
}

// OutputType defines the type of output the stage.Stage associated with this
// Output is expecting.
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

func NewOutput() *Output {
	return &Output{
		typ:         OutputInfer,
		connections: map[string]*link.Link{},
	}
}
