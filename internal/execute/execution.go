package execute

import (
	"context"
	"golang.org/x/sync/errgroup"
)

type execution struct {
	stages *stageMap
	wg     *errgroup.Group
	cancel context.CancelFunc
}

func newExecution(stages *stageMap) *execution {
	return &execution{stages: stages}
}

func (e *execution) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	e.stages.iter(
		func(s Stage) {
			wg.Go(
				func() error {
					return s.Run(ctx)
				},
			)
		},
	)

	e.cancel = cancel
	e.wg = wg
}

func (e *execution) Stop() error {
	e.cancel()
	return e.wg.Wait()
}
