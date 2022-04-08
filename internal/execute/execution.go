package execute

import (
	"context"
	"golang.org/x/sync/errgroup"
	"time"
)

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
}

type Execution struct {
	stages *stageMap
	chans  []chan onlineState
	wg     *errgroup.Group
	cancel context.CancelFunc

	logger Logger
}

func newExecution(stages *stageMap, chans []chan onlineState, logger Logger) *Execution {
	return &Execution{stages: stages, chans: chans, logger: logger}
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
	e.logger.Debugf("Execution started\n")
}

func (e *Execution) Stop() error {
	e.cancel()
	err := e.wg.Wait()
	e.logger.Debugf("Execution stopped\n")
	return err
}
