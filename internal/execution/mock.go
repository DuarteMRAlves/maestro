package execution

type MockInput struct {
	send   []*State
	idx    int
	source bool
	termFn func()

	ch  chan *State
	end chan struct{}
}

func NewMockInput(states []*State, termFn func()) *MockInput {
	ch := make(chan *State)
	end := make(chan struct{})

	i := &MockInput{send: states, idx: 0, termFn: termFn, ch: ch, end: end}
	go func() {
		defer close(i.ch)
		defer close(i.end)
		for {
			if len(i.send) == i.idx {
				termFn()
				<-i.end
				return
			}
			i.ch <- i.send[i.idx]
			i.idx += 1
		}
	}()
	return i
}

func (i *MockInput) Chan() <-chan *State {
	return i.ch
}

func (i *MockInput) Close() {
	i.end <- struct{}{}
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
