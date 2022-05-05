package execute

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
}

type Stage interface {
	Run(context.Context) error
}

type Execution struct {
	stages      *stageMap
	chanDrainer chanDrainer
	wg          *errgroup.Group
	cancel      context.CancelFunc

	logger Logger
}

func newExecution(stages *stageMap, drainFunc chanDrainer, logger Logger) *Execution {
	return &Execution{stages: stages, chanDrainer: drainFunc, logger: logger}
}

func (e *Execution) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		e.chanDrainer(ctx)
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
