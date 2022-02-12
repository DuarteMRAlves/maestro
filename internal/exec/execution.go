package exec

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

// Execution executes an orchestration.
type Execution struct {
	orchestration *api.Orchestration

	stages *StageMap

	term chan struct{}
	errs chan error
	done []<-chan struct{}
}

func (e *Execution) Start() {
	e.term = make(chan struct{})
	e.errs = make(chan error)
	e.done = make([]<-chan struct{}, 0, e.stages.Len())

	go func() {
		err, open := <-e.errs
		if open {
			panic(err)
		}
	}()
	e.stages.Iter(
		func(s Stage) {
			ch := make(chan struct{})
			e.done = append(e.done, ch)
			cfg := &RunCfg{
				term: e.term,
				errs: e.errs,
				done: ch,
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
}
