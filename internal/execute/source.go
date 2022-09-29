package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal/message"
)

// offlineSource is the source of the pipeline. It defines the initial ids of
// the states and sends empty messages of the received type.
type offlineSource struct {
	builder message.Builder
	output  chan<- offlineState
}

func newOfflineSource(
	gen message.Builder, output chan<- offlineState,
) Stage {
	return &offlineSource{
		builder: gen,
		output:  output,
	}
}

func (s *offlineSource) Run(ctx context.Context) error {
	for {
		next := newOfflineState(s.builder.Build())
		select {
		case s.output <- next:
		case <-ctx.Done():
			close(s.output)
			return nil
		}
	}
}
