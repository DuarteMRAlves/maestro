package flow

import (
	"github.com/DuarteMRAlves/maestro/internal/flow/state"
)

type Input struct {
	send []*state.State
	idx  int
}

func NewInput(states []*state.State) *Input {
	return &Input{send: states, idx: 0}
}

func (i *Input) Next() *state.State {
	if len(i.send) == i.idx {
		return nil
	}
	s := i.send[i.idx]
	i.idx += 1
	return s
}

type Output struct {
	States []*state.State
}

func NewOutput() *Output {
	return &Output{States: make([]*state.State, 0)}
}

func (o *Output) Yield(s *state.State) {
	o.States = append(o.States, s)
}
