package arch

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/kv"
	"github.com/dgraph-io/badger/v3"
)

// Manager manages the architecture of created orchestrations.
type Manager interface {
	// CreateOrchestration creates an orchestration from the given request. The
	// function returns an error if the orchestration name is not valid.
	CreateOrchestration(*badger.Txn, *api.CreateOrchestrationRequest) error
	// GetMatchingOrchestration retrieves stored orchestrations that match the
	// received req. The req is an orchestration with the fields that the
	// returned orchestrations should have. If a field is empty, then all values
	// for that field are accepted.
	GetMatchingOrchestration(
		*badger.Txn,
		*api.GetOrchestrationRequest,
	) ([]*api.Orchestration, error)
	// CreateStage creates a new stage with the specified config.
	// It returns an error if the asset can not be created and nil otherwise.
	CreateStage(*badger.Txn, *api.CreateStageRequest) error
	// GetStageByName retrieves a stored stage. It returns the stage and true
	// if the stage exists and nil, false otherwise.
	GetStageByName(*badger.Txn, api.StageName) (*api.Stage, bool)
	// GetMatchingStage retrieves stored stages that match the received req.
	// The req is a stage with the fields that the returned stage should have.
	// If a field is empty, then all values for that field are accepted.
	GetMatchingStage(*badger.Txn, *api.GetStageRequest) ([]*api.Stage, error)
	// CreateLink creates a new link with the specified config. It returns an
	// error if the link is not created and nil otherwise.
	CreateLink(*badger.Txn, *api.CreateLinkRequest) error
	// GetMatchingLinks retrieves stored links that match the received req.
	// The req is a link with the fields that the returned stage should have.
	// If a field is empty, then all values for that field are accepted.
	GetMatchingLinks(*badger.Txn, *api.GetLinkRequest) ([]*api.Link, error)
	// CreateAsset creates a new asset with the specified config. It returns an
	// error if the asset is not created and nil otherwise.
	CreateAsset(*badger.Txn, *api.CreateAssetRequest) error
	// GetMatchingAssets retrieves stored assets that match the received req.
	// The req is an asset with the fields that the returned stage should
	// have. If a field is empty, then all values for that field are accepted.
	GetMatchingAssets(*badger.Txn, *api.GetAssetRequest) ([]*api.Asset, error)
}

type manager struct {
}

// ManagerContext contains configuration to create a new storage manager.
type ManagerContext struct {
	// DB is the underlying db where the data should be stored
	DB *badger.DB
	// CreateDefault specifies whether a default orchestration with the name
	// "default" should be created.
	CreateDefault bool
}

func NewDefaultContext(db *badger.DB) ManagerContext {
	return ManagerContext{
		DB:            db,
		CreateDefault: true,
	}
}

func NewTestContext(db *badger.DB) ManagerContext {
	return ManagerContext{
		DB:            db,
		CreateDefault: false,
	}
}

func NewManager(ctx ManagerContext) (Manager, error) {
	if ctx.CreateDefault {
		err := ctx.DB.Update(
			func(txn *badger.Txn) error {
				helper := kv.NewTxnHelper(txn)
				if !helper.ContainsOrchestration(DefaultOrchestrationName) {
					return helper.SaveOrchestration(defaultOrchestration())
				}
				return nil
			},
		)
		if err != nil {
			return nil, errdefs.PrependMsg(err, "new storage manager")
		}
	}
	m := &manager{}
	return m, nil
}

func (m *manager) CreateOrchestration(
	txn *badger.Txn,
	req *api.CreateOrchestrationRequest,
) error {
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

func (m *manager) GetMatchingOrchestration(
	txn *badger.Txn,
	req *api.GetOrchestrationRequest,
) ([]*api.Orchestration, error) {
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

func (m *manager) CreateStage(
	txn *badger.Txn,
	req *api.CreateStageRequest,
) error {
	var err error
	helper := kv.NewTxnHelper(txn)
	ctx := newCreateStageContext(req, helper)

	if err = ctx.validateAndComplete(); err != nil {
		return err
	}
	return ctx.persist()
}

func (m *manager) GetStageByName(
	txn *badger.Txn,
	name api.StageName,
) (*api.Stage, bool) {
	helper := kv.NewTxnHelper(txn)
	s := &api.Stage{}
	err := helper.LoadStage(s, name)
	if err != nil {
		return nil, false
	}
	return s, true
}

func (m *manager) GetMatchingStage(
	txn *badger.Txn,
	req *api.GetStageRequest,
) ([]*api.Stage, error) {
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

func (m *manager) CreateLink(
	txn *badger.Txn,
	req *api.CreateLinkRequest,
) error {
	var err error
	helper := kv.NewTxnHelper(txn)
	ctx := newCreateLinkContext(req, helper)

	if err = ctx.validateAndComplete(); err != nil {
		return err
	}
	return ctx.persist()
}

func (m *manager) GetMatchingLinks(
	txn *badger.Txn,
	req *api.GetLinkRequest,
) ([]*api.Link, error) {
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

func (m *manager) CreateAsset(
	txn *badger.Txn,
	req *api.CreateAssetRequest,
) error {
	var err error
	if err = m.validateCreateAssetRequest(txn, req); err != nil {
		return err
	}
	asset := &api.Asset{Name: req.Name, Image: req.Image}
	helper := kv.NewTxnHelper(txn)
	if err = helper.SaveAsset(asset); err != nil {
		return errdefs.InternalWithMsg("persist error: %v", err)
	}
	return nil
}

func (m *manager) GetMatchingAssets(
	txn *badger.Txn,
	req *api.GetAssetRequest,
) ([]*api.Asset, error) {
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
