package execution

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/dgraph-io/badger/v3"
	"sync"
)

// Manager handles the created executions.
type Manager interface {
	// StartExecution starts the execution associated with the orchestration
	// with the received name. If no Execution for that name exists, a new
	// execution is created and started.
	StartExecution(*badger.Txn, api.StartOrchestrationRequest) error
}

type manager struct {
	mu         sync.RWMutex
	executions map[api.OrchestrationName]*Execution

	reflectionManager rpc.Manager
}

func NewManager(reflectionManager rpc.Manager) Manager {
	return &manager{
		executions:        map[api.OrchestrationName]*Execution{},
		reflectionManager: reflectionManager,
	}
}

func (m *manager) StartExecution(
	txn *badger.Txn,
	req api.StartOrchestrationRequest,
) error {
	var (
		name api.OrchestrationName
		err  error
	)

	m.mu.Lock()
	defer m.mu.Unlock()

	name = req.Orchestration
	if name == "" {
		name = storage.DefaultOrchestrationName
	}

	_, exists := m.executions[name]
	// Already exists, do nothing.
	if exists {
		return nil
	}

	txnHelper := storage.NewTxnHelper(txn)
	orchestration := &api.Orchestration{}

	err = txnHelper.LoadOrchestration(orchestration, req.Orchestration)
	if err != nil {
		return errdefs.PrependMsg(err, "start execution")
	}

	m.executions[name] = NewExecution(orchestration)

	return nil
}
