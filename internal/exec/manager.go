package exec

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"sync"
)

// Manager handles the created executions.
type Manager interface {
	// StartExecution starts the execution associated with the orchestration
	// with the received name. If no Execution for that name exists, a new
	// execution is created and started.
	StartExecution(*badger.Txn, *api.StartExecutionRequest) error
	// AttachExecution attaches to the Execution with the given Orchestration
	// name. It returns a subscription with list of previous events and a
	// channel where new events will be written. The subscription also provides
	// a token that can be used to unsubscribe.
	AttachExecution(*api.AttachExecutionRequest) (*api.Subscription, error)
}

type manager struct {
	mu         sync.RWMutex
	executions map[api.OrchestrationName]*Execution

	logger *zap.Logger
}

func NewManager(logger *zap.Logger) Manager {
	return &manager{
		executions: map[api.OrchestrationName]*Execution{},
		logger:     logger,
	}
}

func (m *manager) StartExecution(
	txn *badger.Txn,
	req *api.StartExecutionRequest,
) error {
	var err error
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.executions[req.Orchestration]
	// Already exists, do nothing.
	if exists {
		return nil
	}

	m.executions[req.Orchestration], err = newBuilder(txn).
		withOrchestration(req.Orchestration).
		withLogger(m.logger).
		build()
	if err != nil {
		return err
	}

	err = m.setRunning(txn, req.Orchestration)
	if err != nil {
		return err
	}

	m.executions[req.Orchestration].Start()
	return nil
}

func (m *manager) setRunning(
	txn *badger.Txn,
	name api.OrchestrationName,
) error {
	var err error
	helper := kv.NewTxnHelper(txn)
	o := &api.Orchestration{}
	err = helper.LoadOrchestration(o, name)
	if err != nil {
		return errdefs.PrependMsg(err, "start execution")
	}
	o.Phase = api.OrchestrationRunning
	err = helper.SaveOrchestration(o)
	if err != nil {
		return errdefs.PrependMsg(err, "start execution")
	}
	for _, sName := range o.Stages {
		s := &api.Stage{}
		err = helper.LoadStage(s, sName)
		if err != nil {
			return errdefs.PrependMsg(err, "start execution")
		}
		s.Phase = api.StageRunning
		err = helper.SaveStage(s)
		if err != nil {
			return errdefs.PrependMsg(err, "start execution")
		}
	}
	return nil
}

func (m *manager) AttachExecution(req *api.AttachExecutionRequest) (
	*api.Subscription,
	error,
) {
	m.mu.Lock()
	defer m.mu.Unlock()

	e, exists := m.executions[req.Orchestration]
	if !exists {
		return nil, errdefs.NotFoundWithMsg(
			"execution '%s' not found",
			req.Orchestration,
		)
	}
	sub := e.Subscribe()
	return sub, nil
}
