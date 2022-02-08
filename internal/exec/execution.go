package exec

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
)

// Execution executes an orchestration.
type Execution struct {
	orchestration *api.Orchestration

	stages map[api.StageName]Stage

	term chan struct{}
	errs chan error
	done map[api.StageName]<-chan struct{}
}

func (e *Execution) Start() {
	e.term = make(chan struct{})
	e.errs = make(chan error)
	e.done = make(map[api.StageName]<-chan struct{}, len(e.stages))

	go func() {
		err, open := <-e.errs
		if open {
			panic(err)
		}
	}()
	for n, s := range e.stages {
		ch := make(chan struct{})
		e.done[n] = ch
		cfg := &RunCfg{
			term: e.term,
			errs: e.errs,
			done: ch,
		}
		go s.Run(cfg)
	}
}

func (e *Execution) Stop() {
	close(e.term)
	for _, d := range e.done {
		<-d
	}
	close(e.errs)
}
