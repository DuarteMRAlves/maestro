package orchestration

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/events"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"sync"
)

// CreateOrchestration creates an orchestration from the given request. The
// function returns an error if the orchestration name is not valid.
type CreateOrchestration func(*api.CreateOrchestrationRequest) error

// GetOrchestrations retrieves stored orchestrations that match the
// received req. The req is an orchestration with the fields that the
// returned orchestrations should have. If a field is empty, then all values
// for that field are accepted.
type GetOrchestrations func(*api.GetOrchestrationRequest) (
	[]*api.Orchestration,
	error,
)

// CreateStage creates a new stage with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
type CreateStage func(*api.CreateStageRequest) error

// GetStages retrieves stored stages that match the received req.
// The req is a stage with the fields that the returned stage should have.
// If a field is empty, then all values for that field are accepted.
type GetStages func(*api.GetStageRequest) ([]*api.Stage, error)

// CreateLink creates a new link with the specified config. It returns an
// error if the link is not created and nil otherwise.
type CreateLink func(*api.CreateLinkRequest) error

// GetLinks retrieves stored links that match the received req.
// The req is a link with the fields that the returned stage should have.
// If a field is empty, then all values for that field are accepted.
type GetLinks func(*api.GetLinkRequest) ([]*api.Link, error)

// CreateAsset creates a new asset with the specified config. It returns an
// error if the asset is not created and nil otherwise.
type CreateAsset func(*api.CreateAssetRequest) error

// GetAssets retrieves stored assets that match the received req.
// The req is an asset with the fields that the returned stage should
// have. If a field is empty, then all values for that field are accepted.
type GetAssets func(*api.GetAssetRequest) ([]*api.Asset, error)

func CreateOrchestrationWithTxn(txn *badger.Txn) CreateOrchestration {
	return func(req *api.CreateOrchestrationRequest) error {
		var err error
		if err = validateCreateOrchestrationConfig(txn, req); err != nil {
			return err
		}

		o := &api.Orchestration{
			Name:   req.Name,
			Phase:  api.OrchestrationPending,
			Stages: []api.StageName{},
			Links:  []api.LinkName{},
		}
		helper := kv.NewTxnHelper(txn)
		return helper.SaveOrchestration(o)
	}
}

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

		helper := kv.NewTxnHelper(txn)
		err = helper.IterOrchestrations(
			func(o *api.Orchestration) error {
				if filter(o) {
					orchestrationCp := &api.Orchestration{}
					copyOrchestration(orchestrationCp, o)
					res = append(res, orchestrationCp)
				}
				return nil
			},
			kv.DefaultIterOpts(),
		)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func CreateStageWithTxn(txn *badger.Txn) CreateStage {
	return func(req *api.CreateStageRequest) error {
		var err error
		helper := kv.NewTxnHelper(txn)
		ctx := newCreateStageContext(req, helper)

		if err = ctx.validateAndComplete(); err != nil {
			return err
		}
		return ctx.persist()
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

		helper := kv.NewTxnHelper(txn)
		err = helper.IterStages(
			func(s *api.Stage) error {
				if filter(s) {
					stageCp := &api.Stage{}
					copyStage(stageCp, s)
					res = append(res, stageCp)
				}
				return nil
			},
			kv.DefaultIterOpts(),
		)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func CreateLinkWithTxn(txn *badger.Txn) CreateLink {
	return func(req *api.CreateLinkRequest) error {
		var err error
		helper := kv.NewTxnHelper(txn)
		ctx := newCreateLinkContext(req, helper)

		if err = ctx.validateAndComplete(); err != nil {
			return err
		}
		return ctx.persist()
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

		helper := kv.NewTxnHelper(txn)
		err = helper.IterLinks(
			func(l *api.Link) error {
				if filter(l) {
					linkCp := &api.Link{}
					copyLink(linkCp, l)
					res = append(res, linkCp)
				}
				return nil
			},
			kv.DefaultIterOpts(),
		)
		if err != nil {
			return nil, err
		}
		return res, nil
	}
}

func CreateAssetWithTxn(txn *badger.Txn) CreateAsset {
	return func(req *api.CreateAssetRequest) error {
		var err error
		if err = validateCreateAssetRequest(txn, req); err != nil {
			return err
		}
		asset := &api.Asset{Name: req.Name, Image: req.Image}
		helper := kv.NewTxnHelper(txn)
		if err = helper.SaveAsset(asset); err != nil {
			return errdefs.InternalWithMsg("persist error: %v", err)
		}
		return nil
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

		helper := kv.NewTxnHelper(txn)
		err = helper.IterAssets(
			func(a *api.Asset) error {
				if filter(a) {
					assetCp := &api.Asset{}
					copyAsset(assetCp, a)
					res = append(res, assetCp)
				}
				return nil
			},
			kv.DefaultIterOpts(),
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
	// AttachExecution attaches to the Execution with the given Orchestration
	// name. It returns a subscription with list of previous events and a
	// channel where new events will be written. The subscription also provides
	// a token that can be used to unsubscribe.
	AttachExecution(*api.AttachExecutionRequest) (*events.Subscription, error)
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
	*events.Subscription,
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
