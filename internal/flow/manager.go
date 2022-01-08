package flow

import (
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"sync"
)

type Manager struct {
	inputs  sync.Map
	outputs sync.Map
}

func NewManager() *Manager {
	return &Manager{
		inputs:  sync.Map{},
		outputs: sync.Map{},
	}
}

func (m *Manager) Register(
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

	sourceOutputCfg := m.outputCfgForStage(source)
	if err = sourceOutputCfg.register(link); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			sourceOutputCfg.unregisterIfExists(link)
		}
	}()

	targetInputCfg := m.inputCfgForStage(target)
	if err = targetInputCfg.register(link); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			targetInputCfg.unregisterIfExists(link)
		}
	}()

	return nil
}

func (m *Manager) inputCfgForStage(s *stage.Stage) *InputCfg {
	name := s.Name()
	cfg, ok := m.inputs.Load(name)
	if !ok {
		cfg, _ = m.inputs.LoadOrStore(name, newInputCfg())
	}
	return cfg.(*InputCfg)
}

func (m *Manager) outputCfgForStage(s *stage.Stage) *OutputCfg {
	name := s.Name()
	cfg, ok := m.outputs.Load(name)
	if !ok {
		cfg, _ = m.outputs.LoadOrStore(name, newOutputCfg())
	}
	return cfg.(*OutputCfg)
}
