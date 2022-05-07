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

// Execution is an instantiation of a pipeline.
type Execution interface {
	Start()
	Stop() error
}

type offlineExecution struct {
	stages *stageMap

	runner *runner

	logger Logger
}

func newOfflineExecution(stages *stageMap, logger Logger) *offlineExecution {
	return &offlineExecution{stages: stages, logger: logger}
}

func (e *offlineExecution) Start() {
	e.runner = newRunner()

	e.stages.iter(func(s Stage) {
		e.runner.goWithCtx(s.Run)
	})

	e.logger.Debugf("Execution started\n")
}

func (e *offlineExecution) Stop() error {
	err := e.runner.cancelAndWait()
	e.logger.Debugf("Execution stopped\n")
	return err
}

type onlineExecution struct {
	stages      *stageMap
	chanDrainer chanDrainer

	runner *runner

	logger Logger
}

func newOnlineExecution(stages *stageMap, drainFunc chanDrainer, logger Logger) *onlineExecution {
	return &onlineExecution{stages: stages, chanDrainer: drainFunc, logger: logger}
}

func (e *onlineExecution) Start() {
	e.runner = newRunner()

	e.runner.goWithCtx(func(ctx context.Context) error {
		e.chanDrainer(ctx)
		return nil
	})

	e.stages.iter(func(s Stage) {
		e.runner.goWithCtx(s.Run)
	})

	e.logger.Debugf("Execution started\n")
}

func (e *onlineExecution) Stop() error {
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
