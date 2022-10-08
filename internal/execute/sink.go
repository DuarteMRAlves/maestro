package execute

import "context"

type sink struct {
	input <-chan state
}

func newSink(input <-chan state) Stage {
	return &sink{input: input}
}

func (s *sink) Run(ctx context.Context) error {
	for {
		select {
		case <-s.input:
		case <-ctx.Done():
			return nil
		}
	}
}
