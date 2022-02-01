package execution

import "io"

type MockInput struct {
	send   []*State
	idx    int
	source bool
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

func (i *MockInput) IsSource() bool {
	return i.source
}

type MockOutput struct {
	States []*State
	Sink   bool
}

func NewMockOutput() *MockOutput {
	return &MockOutput{States: make([]*State, 0)}
}

func (o *MockOutput) Yield(s *State) {
	o.States = append(o.States, s)
}

func (o *MockOutput) IsSink() bool {
	return o.Sink
}
