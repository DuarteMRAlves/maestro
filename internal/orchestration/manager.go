package orchestration

import (
	"context"
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/discovery"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"google.golang.org/grpc"
	"sync"
	"time"
)

// Manager manages the storage of created orchestrations.
type Manager interface {
	// CreateOrchestration creates an orchestration from the given config. The
	// function returns an error if the orchestration name is not valid.
	CreateOrchestration(*apitypes.Orchestration) error
	// GetMatchingOrchestration retrieves stored orchestrations that match the
	// received query. The query is an orchestration with the fields that the
	// returned orchestrations should have. If a field is empty, then all values
	// for that field are accepted.
	GetMatchingOrchestration(*apitypes.Orchestration) []*apitypes.Orchestration
	// CreateOrchestrationInternal creates a new orchestration without any
	// verification. This method should only be used for tests.
	// FIXME: Remove method: https://github.com/DuarteMRAlves/maestro/issues/245
	CreateOrchestrationInternal(orchestration *Orchestration)
	// CreateStage creates a new stage with the specified config.
	// It returns an error if the asset can not be created and nil otherwise.
	CreateStage(*apitypes.Stage) (*Stage, error)
	// CreateStageInternal creates a new stage without any verification. This
	// method should only be used for tests.
	// FIXME: Remove method: https://github.com/DuarteMRAlves/maestro/issues/245
	CreateStageInternal(*Stage)
	ContainsStage(apitypes.StageName) bool
	GetStageByName(apitypes.StageName) (*Stage, bool)
	// GetMatchingStage retrieves stored stages that match the received query.
	// The query is a stage with the fields that the returned stage should have.
	// If a field is empty, then all values for that field are accepted.
	GetMatchingStage(*apitypes.Stage) []*apitypes.Stage
	CreateLink(*apitypes.Link) (*Link, error)
	// CreateLinkInternal creates a new link without any verification. This
	// method should only be used for tests.
	// FIXME: Remove method: https://github.com/DuarteMRAlves/maestro/issues/245
	CreateLinkInternal(*Link)
	ContainsLink(name apitypes.LinkName) bool
	GetMatchingLinks(query *apitypes.Link) []*apitypes.Link
}

type manager struct {
	orchestrations sync.Map
	stages         sync.Map
	links          sync.Map

	assets asset.Store
}

func NewManager(assets asset.Store) Manager {
	return &manager{
		orchestrations: sync.Map{},
		stages:         sync.Map{},
		links:          sync.Map{},
		assets:         assets,
	}
}

func (m *manager) CreateOrchestration(cfg *apitypes.Orchestration) error {
	var err error
	if err = validateCreateOrchestrationConfig(cfg); err != nil {
		return err
	}

	o := New(cfg.Name)
	_, prev := m.orchestrations.LoadOrStore(o.Name(), o)
	if prev {
		return errdefs.AlreadyExistsWithMsg(
			"orchestration '%v' already exists",
			o.Name())
	}
	return nil
}

func (m *manager) GetMatchingOrchestration(
	query *apitypes.Orchestration,
) []*apitypes.Orchestration {
	if query == nil {
		query = &apitypes.Orchestration{}
	}
	filter := buildOrchestrationQueryFilter(query)
	res := make([]*apitypes.Orchestration, 0)
	m.orchestrations.Range(
		func(key, value interface{}) bool {
			b, ok := value.(*Orchestration)
			if !ok {
				return false
			}
			if filter(b) {
				res = append(res, b.ToApi())
			}
			return true
		})
	return res
}

func (m *manager) CreateOrchestrationInternal(o *Orchestration) {
	m.orchestrations.Store(o.Name(), o)
}

func (m *manager) CreateStage(cfg *apitypes.Stage) (*Stage, error) {
	var (
		rpc reflection.RPC
		err error
	)
	if err = m.validateCreateStageConfig(cfg); err != nil {
		return nil, err
	}
	address := m.inferStageAddress(cfg)
	rpc, err = m.inferRpc(address, cfg)
	if err != nil {
		return nil, err
	}
	s := NewStage(cfg.Name, address, cfg.Asset, rpc, nil)
	_, prev := m.stages.LoadOrStore(s.name, s)
	if prev {
		return nil, errdefs.AlreadyExistsWithMsg(
			"stage '%v' already exists",
			s.name)
	}
	return s, nil
}

func (m *manager) CreateStageInternal(s *Stage) {
	m.stages.Store(s.name, s)
}

func (m *manager) ContainsStage(name apitypes.StageName) bool {
	_, ok := m.stages.Load(name)
	return ok
}

func (m *manager) GetStageByName(name apitypes.StageName) (*Stage, bool) {
	loaded, ok := m.stages.Load(name)
	if !ok {
		return nil, false
	}
	stage, ok := loaded.(*Stage)
	return stage, ok
}

func (m *manager) GetMatchingStage(query *apitypes.Stage) []*apitypes.Stage {
	if query == nil {
		query = &apitypes.Stage{}
	}
	filter := buildStageQueryFilter(query)
	res := make([]*apitypes.Stage, 0)
	m.stages.Range(
		func(key, value interface{}) bool {
			s, ok := value.(*Stage)
			if !ok {
				return false
			}
			if filter(s) {
				res = append(res, s.ToApi())
			}
			return true
		})
	return res
}

func (m *manager) CreateLink(cfg *apitypes.Link) (*Link, error) {
	if err := m.validateCreateLinkConfig(cfg); err != nil {
		return nil, err
	}
	l := NewLink(
		cfg.Name,
		cfg.SourceStage,
		cfg.SourceField,
		cfg.TargetStage,
		cfg.TargetField,
	)
	_, prev := m.links.LoadOrStore(l.name, l)
	if prev {
		return nil, errdefs.AlreadyExistsWithMsg(
			"link '%v' already exists",
			l.name)
	}
	return l, nil
}

func (m *manager) CreateLinkInternal(l *Link) {
	m.links.Store(l.Name(), l)
}

// ContainsLink returns true if a link with the given name exists and false
// otherwise.
func (m *manager) ContainsLink(name apitypes.LinkName) bool {
	_, ok := m.links.Load(name)
	return ok
}

func (m *manager) GetMatchingLinks(query *apitypes.Link) []*apitypes.Link {
	if query == nil {
		query = &apitypes.Link{}
	}
	filter := buildLinkQueryFilter(query)
	res := make([]*apitypes.Link, 0)
	m.links.Range(
		func(key, value interface{}) bool {
			l, ok := value.(*Link)
			if !ok {
				return false
			}
			if filter(l) {
				res = append(res, l.ToApi())
			}
			return true
		})
	return res
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
) (reflection.RPC, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return nil, errdefs.InternalWithMsg(
			"connect to %s for stage %s: %s",
			address,
			cfg.Name,
			err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rpcDiscoveryCfg := &discovery.Config{
		Service: cfg.Service,
		Rpc:     cfg.Rpc,
	}
	rpc, err := discovery.FindRpc(ctx, conn, rpcDiscoveryCfg)
	if err != nil {
		return nil, errdefs.PrependMsg(err, "stage %v", cfg.Name)
	}
	return rpc, nil
}
