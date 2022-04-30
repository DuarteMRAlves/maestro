package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal"
)

// offlineSourceStage is the source of the pipeline. It defines the initial ids of
// the states and sends empty messages of the received type.
type offlineSourceStage struct {
	gen    internal.EmptyMessageGen
	output chan<- offlineState
}

func newOfflineSourceStage(
	gen internal.EmptyMessageGen,
	output chan<- offlineState,
) *offlineSourceStage {
	return &offlineSourceStage{
		gen:    gen,
		output: output,
	}
}

func (s *offlineSourceStage) Run(ctx context.Context) error {
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
