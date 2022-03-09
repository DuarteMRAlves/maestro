package execute

import (
	"context"
	"github.com/DuarteMRAlves/maestro/internal/domain"
	"github.com/DuarteMRAlves/maestro/internal/invoke"
	"sync"
	"time"
)

type StageInformation interface {
	Name() domain.StageName
	MethodContext() domain.MethodContext
}

type MethodFinder interface {
	Find(domain.MethodContext) invoke.UnaryInvoke
}

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
	)
	for {
		select {
		case in = <-s.input:
		case <-ctx.Done():
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
