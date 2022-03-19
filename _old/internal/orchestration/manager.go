package orchestration

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/storage"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"sync"
)

// GetOrchestrations retrieves stored orchestrations that match the
// received req. The req is an orchestration with the fields that the
// returned orchestrations should have. If a field is empty, then all values
// for that field are accepted.
type GetOrchestrations func(*api.GetOrchestrationRequest) (
	[]*api.Orchestration,
	error,
)

// GetStages retrieves stored stages that match the received req.
// The req is a stage with the fields that the returned stage should have.
// If a field is empty, then all values for that field are accepted.
type GetStages func(*api.GetStageRequest) ([]*api.Stage, error)

// GetLinks retrieves stored links that match the received req.
// The req is a link with the fields that the returned stage should have.
// If a field is empty, then all values for that field are accepted.
type GetLinks func(*api.GetLinkRequest) ([]*api.Link, error)

// GetAssets retrieves stored assets that match the received req.
// The req is an asset with the fields that the returned stage should
// have. If a field is empty, then all values for that field are accepted.
type GetAssets func(*api.GetAssetRequest) ([]*api.Asset, error)

func GetOrchestrationsWithTxn(txn *badger.Txn) GetOrchestrations {
	return func(req *api.GetOrchestrationRequest) (
		[]*api.Orchestration,
		error,
	) {
		var err error

		if req == nil {
			req = &api.GetOrchestrationRequest{}
		}
		filter := buildOrchestrationQueryFilter(req)
		res := make([]*api.Orchestration, 0)

		helper := storage.NewTxnHelper(txn)
		err = helper.IterOrchestrations(
			func(o *api.Orchestration) error {
				if filter(o) {
					orchestrationCp := &api.Orchestration{}
					copyOrchestration(orchestrationCp, o)
					res = append(res, orchestrationCp)
				}
				return nil
			},
			storage.DefaultIterOpts(),
		)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func GetStagesWithTxn(txn *badger.Txn) GetStages {
	return func(req *api.GetStageRequest) ([]*api.Stage, error) {
		var err error

		if req == nil {
			req = &api.GetStageRequest{}
		}
		filter := buildStageQueryFilter(req)
		res := make([]*api.Stage, 0)

		helper := storage.NewTxnHelper(txn)
		err = helper.IterStages(
			func(s *api.Stage) error {
				if filter(s) {
					stageCp := &api.Stage{}
					copyStage(stageCp, s)
					res = append(res, stageCp)
				}
				return nil
			},
			storage.DefaultIterOpts(),
		)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func GetLinksWithTxn(txn *badger.Txn) GetLinks {
	return func(req *api.GetLinkRequest) ([]*api.Link, error) {
		var err error

		if req == nil {
			req = &api.GetLinkRequest{}
		}
		filter := buildLinkQueryFilter(req)
		res := make([]*api.Link, 0)

		helper := storage.NewTxnHelper(txn)
		err = helper.IterLinks(
			func(l *api.Link) error {
				if filter(l) {
					linkCp := &api.Link{}
					copyLink(linkCp, l)
					res = append(res, linkCp)
				}
				return nil
			},
			storage.DefaultIterOpts(),
		)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func GetAssetsWithTxn(txn *badger.Txn) GetAssets {
	return func(req *api.GetAssetRequest) ([]*api.Asset, error) {
		var err error

		if req == nil {
			req = &api.GetAssetRequest{}
		}
		filter := buildAssetQueryFilter(req)
		res := make([]*api.Asset, 0)

		helper := storage.NewTxnHelper(txn)
		err = helper.IterAssets(
			func(a *api.Asset) error {
				if filter(a) {
					assetCp := &api.Asset{}
					copyAsset(assetCp, a)
					res = append(res, assetCp)
				}
				return nil
			},
			storage.DefaultIterOpts(),
		)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

// Manager handles the created executions.
type Manager interface {
	// StartExecution starts the execution associated with the orchestration
	// with the received name. If no Execution for that name exists, a new
	// execution is created and started.
	StartExecution(*badger.Txn, *api.StartExecutionRequest) error
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
	helper := storage.NewTxnHelper(txn)
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