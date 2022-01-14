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
	inputs  sync.Map
	outputs sync.Map
	flows   sync.Map
}

func NewManager() Manager {
	return &manager{
		inputs:  sync.Map{},
		outputs: sync.Map{},
		flows:   sync.Map{},
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

	flow, err := m.flowForLink(link)
	if err != nil {
		return err
	}

	sourceOutputCfg := m.outputCfgForStage(source)
	if err = sourceOutputCfg.Register(flow); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sourceOutputCfg.UnregisterIfExists(flow)
		}
	}()

	targetInputCfg := m.inputCfgForStage(target)
	if err = targetInputCfg.Register(flow); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			targetInputCfg.UnregisterIfExists(flow)
		}
	}()

	return nil
}

func (m *manager) inputCfgForStage(s *stage.Stage) *input.Cfg {
	name := s.Name()
	cfg, ok := m.inputs.Load(name)
	if !ok {
		cfg, _ = m.inputs.LoadOrStore(name, input.NewInputCfg())
	}
	return cfg.(*input.Cfg)
}

func (m *manager) outputCfgForStage(s *stage.Stage) *output.Cfg {
	name := s.Name()
	cfg, ok := m.outputs.Load(name)
	if !ok {
		cfg, _ = m.outputs.LoadOrStore(name, output.NewOutputCfg())
	}
	return cfg.(*output.Cfg)
}

func workerCfgForStage(s *stage.Stage) *worker.Cfg {
	return &worker.Cfg{
		Address: s.Address(),
		Rpc:     s.Rpc(),
		Input:   nil,
		Output:  nil,
		Done:    nil,
		MaxMsg:  0,
	}
}

func (m *manager) flowForLink(l *link.Link) (*flow.Flow, error) {
	var err error
	name := l.Name()
	f, ok := m.flows.Load(name)
	if !ok {
		if f, err = flow.NewFlow(l); err != nil {
			return nil, err
		}
		f, _ = m.flows.LoadOrStore(name, f)
	}
	return f.(*flow.Flow), nil
}
