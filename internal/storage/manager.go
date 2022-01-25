package storage

import (
	"context"
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/dgraph-io/badger/v3"
	"google.golang.org/grpc"
	"time"
)

// Manager manages the storage of created orchestrations.
type Manager interface {
	// CreateOrchestration creates an orchestration from the given config. The
	// function returns an error if the orchestration name is not valid.
	CreateOrchestration(*badger.Txn, *apitypes.Orchestration) error
	// GetMatchingOrchestration retrieves stored orchestrations that match the
	// received query. The query is an orchestration with the fields that the
	// returned orchestrations should have. If a field is empty, then all values
	// for that field are accepted.
	GetMatchingOrchestration(
		*badger.Txn,
		*apitypes.Orchestration,
	) ([]*apitypes.Orchestration, error)
	// CreateStage creates a new stage with the specified config.
	// It returns an error if the asset can not be created and nil otherwise.
	CreateStage(*badger.Txn, *apitypes.Stage) (*Stage, error)
	// ContainsStage returns true if the stage exists and false otherwise.
	ContainsStage(*badger.Txn, apitypes.StageName) bool
	// GetStageByName retrieves a stored stage. It returns the stage and true
	// if the stage exists and nil, false otherwise.
	GetStageByName(*badger.Txn, apitypes.StageName) (*Stage, bool)
	// GetMatchingStage retrieves stored stages that match the received query.
	// The query is a stage with the fields that the returned stage should have.
	// If a field is empty, then all values for that field are accepted.
	GetMatchingStage(*badger.Txn, *apitypes.Stage) ([]*apitypes.Stage, error)
	// CreateLink creates a new link with the specified config. It returns an
	// error if the link is not created and nil otherwise.
	CreateLink(*badger.Txn, *apitypes.Link) (*Link, error)
	// ContainsLink returns true if a link with the given name exists and false
	// otherwise.
	ContainsLink(*badger.Txn, apitypes.LinkName) bool
	// GetMatchingLinks retrieves stored links that match the received query.
	// The query is a link with the fields that the returned stage should have.
	// If a field is empty, then all values for that field are accepted.
	GetMatchingLinks(*badger.Txn, *apitypes.Link) ([]*apitypes.Link, error)
	// CreateAsset creates a new asset with the specified config. It returns an
	// error if the asset is not created and nil otherwise.
	CreateAsset(*badger.Txn, *apitypes.Asset) error
	// ContainsAsset returns true if an asset with the given name exists and
	// false otherwise.
	ContainsAsset(*badger.Txn, apitypes.AssetName) bool
	// GetMatchingAssets retrieves stored assets that match the received query.
	// The query is an asset with the fields that the returned stage should
	// have. If a field is empty, then all values for that field are accepted.
	GetMatchingAssets(*badger.Txn, *apitypes.Asset) ([]*apitypes.Asset, error)
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
	cfg *apitypes.Orchestration,
) error {
	var err error
	if err = validateCreateOrchestrationConfig(txn, cfg); err != nil {
		return err
	}

	o := NewOrchestration(cfg.Name)
	err = persistOrchestration(txn, o)
	if err != nil {
		return errdefs.InternalWithMsg("persist error: %v", err)
	}
	return nil
}

func (m *manager) GetMatchingOrchestration(
	txn *badger.Txn,
	query *apitypes.Orchestration,
) ([]*apitypes.Orchestration, error) {
	var (
		o   Orchestration
		cp  []byte
		err error
	)
	if query == nil {
		query = &apitypes.Orchestration{}
	}
	filter := buildOrchestrationQueryFilter(query)
	res := make([]*apitypes.Orchestration, 0)

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
		fmt.Println("outside")
		fmt.Println(o)
		if filter(&o) {
			fmt.Println("inside")
			fmt.Println(o)
			res = append(res, o.ToApi())
		}
	}
	return res, nil
}

func (m *manager) CreateStage(
	txn *badger.Txn,
	cfg *apitypes.Stage,
) (*Stage, error) {
	var err error
	if err = m.validateCreateStageConfig(txn, cfg); err != nil {
		return nil, err
	}
	address := m.inferStageAddress(cfg)
	err = m.inferRpc(address, cfg)
	if err != nil {
		return nil, err
	}
	spec := &RpcSpec{
		address: address,
		service: cfg.Service,
		rpc:     cfg.Rpc,
	}
	s := NewStage(cfg.Name, spec, cfg.Asset, nil)
	err = PersistStage(txn, s)
	if err != nil {
		return nil, errdefs.InternalWithMsg("persist error: %v", err)
	}
	return s, nil
}

func (m *manager) ContainsStage(txn *badger.Txn, name apitypes.StageName) bool {
	item, _ := txn.Get(stageKey(name))
	return item != nil
}

func (m *manager) GetStageByName(
	txn *badger.Txn,
	name apitypes.StageName,
) (*Stage, bool) {
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
	s := &Stage{}
	s.rpcSpec = &RpcSpec{}
	err = loadStage(s, data)
	if err != nil {
		return nil, false
	}
	return s, true
}

func (m *manager) GetMatchingStage(
	txn *badger.Txn,
	query *apitypes.Stage,
) ([]*apitypes.Stage, error) {
	var (
		s    Stage
		data []byte
		err  error
	)
	s.rpcSpec = &RpcSpec{}

	if query == nil {
		query = &apitypes.Stage{}
	}
	filter := buildStageQueryFilter(query)
	res := make([]*apitypes.Stage, 0)

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
			res = append(res, s.ToApi())
		}
	}
	return res, nil
}

func (m *manager) CreateLink(
	txn *badger.Txn,
	cfg *apitypes.Link,
) (*Link, error) {
	var err error
	if err = m.validateCreateLinkConfig(txn, cfg); err != nil {
		return nil, err
	}
	l := NewLink(
		cfg.Name,
		cfg.SourceStage,
		cfg.SourceField,
		cfg.TargetStage,
		cfg.TargetField,
	)
	if err = PersistLink(txn, l); err != nil {
		return nil, errdefs.InternalWithMsg("persist error: %v", err)
	}
	return l, nil
}

// ContainsLink returns true if a link with the given name exists and false
// otherwise.
func (m *manager) ContainsLink(txn *badger.Txn, name apitypes.LinkName) bool {
	item, _ := txn.Get(linkKey(name))
	return item != nil
}

func (m *manager) GetMatchingLinks(
	txn *badger.Txn,
	query *apitypes.Link,
) ([]*apitypes.Link, error) {
	var (
		l    Link
		data []byte
		err  error
	)

	if query == nil {
		query = &apitypes.Link{}
	}
	filter := buildLinkQueryFilter(query)
	res := make([]*apitypes.Link, 0)

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
			res = append(res, l.ToApi())
		}
	}
	return res, nil
}

func (m *manager) inferStageAddress(cfg *apitypes.Stage) string {
	address := cfg.Address
	// If address is empty, fill it from cfg host and port.
	if address == "" {
		host, port := cfg.Host, cfg.Port
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
	cfg *apitypes.Stage,
) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return errdefs.InternalWithMsg(
			"connect to %s for stage %s: %s",
			address,
			cfg.Name,
			err,
		)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rpcDiscoveryCfg := &reflection.FindQuery{
		Conn:    conn,
		Service: cfg.Service,
		Rpc:     cfg.Rpc,
	}
	err = m.reflectionManager.FindRpc(ctx, cfg.Name, rpcDiscoveryCfg)
	if err != nil {
		return errdefs.PrependMsg(err, "stage %v", cfg.Name)
	}
	return nil
}

func (m *manager) CreateAsset(txn *badger.Txn, cfg *apitypes.Asset) error {
	var err error
	if err = validateCreateAssetConfig(cfg); err != nil {
		return errdefs.InvalidArgumentWithError(err)
	}

	if m.ContainsAsset(txn, cfg.Name) {
		return errdefs.AlreadyExistsWithMsg(
			"asset '%v' already exists",
			cfg.Name,
		)
	}
	asset := NewAsset(cfg.Name, cfg.Image)
	if err = PersistAsset(txn, asset); err != nil {
		return errdefs.InternalWithMsg("persist error: %v", err)
	}
	return nil
}

func (m *manager) ContainsAsset(txn *badger.Txn, name apitypes.AssetName) bool {
	item, _ := txn.Get(assetKey(name))
	return item != nil
}

func (m *manager) GetMatchingAssets(
	txn *badger.Txn,
	query *apitypes.Asset,
) ([]*apitypes.Asset, error) {
	var (
		asset Asset
		cp    []byte
		err   error
	)

	if query == nil {
		query = &apitypes.Asset{}
	}
	filter := buildQueryFilter(query)
	res := make([]*apitypes.Asset, 0)
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
			res = append(res, asset.ToApi())
		}
	}
	return res, nil
}
