package execute

import (
	"context"
	"time"

	"golang.org/x/sync/errgroup"
)

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
}

type Stage interface {
	Run(context.Context) error
}

type Execution[T any] struct {
	stages *stageMap
	chans  []chan T
	wg     *errgroup.Group
	cancel context.CancelFunc

	logger Logger
}

func newExecution[T any](stages *stageMap, chans []chan T, logger Logger) *Execution[T] {
	return &Execution[T]{stages: stages, chans: chans, logger: logger}
}

func (e *Execution[T]) Start() {
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

func (e *Execution[T]) Stop() error {
	e.cancel()
	err := e.wg.Wait()
	e.logger.Debugf("Execution stopped\n")
	return err
}
