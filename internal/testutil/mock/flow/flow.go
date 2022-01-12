package flow

import "github.com/DuarteMRAlves/maestro/internal/flow"

type Input struct {
	send []*flow.State
	idx  int
}

func NewInput(states []*flow.State) *Input {
	return &Input{send: states, idx: 0}
}

func (i *Input) Next() *flow.State {
	if len(i.send) == i.idx {
		return nil
	}
	s := i.send[i.idx]
	i.idx += 1
	return s
}

type Output struct {
	States []*flow.State
}

func NewOutput() *Output {
	return &Output{States: make([]*flow.State, 0)}
}

func (o *Output) Yield(s *flow.State) {
	o.States = append(o.States, s)
}
