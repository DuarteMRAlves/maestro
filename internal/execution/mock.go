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
	States chan *State
	Sink   bool

	ch  chan *State
	end chan struct{}
}

func NewMockOutput(expected int) *MockOutput {
	ch := make(chan *State)
	end := make(chan struct{})

	o := &MockOutput{States: make(chan *State, expected), ch: ch, end: end}

	go func() {
		defer close(o.ch)
		defer close(o.end)
		defer close(o.States)
		for {
			select {
			case s := <-o.ch:
				o.States <- s
			case <-o.end:
				return
			}
		}
	}()

	return o
}

func (o *MockOutput) Chan() chan<- *State {
	return o.ch
}

func (o *MockOutput) Close() {
	o.end <- struct{}{}
}

func (o *MockOutput) IsSink() bool {
	return o.Sink
}
