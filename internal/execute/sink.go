package execute

import "context"

type offlineSink struct {
	input <-chan offlineState
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

type onlineSink struct {
	input <-chan onlineState
}

func (s *onlineSink) Run(ctx context.Context) error {
	for {
		select {
		case <-s.input:
		case <-ctx.Done():
			return nil
		}
	}
}

type sinkBuildFunc[T any] func(input <-chan T) Stage

func offlineSinkBuildFunc() sinkBuildFunc[offlineState] {
	return func(input <-chan offlineState) Stage {
		return &offlineSink{input: input}
	}
}

func onlineSinkBuildFunc() sinkBuildFunc[onlineState] {
	return func(input <-chan onlineState) Stage {
		return &onlineSink{input: input}
	}
}
