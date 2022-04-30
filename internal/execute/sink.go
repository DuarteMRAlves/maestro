package execute

import "context"

type offlineSinkStage struct {
	input <-chan offlineState
}

func newOfflineSinkStage(input <-chan offlineState) *offlineSinkStage {
	return &offlineSinkStage{input: input}
}

func (s *offlineSinkStage) Run(ctx context.Context) error {
	for {
		select {
		case <-s.input:
		case <-ctx.Done():
			return nil
		}
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
