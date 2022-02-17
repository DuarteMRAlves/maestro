package exec

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/events"
	"go.uber.org/zap"
)

// Execution executes an orchestration.
type Execution struct {
	orchestration *api.Orchestration

	stages *StageMap
	pubSub events.PubSub

	logger *zap.Logger

	// Runtime structs should not be specified in initialization
	term chan struct{}
	errs chan error
	done []<-chan struct{}
	// cleanup is a channel to signal any cleanup required after all the stages
	// have finished.
	cleanup chan struct{}
}

func (e *Execution) Start() {
	e.term = make(chan struct{})
	e.errs = make(chan error)
	e.done = make([]<-chan struct{}, 0, e.stages.Len())
	e.cleanup = make(chan struct{})

	go func() {
		err, open := <-e.errs
		if open {
			panic(err)
		}
	}()
	go func() {
		logSub := e.pubSub.Subscribe()
		for {
			select {
			case <-e.cleanup:
				err := e.pubSub.Unsubscribe(logSub.Token)
				if err != nil {
					e.errs <- errdefs.PrependMsg(err, "execution log")
				}
				return
			case event := <-logSub.Future:
				e.logger.Info(
					event.Description,
					zap.Time("time", event.Timestamp),
				)
			}
		}
	}()
	e.stages.Iter(
		func(s Stage) {
			ch := make(chan struct{})
			e.done = append(e.done, ch)
			cfg := &RunCfg{
				pubSub: e.pubSub,
				term:   e.term,
				errs:   e.errs,
				done:   ch,
			}
			go s.Run(cfg)
		},
	)
}

func (e *Execution) Stop() {
	close(e.term)
	for _, d := range e.done {
		<-d
	}
	close(e.errs)
	e.stages.Iter(
		func(s Stage) {
			s.Close()
		},
	)
	close(e.cleanup)
}

func (e *Execution) Subscribe() *api.Subscription {
	return e.pubSub.Subscribe()
}
