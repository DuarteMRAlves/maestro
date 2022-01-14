package output

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/flow/flow"
	"github.com/DuarteMRAlves/maestro/internal/flow/state"
)

// Output receives the output flow.State for a given stage and sends it to the
// next stages.
type Output interface {
	Yield(s *state.State)
}

// Cfg represents the several output flows for a stage
type Cfg struct {
	typ   Type
	flows []*flow.Flow
}

// Type defines the type of output the stage.Stage associated with this
// Cfg is expecting.
type Type string

const (
	// OutputInfer means the output type is not specified and should be inferred
	// from the received connections.
	OutputInfer Type = "Infer"
	// OutputSingle means the stage only sends its entire output to another
	// stage.
	OutputSingle Type = "Single"
	// OutputSink means the stage outputs the empty message and its output
	// should be discarded as it does not connect to any other stage.
	OutputSink Type = "Sink"
	// OutputSplit means the stage output should be split into multiple messages
	// sent to different stages.
	OutputSplit Type = "Split"
	// OutputDuplicate means the stage output should be duplicated and sent as
	// input to multiple stages.
	OutputDuplicate Type = "Duplicate"
)

func NewOutputCfg() *Cfg {
	return &Cfg{
		typ:   OutputInfer,
		flows: []*flow.Flow{},
	}
}

func (o *Cfg) Register(f *flow.Flow) error {
	l := f.Link
	for _, prev := range o.flows {
		if l.Name() == prev.Link.Name() {
			return errdefs.InvalidArgumentWithMsg(
				"Link with an equal name already registered: %s",
				prev.Link.Name())
		}
	}

	o.flows = append(o.flows, f)
	return nil
}

func (o *Cfg) UnregisterIfExists(search *flow.Flow) {
	idx := -1
	for i, f := range o.flows {
		if f.Link.Name() == search.Link.Name() {
			idx = i
			break
		}
	}
	if idx != -1 {
		o.flows[idx] = o.flows[len(o.flows)-1]
		o.flows = o.flows[:len(o.flows)-1]
	}
}

func (o *Cfg) ToFlow() Output {
	switch len(o.flows) {
	case 1:
		return &SingleOutput{flow: o.flows[0]}
	}
	return nil
}

// SingleOutput is a struct that implements Output for a single output flow.
type SingleOutput struct {
	flow *flow.Flow
}

func (o *SingleOutput) Yield(s *state.State) {
	o.flow.Queue.Push(s)
}
