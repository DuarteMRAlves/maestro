package flow

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
)

// Output receives the output flow.State for a given stage and sends it to the
// next stages.
type Output interface {
	Out() chan<- *State
}

// OutputCfg represents the several output flows for a stage
type OutputCfg struct {
	typ   OutputType
	links []*link.Link
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

func newOutputCfg() *OutputCfg {
	return &OutputCfg{
		typ:   OutputInfer,
		links: []*link.Link{},
	}
}

func (o *OutputCfg) register(link *link.Link) error {
	for _, l := range o.links {
		if link.Name() == l.Name() {
			return errdefs.InvalidArgumentWithMsg(
				"link with an equal name already registered: %s",
				l.Name())
		}
	}

	o.links = append(o.links, link)
	return nil
}

func (o *OutputCfg) unregisterIfExists(search *link.Link) {
	idx := -1
	for i, l := range o.links {
		if l.Name() == search.Name() {
			idx = i
			break
		}
	}
	if idx != -1 {
		o.links[idx] = o.links[len(o.links)-1]
		o.links = o.links[:len(o.links)-1]
	}
}

func (o *OutputCfg) ToFlow() Output {
	return nil
}
