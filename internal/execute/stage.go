package execute

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal"
	"time"
)

type Stage interface {
	Run(context.Context) error
}

type unaryStage struct {
	name    internal.StageName
	address internal.Address

	input  <-chan onlineState
	output chan<- onlineState

	clientBuilder internal.UnaryClientBuilder

	logger Logger
}

func newUnaryStage(
	name internal.StageName,
	input <-chan onlineState,
	output chan<- onlineState,
	address internal.Address,
	clientBuilder internal.UnaryClientBuilder,
	logger Logger,
) *unaryStage {
	return &unaryStage{
		name:          name,
		input:         input,
		output:        output,
		address:       address,
		clientBuilder: clientBuilder,
		logger:        logger,
	}
}

func (s *unaryStage) Run(ctx context.Context) error {
	var (
		in, out onlineState
		more    bool
	)
	client, err := s.clientBuilder(s.address)
	if err != nil {
		return err
	}
	defer client.Close()
	s.logger.Infof("'%s': started\n", s.name)
	for {
		select {
		case in, more = <-s.input:
		case <-ctx.Done():
			close(s.output)
			s.logger.Infof("'%s': finished\n", s.name)
			return nil
		}
		// channel is closed
		if !more {
			close(s.output)
			s.logger.Infof("'%s': finished\n", s.name)
			return nil
		}
		s.logger.Debugf("'%s': recv id: %d, msg: %#v\n", s.name, in.id, in.msg)
		req := in.msg
		rep, err := s.call(ctx, client, req)
		if err != nil {
			return err
		}

		out = fromOnlineState(in, rep)
		s.logger.Debugf("'%s': send id: %d, msg: %#v\n", s.name, out.id, out.msg)
		select {
		case s.output <- out:
		case <-ctx.Done():
			close(s.output)
			s.logger.Infof("'%s': finished\n", s.name)
			return nil
		}
	}
}

func (s *unaryStage) call(
	ctx context.Context,
	client internal.UnaryClient,
	req internal.Message,
) (internal.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	return client.Call(ctx, req)
}

// sourceStage is the source of the pipeline. It defines the initial ids of
// the states and sends empty messages of the received type.
type sourceStage struct {
	count  int32
	gen    internal.EmptyMessageGen
	output chan<- onlineState
}

func newSourceStage(
	start int32,
	gen internal.EmptyMessageGen,
	output chan<- onlineState,
) *sourceStage {
	return &sourceStage{
		count:  start,
		gen:    gen,
		output: output,
	}
}

func (s *sourceStage) Run(ctx context.Context) error {
	for {
		next := newOnlineState(id(s.count), s.gen())
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
	input <-chan onlineState
}

func newSinkStage(input <-chan onlineState) *sinkStage {
	return &sinkStage{input: input}
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
	fields []internal.MessageField
	// inputs are the several input channels from which to collect the messages.
	inputs []<-chan onlineState
	// output is the channel used to send messages to the downstream stage.
	output chan<- onlineState
	// gen generates empty messages for the output type.
	gen internal.EmptyMessageGen
	// currId is the current id being constructed.
	currId id
}

func newMergeStage(
	fields []internal.MessageField,
	inputs []<-chan onlineState,
	output chan<- onlineState,
	gen internal.EmptyMessageGen,
) *mergeStage {
	return &mergeStage{
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
		partial   internal.Message
		currState onlineState
		done      bool
	)
	// latest stores the most recent state received from any channel.
	latest := make([]onlineState, 0, len(s.inputs))
	for i := 0; i < len(s.inputs); i++ {
		latest = append(latest, emptyOnlineState)
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
			err := partial.SetField(s.fields[i], currState.msg)
			if err != nil {
				return err
			}
			setFields++
		}
		// All fields from inputs were set. The message can be sent
		if setFields == len(s.inputs) {
			sendState := newOnlineState(s.currId, partial)
			select {
			case s.output <- sendState:
			case <-ctx.Done():
				close(s.output)
				return nil
			}
			s.currId++
			for i := 0; i < len(s.inputs); i++ {
				latest[i] = emptyOnlineState
			}
		}
	}
}

func (s *mergeStage) takeUntilCurrId(
	ctx context.Context,
	input <-chan onlineState,
) (onlineState, bool) {
	for {
		select {
		case st, more := <-input:
			if !more {
				return emptyOnlineState, true
			}
			if st.id >= s.currId {
				return st, false
			}
		case <-ctx.Done():
			return emptyOnlineState, true
		}
	}
}

type splitStage struct {
	// fields are the names of the fields of the received message that should
	// be sent through the respective channel. If field is empty, the
	// entire message is sent.
	fields []internal.MessageField
	// input is the channel from which to receive the messages.
	input <-chan onlineState
	// outputs are the several channels where to send messages.
	outputs []chan<- onlineState
}

func newSplitStage(
	fields []internal.MessageField,
	input <-chan onlineState,
	outputs []chan<- onlineState,
) *splitStage {
	return &splitStage{
		fields:  fields,
		input:   input,
		outputs: outputs,
	}
}

func (s *splitStage) Run(ctx context.Context) error {
	var currState onlineState
	for {
		select {
		case currState = <-s.input:
		case <-ctx.Done():
			for _, c := range s.outputs {
				close(c)
			}
			return nil
		}
		msg := currState.msg
		for i, out := range s.outputs {
			send := msg
			field := s.fields[i]
			if !field.IsEmpty() {
				fieldMsg, err := msg.GetField(field)
				if err != nil {
					return err
				}
				send = fieldMsg
			}
			sendState := newOnlineState(currState.id, send)
			select {
			case out <- sendState:
			case <-ctx.Done():
				for _, c := range s.outputs {
					close(c)
				}
				return nil
			}
		}
	}
}
