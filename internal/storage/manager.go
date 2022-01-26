package storage

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/grpc"
	"time"
)

// Manager manages the storage of created orchestrations.
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
	CreateStage(*badger.Txn, *api.CreateStageRequest) (*api.Stage, error)
	// ContainsStage returns true if the stage exists and false otherwise.
	ContainsStage(*badger.Txn, api.StageName) bool
	// GetStageByName retrieves a stored stage. It returns the stage and true
	// if the stage exists and nil, false otherwise.
	GetStageByName(*badger.Txn, api.StageName) (*api.Stage, bool)
	// GetMatchingStage retrieves stored stages that match the received req.
	// The req is a stage with the fields that the returned stage should have.
	// If a field is empty, then all values for that field are accepted.
	GetMatchingStage(*badger.Txn, *api.GetStageRequest) ([]*api.Stage, error)
	// CreateLink creates a new link with the specified config. It returns an
	// error if the link is not created and nil otherwise.
	CreateLink(*badger.Txn, *api.CreateLinkRequest) (*api.Link, error)
	// ContainsLink returns true if a link with the given name exists and false
	// otherwise.
	ContainsLink(*badger.Txn, api.LinkName) bool
	// GetMatchingLinks retrieves stored links that match the received req.
	// The req is a link with the fields that the returned stage should have.
	// If a field is empty, then all values for that field are accepted.
	GetMatchingLinks(*badger.Txn, *api.GetLinkRequest) ([]*api.Link, error)
	// CreateAsset creates a new asset with the specified config. It returns an
	// error if the asset is not created and nil otherwise.
	CreateAsset(*badger.Txn, *api.CreateAssetRequest) error
	// ContainsAsset returns true if an asset with the given name exists and
	// false otherwise.
	ContainsAsset(*badger.Txn, api.AssetName) bool
	// GetMatchingAssets retrieves stored assets that match the received req.
	// The req is an asset with the fields that the returned stage should
	// have. If a field is empty, then all values for that field are accepted.
	GetMatchingAssets(*badger.Txn, *api.GetAssetRequest) ([]*api.Asset, error)
}

type manager struct {
	reflectionManager reflection.Manager
}

func NewManager(reflectionManager reflection.Manager) Manager {
	return &manager{
		reflectionManager: reflectionManager,
	}
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
	err = persistOrchestration(txn, o)
	if err != nil {
		return errdefs.InternalWithMsg("persist error: %v", err)
	}
	return nil
}

func (m *manager) GetMatchingOrchestration(
	txn *badger.Txn,
	req *api.GetOrchestrationRequest,
) ([]*api.Orchestration, error) {
	var (
		o   api.Orchestration
		cp  []byte
		err error
	)
	if req == nil {
		req = &api.GetOrchestrationRequest{}
	}
	filter := buildOrchestrationQueryFilter(req)
	res := make([]*api.Orchestration, 0)

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	prefix := []byte("orchestration:")
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		cp, err = item.ValueCopy(cp)
		if err != nil {
			return nil, errdefs.InternalWithMsg("read: %v", err)
		}
		err = loadOrchestration(&o, cp)
		if err != nil {
			return nil, errdefs.InternalWithMsg("decoding: %v", err)
		}
		if filter(&o) {
			orchestrationCp := &api.Orchestration{}
			copyOrchestration(orchestrationCp, &o)
			res = append(res, orchestrationCp)
		}
	}
	return res, nil
}

func (m *manager) CreateStage(
	txn *badger.Txn,
	req *api.CreateStageRequest,
) (*api.Stage, error) {
	var err error
	if err = m.validateCreateStageConfig(txn, req); err != nil {
		return nil, err
	}
	address := m.inferStageAddress(req)
	err = m.inferRpc(address, req)
	if err != nil {
		return nil, err
	}
	s := &api.Stage{
		Name:    req.Name,
		Phase:   api.StagePending,
		Service: req.Service,
		Rpc:     req.Rpc,
		Address: address,
		Asset:   req.Asset,
	}
	err = PersistStage(txn, s)
	if err != nil {
		return nil, errdefs.InternalWithMsg("persist error: %v", err)
	}
	return s, nil
}

func (m *manager) ContainsStage(txn *badger.Txn, name api.StageName) bool {
	item, _ := txn.Get(stageKey(name))
	return item != nil
}

func (m *manager) GetStageByName(
	txn *badger.Txn,
	name api.StageName,
) (*api.Stage, bool) {
	var (
		data []byte
		err  error
	)
	item, _ := txn.Get(stageKey(name))
	if item == nil {
		return nil, false
	}
	data, err = item.ValueCopy(nil)
	if err != nil {
		return nil, false
	}
	s := &api.Stage{}
	err = loadStage(s, data)
	if err != nil {
		return nil, false
	}
	return s, true
}

func (m *manager) GetMatchingStage(
	txn *badger.Txn,
	req *api.GetStageRequest,
) ([]*api.Stage, error) {
	var (
		s    api.Stage
		data []byte
		err  error
	)

	if req == nil {
		req = &api.GetStageRequest{}
	}
	filter := buildStageQueryFilter(req)
	res := make([]*api.Stage, 0)

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	prefix := []byte("stage:")

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		data, err = item.ValueCopy(data)
		if err != nil {
			return nil, errdefs.InternalWithMsg("read: %v", err)
		}
		err = loadStage(&s, data)
		if err != nil {
			return nil, errdefs.InternalWithMsg("decoding: %v", err)
		}
		if filter(&s) {
			stageCp := &api.Stage{
				Name:    s.Name,
				Phase:   s.Phase,
				Service: s.Service,
				Rpc:     s.Rpc,
				Address: s.Address,
				Asset:   s.Asset,
			}
			res = append(res, stageCp)
		}
	}
	return res, nil
}

func (m *manager) CreateLink(
	txn *badger.Txn,
	req *api.CreateLinkRequest,
) (*api.Link, error) {
	var err error
	if err = m.validateCreateLinkConfig(txn, req); err != nil {
		return nil, err
	}
	l := &api.Link{
		Name:        req.Name,
		SourceStage: req.SourceStage,
		SourceField: req.SourceField,
		TargetStage: req.TargetStage,
		TargetField: req.TargetField,
	}
	if err = PersistLink(txn, l); err != nil {
		return nil, errdefs.InternalWithMsg("persist error: %v", err)
	}
	return l, nil
}

// ContainsLink returns true if a link with the given name exists and false
// otherwise.
func (m *manager) ContainsLink(txn *badger.Txn, name api.LinkName) bool {
	item, _ := txn.Get(linkKey(name))
	return item != nil
}

func (m *manager) GetMatchingLinks(
	txn *badger.Txn,
	req *api.GetLinkRequest,
) ([]*api.Link, error) {
	var (
		l    api.Link
		data []byte
		err  error
	)

	if req == nil {
		req = &api.GetLinkRequest{}
	}
	filter := buildLinkQueryFilter(req)
	res := make([]*api.Link, 0)

	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	prefix := []byte("link:")

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		data, err = item.ValueCopy(data)
		if err != nil {
			return nil, errdefs.InternalWithMsg("read: %v", err)
		}
		err = loadLink(&l, data)
		if err != nil {
			return nil, errdefs.InternalWithMsg("decoding: %v", err)
		}
		if filter(&l) {
			linkCp := &api.Link{
				Name:        l.Name,
				SourceStage: l.SourceStage,
				SourceField: l.SourceField,
				TargetStage: l.TargetStage,
				TargetField: l.TargetField,
			}
			res = append(res, linkCp)
		}
	}
	return res, nil
}

func (m *manager) inferStageAddress(req *api.CreateStageRequest) string {
	address := req.Address
	// If address is empty, fill it from req host and port.
	if address == "" {
		host, port := req.Host, req.Port
		if host == "" {
			host = "localhost"
		}
		if port == 0 {
			port = 8061
		}
		address = fmt.Sprintf("%s:%d", host, port)
	}
	return address
}

func (m *manager) inferRpc(
	address string,
	req *api.CreateStageRequest,
) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return errdefs.InternalWithMsg(
			"connect to %s for stage %s: %s",
			address,
			req.Name,
			err,
		)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rpcDiscoveryCfg := &reflection.FindQuery{
		Conn:    conn,
		Service: req.Service,
		Rpc:     req.Rpc,
	}
	err = m.reflectionManager.FindRpc(ctx, req.Name, rpcDiscoveryCfg)
	if err != nil {
		return errdefs.PrependMsg(err, "stage %v", req.Name)
	}
	return nil
}

func (m *manager) CreateAsset(
	txn *badger.Txn,
	req *api.CreateAssetRequest,
) error {
	var err error
	if err = validateCreateAssetRequest(req); err != nil {
		return errdefs.InvalidArgumentWithError(err)
	}

	if m.ContainsAsset(txn, req.Name) {
		return errdefs.AlreadyExistsWithMsg(
			"asset '%v' already exists",
			req.Name,
		)
	}
	asset := &api.Asset{Name: req.Name, Image: req.Image}
	if err = PersistAsset(txn, asset); err != nil {
		return errdefs.InternalWithMsg("persist error: %v", err)
	}
	return nil
}

func (m *manager) ContainsAsset(txn *badger.Txn, name api.AssetName) bool {
	item, _ := txn.Get(assetKey(name))
	return item != nil
}

func (m *manager) GetMatchingAssets(
	txn *badger.Txn,
	req *api.GetAssetRequest,
) ([]*api.Asset, error) {
	var (
		asset api.Asset
		cp    []byte
		err   error
	)

	if req == nil {
		req = &api.GetAssetRequest{}
	}
	filter := buildAssetQueryFilter(req)
	res := make([]*api.Asset, 0)
	it := txn.NewIterator(badger.DefaultIteratorOptions)

	defer it.Close()
	prefix := []byte("asset:")

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		cp, err = item.ValueCopy(cp)
		if err != nil {
			return nil, errdefs.InternalWithMsg("read: %v", err)
		}
		err = loadAsset(&asset, cp)
		if err != nil {
			return nil, errdefs.InternalWithMsg("decoding: %v", err)
		}
		if filter(&asset) {
			assetCp := &api.Asset{
				Name:  asset.Name,
				Image: asset.Image,
			}
			res = append(res, assetCp)
		}
	}
	return res, nil
}
