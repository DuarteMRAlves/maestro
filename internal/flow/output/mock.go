package output

import "github.com/DuarteMRAlves/maestro/internal/flow/state"

type MockOutput struct {
	States []*state.State
}

func NewMockOutput() *MockOutput {
	return &MockOutput{States: make([]*state.State, 0)}
}

func (o *MockOutput) Yield(s *state.State) {
	o.States = append(o.States, s)
}
