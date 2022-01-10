package flow

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
)

// Output receives the output flow.State for a given stage and sends it to the
// next stages.
type Output interface {
	Yield(s *State)
}

// OutputCfg represents the several output flows for a stage
type OutputCfg struct {
	typ   OutputType
	flows []*Flow
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
		flows: []*Flow{},
	}
}

func (o *OutputCfg) register(f *Flow) error {
	l := f.link
	for _, prev := range o.flows {
		if l.Name() == prev.link.Name() {
			return errdefs.InvalidArgumentWithMsg(
				"link with an equal name already registered: %s",
				prev.link.Name())
		}
	}

	o.flows = append(o.flows, f)
	return nil
}

func (o *OutputCfg) unregisterIfExists(search *Flow) {
	idx := -1
	for i, f := range o.flows {
		if f.link.Name() == search.link.Name() {
			idx = i
			break
		}
	}
	if idx != -1 {
		o.flows[idx] = o.flows[len(o.flows)-1]
		o.flows = o.flows[:len(o.flows)-1]
	}
}

func (o *OutputCfg) ToFlow() Output {
	switch len(o.flows) {
	case 1:
		return &SingleOutput{flow: o.flows[0]}
	}
	return nil
}

// SingleOutput is a struct that implements Output for a single output flow.
type SingleOutput struct {
	flow *Flow
}

func (o *SingleOutput) Yield(s *State) {
	o.flow.queue.Push(s)
}
