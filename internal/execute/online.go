package execute

import (
	"context"
	"fmt"
	"time"

	"github.com/DuarteMRAlves/maestro/internal"
)

type id int

// onlineState defines a structure to store the state of an online pipeline. It
// is created in a source stage with a unique id, that is transferred throughout
// the pipeline, identifying messages that were generated by applying
// transformations to the original message and allowing for parallel branches to
// be synchronized.
type onlineState struct {
	id  id
	msg internal.Message
}

var emptyOnlineState = newOnlineState(-1, nil)

func newOnlineState(id id, msg internal.Message) onlineState {
	return onlineState{
		id:  id,
		msg: msg,
	}
}

func fromOnlineState(s onlineState, msg internal.Message) onlineState {
	return newOnlineState(s.id, msg)
}

func (s onlineState) String() string {
	return fmt.Sprintf("onlineState{id:%d,msg:%v}", s.id, s.msg)
}

type onlineUnaryStage struct {
	name    internal.StageName
	address internal.Address

	input  <-chan onlineState
	output chan<- onlineState

	clientBuilder internal.UnaryClientBuilder

	logger Logger
}

func newOnlineUnaryStage(
	name internal.StageName,
	input <-chan onlineState,
	output chan<- onlineState,
	address internal.Address,
	clientBuilder internal.UnaryClientBuilder,
	logger Logger,
) *onlineUnaryStage {
	return &onlineUnaryStage{
		name:          name,
		input:         input,
		output:        output,
		address:       address,
		clientBuilder: clientBuilder,
		logger:        logger,
	}
}

func (s *onlineUnaryStage) Run(ctx context.Context) error {
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

func (s *onlineUnaryStage) call(
	ctx context.Context,
	client internal.UnaryClient,
	req internal.Message,
) (internal.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	return client.Call(ctx, req)
}

// onlineSourceStage is the source of the pipeline. It defines the initial ids of
// the states and sends empty messages of the received type.
type onlineSourceStage struct {
	count  int32
	gen    internal.EmptyMessageGen
	output chan<- onlineState
}

func newOnlineSourceStage(
	start int32,
	gen internal.EmptyMessageGen,
	output chan<- onlineState,
) *onlineSourceStage {
	return &onlineSourceStage{
		count:  start,
		gen:    gen,
		output: output,
	}
}

func (s *onlineSourceStage) Run(ctx context.Context) error {
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

type onlineSinkStage struct {
	input <-chan onlineState
}

func newOnlineSinkStage(input <-chan onlineState) *onlineSinkStage {
	return &onlineSinkStage{input: input}
}

func (s *onlineSinkStage) Run(ctx context.Context) error {
	for {
		select {
		case <-s.input:
		case <-ctx.Done():
			return nil
		}
	}
}

type onlineSplitStage struct {
	// fields are the names of the fields of the received message that should
	// be sent through the respective channel. If field is empty, the
	// entire message is sent.
	fields []internal.MessageField
	// input is the channel from which to receive the messages.
	input <-chan onlineState
	// outputs are the several channels where to send messages.
	outputs []chan<- onlineState
}

func newOnlineSplitStage(
	fields []internal.MessageField,
	input <-chan onlineState,
	outputs []chan<- onlineState,
) *onlineSplitStage {
	return &onlineSplitStage{
		fields:  fields,
		input:   input,
		outputs: outputs,
	}
}

func (s *onlineSplitStage) Run(ctx context.Context) error {
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
