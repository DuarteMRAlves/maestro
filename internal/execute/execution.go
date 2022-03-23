package execute

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"
)

type execution struct {
	stages *stageMap
	chans  []chan state
	wg     *errgroup.Group
	cancel context.CancelFunc
}

func newExecution(stages *stageMap, chans []chan state) *execution {
	return &execution{stages: stages, chans: chans}
}

func (e *execution) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	chanDrainer := newChanDrainer(5*time.Millisecond, e.chans...)
	wg.Go(func() error {
		chanDrainer(ctx)
		return nil
	})

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
