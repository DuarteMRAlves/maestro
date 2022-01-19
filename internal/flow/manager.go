package flow

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/flow/connection"
	"github.com/DuarteMRAlves/maestro/internal/flow/input"
	"github.com/DuarteMRAlves/maestro/internal/flow/output"
	"github.com/DuarteMRAlves/maestro/internal/flow/worker"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"sync"
)

// Manager handles the flows that are orchestrated.
type Manager interface {
	// RegisterLink registers a link between two stages. The first
	// stage is the source of the link and the second is the target.
	RegisterLink(*stage.Stage, *stage.Stage, *link.Link) error
}

type manager struct {
	mu          sync.RWMutex
	workers     map[apitypes.StageName]worker.Worker
	inputs      map[apitypes.StageName]*input.Cfg
	outputs     map[apitypes.StageName]*output.Cfg
	connections map[apitypes.LinkName]*connection.Connection
}

func NewManager() Manager {
	return &manager{
		workers:     map[apitypes.StageName]worker.Worker{},
		inputs:      map[apitypes.StageName]*input.Cfg{},
		outputs:     map[apitypes.StageName]*output.Cfg{},
		connections: map[apitypes.LinkName]*connection.Connection{},
	}
}

func (m *manager) RegisterStage(s *stage.Stage) error {
	cfg := workerCfgForStage(s)

	w, err := worker.NewWorker(cfg)
	if err != nil {
		return err
	}

	name := s.Name()

	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.workers[name]
	if exists {
		return errdefs.AlreadyExistsWithMsg(
			"Worker for stage %s already exists",
			name)
	}
	m.workers[name] = w

	return err
}

func (m *manager) RegisterLink(
	source *stage.Stage,
	target *stage.Stage,
	link *link.Link,
) error {
	var (
		ok  bool
		err error
	)

	sourceOutput := source.Rpc().Output()
	targetInput := target.Rpc().Input()

	if link.SourceField() != "" {
		sourceOutput, ok = sourceOutput.GetMessageField(link.SourceField())
		if !ok {
			return errdefs.NotFoundWithMsg(
				"field with name %s not found for message %s for source stage "+
					"in link %s",
				link.SourceField(),
				source.Rpc().Output().FullyQualifiedName(),
				link.Name())
		}
	}
	if link.TargetField() != "" {
		targetInput, ok = targetInput.GetMessageField(link.TargetField())
		if !ok {
			return errdefs.NotFoundWithMsg(
				"field with name %s not found for message %s for target stage "+
					"in link %v",
				link.TargetField(),
				target.Rpc().Input().FullyQualifiedName(),
				link.Name())
		}
	}
	if !sourceOutput.Compatible(targetInput) {
		return errdefs.InvalidArgumentWithMsg(
			"incompatible message types between source output %s and target"+
				" input %s in link %s",
			sourceOutput.FullyQualifiedName(),
			targetInput.FullyQualifiedName(),
			link.Name())
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	conn, err := m.connectionForLink(link)
	if err != nil {
		return err
	}

	sourceOutputCfg := m.outputCfgForStage(source)
	if err = sourceOutputCfg.Register(conn); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sourceOutputCfg.UnregisterIfExists(conn)
		}
	}()

	targetInputCfg := m.inputCfgForStage(target)
	if err = targetInputCfg.Register(conn); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			targetInputCfg.UnregisterIfExists(conn)
		}
	}()

	return nil
}

func (m *manager) inputCfgForStage(s *stage.Stage) *input.Cfg {
	name := s.Name()
	cfg, ok := m.inputs[name]
	if !ok {
		cfg = input.NewInputCfg()
		m.inputs[name] = cfg
	}
	return cfg
}

func (m *manager) outputCfgForStage(s *stage.Stage) *output.Cfg {
	name := s.Name()
	cfg, ok := m.outputs[name]
	if !ok {
		cfg = output.NewOutputCfg()
		m.outputs[name] = cfg
	}
	return cfg
}

func workerCfgForStage(s *stage.Stage) *worker.Cfg {
	return &worker.Cfg{
		Address: s.Address(),
		Rpc:     s.Rpc(),
		Input:   nil,
		Output:  nil,
		Done:    nil,
	}
}

func (m *manager) connectionForLink(
	l *link.Link,
) (*connection.Connection, error) {
	var (
		c   *connection.Connection
		ok  bool
		err error
	)
	name := l.Name()
	c, ok = m.connections[name]
	if !ok {
		if c, err = connection.New(l); err != nil {
			return nil, err
		}
		m.connections[name] = c
	}
	return c, nil
}
