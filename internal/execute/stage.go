package execute

import (
	"context"
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
