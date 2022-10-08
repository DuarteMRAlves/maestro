package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal/message"
)

// source is the source of the pipeline. It defines the initial ids of
// the states and sends empty messages of the received type.
type source struct {
	builder message.Builder
	output  chan<- state
}

func newSource(
	gen message.Builder, output chan<- state,
) Stage {
	return &source{
		builder: gen,
		output:  output,
	}
}

func (s *source) Run(ctx context.Context) error {
	for {
		next := newState(s.builder.Build())
		select {
		case s.output <- next:
		case <-ctx.Done():
			close(s.output)
			return nil
		}
	}
}
