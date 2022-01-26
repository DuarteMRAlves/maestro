package execution

import (
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/execution/connection"
	"github.com/DuarteMRAlves/maestro/internal/execution/flow"
	"github.com/DuarteMRAlves/maestro/internal/execution/input"
	"github.com/DuarteMRAlves/maestro/internal/execution/output"
	"github.com/DuarteMRAlves/maestro/internal/execution/worker"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"sync"
)

// Manager handles the flows that are orchestrated.
type Manager interface {
	// RegisterStage registers a stage to be later included in an orchestration.
	RegisterStage(*api.Stage) error
	// RegisterLink registers a link between two stages. The first
	// stage is the source of the link and the second is the target.
	RegisterLink(*api.Stage, *api.Stage, *storage.Link) error
	// RegisterOrchestration registers an orchestration with multiple links.
	RegisterOrchestration(*api.Orchestration) error
}

type manager struct {
	mu          sync.RWMutex
	workers     map[api.StageName]worker.Worker
	inputs      map[api.StageName]*input.Cfg
	outputs     map[api.StageName]*output.Cfg
	connections map[api.LinkName]*connection.Connection
	flows       map[api.OrchestrationName]*flow.Flow

	reflectionManager reflection.Manager
}

func NewManager(reflectionManager reflection.Manager) Manager {
	return &manager{
		workers:           map[api.StageName]worker.Worker{},
		inputs:            map[api.StageName]*input.Cfg{},
		outputs:           map[api.StageName]*output.Cfg{},
		connections:       map[api.LinkName]*connection.Connection{},
		flows:             map[api.OrchestrationName]*flow.Flow{},
		reflectionManager: reflectionManager,
	}
}

func (m *manager) RegisterStage(s *api.Stage) error {
	rpc, ok := m.reflectionManager.GetRpc(s.Name)
	if !ok {
		return errdefs.NotFoundWithMsg("Rpc not found for stage %s", s.Name)
	}
	cfg := workerCfgForStage(s, rpc)
	w, err := worker.NewWorker(cfg)
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
	link *storage.Link,
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

	if link.SourceField() != "" {
		sourceOutput, ok = sourceOutput.GetMessageField(link.SourceField())
		if !ok {
			return errdefs.NotFoundWithMsg(
				"field with name %s not found for message %s for source stage "+
					"in link %s",
				link.SourceField(),
				sourceRpc.Output().FullyQualifiedName(),
				link.Name(),
			)
		}
	}
	if link.TargetField() != "" {
		targetInput, ok = targetInput.GetMessageField(link.TargetField())
		if !ok {
			return errdefs.NotFoundWithMsg(
				"field with name %s not found for message %s for target stage "+
					"in link %v",
				link.TargetField(),
				targetRpc.Input().FullyQualifiedName(),
				link.Name(),
			)
		}
	}
	if !sourceOutput.Compatible(targetInput) {
		return errdefs.InvalidArgumentWithMsg(
			"incompatible message types between source output %s and target"+
				" input %s in link %s",
			sourceOutput.FullyQualifiedName(),
			targetInput.FullyQualifiedName(),
			link.Name(),
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

	m.flows[o.Name] = flow.New(o)

	return nil
}

func (m *manager) inputCfgForStage(s *api.Stage) *input.Cfg {
	name := s.Name
	cfg, ok := m.inputs[name]
	if !ok {
		cfg = input.NewCfg()
		m.inputs[name] = cfg
	}
	return cfg
}

func (m *manager) outputCfgForStage(s *api.Stage) *output.Cfg {
	name := s.Name
	cfg, ok := m.outputs[name]
	if !ok {
		cfg = output.NewCfg()
		m.outputs[name] = cfg
	}
	return cfg
}

func workerCfgForStage(s *api.Stage, rpc reflection.RPC) *worker.Cfg {
	return &worker.Cfg{
		Address: s.Address,
		Rpc:     rpc,
		Input:   nil,
		Output:  nil,
		Done:    nil,
	}
}

func (m *manager) connectionForLink(
	l *storage.Link,
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
