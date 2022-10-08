package execute

import (
	"context"

	"github.com/DuarteMRAlves/maestro/internal/compiled"
	"golang.org/x/sync/errgroup"
)

type Stage interface {
	Run(context.Context) error
}

// Execution is an instantiation of a pipeline.
type Execution interface {
	Start()
	Stop() error
}

type execution struct {
	stages map[compiled.StageName]Stage

	runner *runner

	logger Logger
}

func newExecution(stages map[compiled.StageName]Stage, logger Logger) *execution {
	return &execution{stages: stages, logger: logger}
}

func (e *execution) Start() {
	e.runner = newRunner()

	for _, s := range e.stages {
		e.runner.goWithCtx(s.Run)
	}

	e.logger.Debugf("Execution started\n")
}

func (e *execution) Stop() error {
	err := e.runner.cancelAndWait()
	e.logger.Debugf("Execution stopped\n")
	return err
}

type runner struct {
	wg *errgroup.Group

	ctx        context.Context
	cancelFunc context.CancelFunc
}

func newRunner() *runner {
	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)
	return &runner{wg: wg, cancelFunc: cancel, ctx: ctx}
}

func (r *runner) goWithCtx(f func(context.Context) error) {
	r.wg.Go(func() error {
		return f(r.ctx)
	})
}

func (r *runner) cancelAndWait() error {
	r.cancelFunc()
	return r.wg.Wait()
}
