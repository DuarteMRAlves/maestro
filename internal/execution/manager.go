package execution

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"sync"
)

// Manager handles the flows that are orchestrated.
type Manager interface {
	// RegisterStage registers a stage to be later included in an orchestration.
	RegisterStage(*api.Stage) error
	// RegisterLink registers a link between two stages. The first
	// stage is the source of the link and the second is the target.
	RegisterLink(*api.Stage, *api.Stage, *api.Link) error
	// RegisterOrchestration registers an orchestration with multiple links.
	RegisterOrchestration(*api.Orchestration) error
}

type manager struct {
	mu          sync.RWMutex
	workers     map[api.StageName]Worker
	inputs      map[api.StageName]*InputCfg
	outputs     map[api.StageName]*InputCfg
	connections map[api.LinkName]*Connection
	flows       map[api.OrchestrationName]*Flow

	reflectionManager rpc.Manager
}

func NewManager(reflectionManager rpc.Manager) Manager {
	return &manager{
		workers:           map[api.StageName]Worker{},
		inputs:            map[api.StageName]*InputCfg{},
		outputs:           map[api.StageName]*InputCfg{},
		connections:       map[api.LinkName]*Connection{},
		flows:             map[api.OrchestrationName]*Flow{},
		reflectionManager: reflectionManager,
	}
}

func (m *manager) RegisterStage(s *api.Stage) error {
	rpc, ok := m.reflectionManager.GetRpc(s.Name)
	if !ok {
		return errdefs.NotFoundWithMsg("Rpc not found for stage %s", s.Name)
	}
	cfg := workerCfgForStage(s, rpc)
	w, err := NewWorker(cfg)
	if err != nil {
		return err
	}

	name := s.Name

	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.workers[name]
	if exists {
		return errdefs.AlreadyExistsWithMsg(
			"Worker for stage %s already exists",
			name,
		)
	}
	m.workers[name] = w

	return err
}

func (m *manager) RegisterLink(
	source *api.Stage,
	target *api.Stage,
	link *api.Link,
) error {
	var (
		ok  bool
		err error
	)

	sourceRpc, ok := m.reflectionManager.GetRpc(source.Name)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"Rpc not found for source %s",
			source.Name,
		)
	}
	targetRpc, ok := m.reflectionManager.GetRpc(target.Name)
	if !ok {
		return errdefs.NotFoundWithMsg(
			"Rpc not found for target %s",
			target.Name,
		)
	}
	fmt.Println("On get input/output")
	sourceOutput := sourceRpc.Output()
	targetInput := targetRpc.Input()

	if link.SourceField != "" {
		sourceOutput, ok = sourceOutput.GetMessageField(link.SourceField)
		if !ok {
			return errdefs.NotFoundWithMsg(
				"field with name %s not found for message %s for source stage "+
					"in link %s",
				link.SourceField,
				sourceRpc.Output().FullyQualifiedName(),
				link.Name,
			)
		}
	}
	if link.TargetField != "" {
		targetInput, ok = targetInput.GetMessageField(link.TargetField)
		if !ok {
			return errdefs.NotFoundWithMsg(
				"field with name %s not found for message %s for target stage "+
					"in link %v",
				link.TargetField,
				targetRpc.Input().FullyQualifiedName(),
				link.Name,
			)
		}
	}
	if !sourceOutput.Compatible(targetInput) {
		return errdefs.InvalidArgumentWithMsg(
			"incompatible message types between source output %s and target"+
				" input %s in link %s",
			sourceOutput.FullyQualifiedName(),
			targetInput.FullyQualifiedName(),
			link.Name,
		)
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

func (m *manager) RegisterOrchestration(o *api.Orchestration) error {
	var exists bool

	m.mu.Lock()
	defer m.mu.Unlock()
	for _, l := range o.Links {
		_, exists = m.connections[l]
		if !exists {
			return errdefs.NotFoundWithMsg("link not registered: %v", l)
		}
	}

	m.flows[o.Name] = NewFlow(o)

	return nil
}

func (m *manager) inputCfgForStage(s *api.Stage) *InputCfg {
	name := s.Name
	cfg, ok := m.inputs[name]
	if !ok {
		cfg = NewInputCfg()
		m.inputs[name] = cfg
	}
	return cfg
}

func (m *manager) outputCfgForStage(s *api.Stage) *InputCfg {
	name := s.Name
	cfg, ok := m.outputs[name]
	if !ok {
		cfg = NewInputCfg()
		m.outputs[name] = cfg
	}
	return cfg
}

func workerCfgForStage(s *api.Stage, rpc rpc.RPC) *WorkerCfg {
	return &WorkerCfg{
		Address: s.Address,
		Rpc:     rpc,
		Input:   nil,
		Output:  nil,
		Done:    nil,
	}
}

func (m *manager) connectionForLink(
	l *api.Link,
) (*Connection, error) {
	var (
		c   *Connection
		ok  bool
		err error
	)
	name := l.Name
	c, ok = m.connections[name]
	if !ok {
		if c, err = NewConnection(l); err != nil {
			return nil, err
		}
		m.connections[name] = c
	}
	return c, nil
}
