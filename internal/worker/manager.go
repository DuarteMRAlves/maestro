package worker

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"sync"
)

type Manager interface {
	// Run creates and runs a new Worker associated with the stage.Stage with
	// the given name, creating worker from the given cfg.
	Run(stage apitypes.StageName, cfg *Cfg) error
}

type manager struct {
	workers map[apitypes.StageName]Worker
	mu      sync.RWMutex
}

func (m *manager) Run(stage apitypes.StageName, cfg *Cfg) error {
	w, err := NewWorker(cfg)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.workers[stage]
	if exists {
		return errdefs.AlreadyExistsWithMsg(
			"Worker for stage %s already exists",
			stage)
	}
	m.workers[stage] = w
	go w.Run()
	return nil
}
