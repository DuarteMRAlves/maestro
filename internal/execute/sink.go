package execute

import "context"

type offlineSink struct {
	input <-chan offlineState
}

func newOfflineSink(input <-chan offlineState) Stage {
	return &offlineSink{input: input}
}

func (s *offlineSink) Run(ctx context.Context) error {
	for {
		select {
		case <-s.input:
		case <-ctx.Done():
			return nil
		}
	}
}
