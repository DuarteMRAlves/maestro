package execution

import "io"

type MockInput struct {
	send []*State
	idx  int
}

func NewMockInput(states []*State) *MockInput {
	return &MockInput{send: states, idx: 0}
}

func (i *MockInput) Next() (*State, error) {
	if len(i.send) == i.idx {
		return nil, io.EOF
	}
	s := i.send[i.idx]
	i.idx += 1
	return s, nil
}

type MockOutput struct {
	States []*State
}

func NewMockOutput() *MockOutput {
	return &MockOutput{States: make([]*State, 0)}
}

func (o *MockOutput) Yield(s *State) {
	o.States = append(o.States, s)
}
