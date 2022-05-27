package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal"
)

// offlineSource is the source of the pipeline. It defines the initial ids of
// the states and sends empty messages of the received type.
type offlineSource struct {
	gen    internal.EmptyMessageGen
	output chan<- offlineState
}

func newOfflineSource(
	gen internal.EmptyMessageGen, output chan<- offlineState,
) Stage {
	return &offlineSource{
		gen:    gen,
		output: output,
	}
}

func (s *offlineSource) Run(ctx context.Context) error {
	for {
		next := newOfflineState(s.gen())
		select {
		case s.output <- next:
		case <-ctx.Done():
			close(s.output)
			return nil
		}
	}
}

// onlineSource is the source of the pipeline. It defines the initial ids of
// the states and sends empty messages of the received type.
type onlineSource struct {
	count  int32
	gen    internal.EmptyMessageGen
	output chan<- onlineState
}

func newOnlineSource(
	start int32, gen internal.EmptyMessageGen, output chan<- onlineState,
) Stage {
	return &onlineSource{
		count:  start,
		gen:    gen,
		output: output,
	}
}

func (s *onlineSource) Run(ctx context.Context) error {
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
