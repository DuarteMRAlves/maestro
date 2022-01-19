package flow

import (
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/flow/flow"
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
	mu      sync.RWMutex
	workers map[apitypes.StageName]worker.Worker
	inputs  map[apitypes.StageName]*input.Cfg
	outputs map[apitypes.StageName]*output.Cfg
	flows   map[apitypes.LinkName]*flow.Flow
}

func NewManager() Manager {
	return &manager{
		workers: map[apitypes.StageName]worker.Worker{},
		inputs:  map[apitypes.StageName]*input.Cfg{},
		outputs: map[apitypes.StageName]*output.Cfg{},
		flows:   map[apitypes.LinkName]*flow.Flow{},
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

	linkFlow, err := m.flowForLink(link)
	if err != nil {
		return err
	}

	sourceOutputCfg := m.outputCfgForStage(source)
	if err = sourceOutputCfg.Register(linkFlow); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sourceOutputCfg.UnregisterIfExists(linkFlow)
		}
	}()

	targetInputCfg := m.inputCfgForStage(target)
	if err = targetInputCfg.Register(linkFlow); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			targetInputCfg.UnregisterIfExists(linkFlow)
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

func (m *manager) flowForLink(l *link.Link) (*flow.Flow, error) {
	var (
		f   *flow.Flow
		ok  bool
		err error
	)
	name := l.Name()
	f, ok = m.flows[name]
	if !ok {
		if f, err = flow.NewFlow(l); err != nil {
			return nil, err
		}
		m.flows[name] = f
	}
	return f, nil
}
