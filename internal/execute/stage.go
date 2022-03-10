package execute

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/invoke"
	"sync"
	"time"
)

type unaryStage struct {
	running bool

	input  <-chan state
	output chan<- state

	gen    invoke.MessageGenerator
	invoke invoke.UnaryInvoke

	mu sync.Mutex
}

func newUnaryStage(
	input <-chan state,
	output chan<- state,
	gen invoke.MessageGenerator,
	invoke invoke.UnaryInvoke,
) *unaryStage {
	return &unaryStage{
		running: false,
		input:   input,
		output:  output,
		gen:     gen,
		invoke:  invoke,
		mu:      sync.Mutex{},
	}
}

func (s *unaryStage) Run(ctx context.Context) error {
	var (
		in, out  state
		req, rep invoke.DynamicMessage
		more     bool
	)
	for {
		select {
		case in, more = <-s.input:
		case <-ctx.Done():
			close(s.output)
			return nil
		}
		// channel is closed
		if !more {
			close(s.output)
			return nil
		}
		req = in.msg
		rep = s.gen()

		err := s.call(ctx, req.GrpcMessage(), rep.GrpcMessage())
		if err != nil {
			return err
		}

		out = updateStateMsg(in, rep)

		select {
		case s.output <- out:
		case <-ctx.Done():
			close(s.output)
			return nil
		}
	}
}

func (s *unaryStage) call(ctx context.Context, req, rep interface{}) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return s.invoke(ctx, req, rep)
}

// sourceStage is the source of the orchestration. It defines the initial ids of
// the states and sends empty messages of the received type.
type sourceStage struct {
	count  int32
	gen    invoke.MessageGenerator
	output chan<- state
}

func newSourceStage(
	start int32,
	gen invoke.MessageGenerator,
	output chan<- state,
) sourceStage {
	return sourceStage{
		count:  start,
		gen:    gen,
		output: output,
	}
}

func (s *sourceStage) Run(ctx context.Context) error {
	for {
		next := newState(id(s.count), s.gen())
		select {
		case s.output <- next:
		case <-ctx.Done():
			close(s.output)
			return nil
		}
		s.count++
	}
}

type sinkStage struct {
	input <-chan state
}

func newSinkStage(input <-chan state) sinkStage {
	return sinkStage{input: input}
}

func (s *sinkStage) Run(ctx context.Context) error {
	for {
		select {
		case <-s.input:
		case <-ctx.Done():
			return nil
		}
	}
}

type mergeStage struct {
	// fields are the names of the fields of the generated message that should
	// be filled with the collected messages.
	fields []domain.MessageField
	// inputs are the several input channels from which to collect the messages.
	inputs []<-chan state
	// output is the channel used to send messages to the downstream stage.
	output chan<- state
	// gen generates empty messages for the output type.
	gen invoke.MessageGenerator
	// currId is the current id being constructed.
	currId id
}

func newMergeStage(
	fields []domain.MessageField,
	inputs []<-chan state,
	output chan<- state,
	gen invoke.MessageGenerator,
) mergeStage {
	return mergeStage{
		fields: fields,
		inputs: inputs,
		output: output,
		gen:    gen,
		currId: 0,
	}
}

func (s *mergeStage) Run(ctx context.Context) error {
	var (
		// partial is the current message being constructed.
		partial   invoke.DynamicMessage
		currState state
		done      bool
	)
	// latest stores the most recent state received from any channel.
	latest := make([]state, 0, len(s.inputs))
	for i := 0; i < len(s.inputs); i++ {
		latest = append(latest, emptyState)
	}
	fieldSetters := make([]invoke.FieldSetter, 0, len(s.fields))
	for _, f := range s.fields {
		fieldSetters = append(fieldSetters, invoke.NewFieldSetter(f))
	}
	for {
		partial = s.gen()
		// number of fields in the partial message that are set.
		setFields := 0
		for i, input := range s.inputs {
			currState = latest[i]
			if currState.id < s.currId {
				currState, done = s.takeUntilCurrId(ctx, input)
				if done {
					close(s.output)
					return nil
				}
				latest[i] = currState
			}
			// If the currState id is higher than the current id, we will never
			// be able to construct the message for the current id. As such, we
			// need to break this cycle, discard the work and move to the next
			// iteration with the currState id.
			if currState.id > s.currId {
				s.currId = currState.id
				break
			}
			res := fieldSetters[i](partial, currState.msg.GrpcMessage())
			if res.IsError() {
				return res.Error()
			}
			partial = res.Unwrap()
			setFields++
		}
		// All fields from inputs were set. The message can be sent
		if setFields == len(s.inputs) {
			sendState := newState(s.currId, partial)
			select {
			case s.output <- sendState:
			case <-ctx.Done():
				close(s.output)
				return nil
			}
			s.currId++
			for i := 0; i < len(s.inputs); i++ {
				latest[i] = emptyState
			}
		}
	}
}

func (s *mergeStage) takeUntilCurrId(
	ctx context.Context,
	input <-chan state,
) (state, bool) {
	for {
		select {
		case st, more := <-input:
			if !more {
				return emptyState, true
			}
			if st.id >= s.currId {
				return st, false
			}
		case <-ctx.Done():
			return emptyState, true
		}
	}
}
