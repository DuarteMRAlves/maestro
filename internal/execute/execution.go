package execute

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"
)

type Execution struct {
	stages *stageMap
	chans  []chan state
	wg     *errgroup.Group
	cancel context.CancelFunc
}

func newExecution(stages *stageMap, chans []chan state) *Execution {
	return &Execution{stages: stages, chans: chans}
}

func (e *Execution) Start() {
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

func (e *Execution) Stop() error {
	e.cancel()
	return e.wg.Wait()
}
