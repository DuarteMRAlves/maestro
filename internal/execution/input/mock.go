package input

import (
	"github.com/DuarteMRAlves/maestro/internal/execution/state"
	"io"
)

type MockInput struct {
	send []*state.State
	idx  int
}

func NewMockInput(states []*state.State) *MockInput {
	return &MockInput{send: states, idx: 0}
}

func (i *MockInput) Next() (*state.State, error) {
	if len(i.send) == i.idx {
		return nil, io.EOF
	}
	s := i.send[i.idx]
	i.idx += 1
	return s, nil
}
