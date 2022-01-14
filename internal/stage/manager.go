package stage

import (
	"context"
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/discovery"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"google.golang.org/grpc"
	"time"
)

type Manager interface {
	Create(cfg *apitypes.Stage) error
	Get(query *apitypes.Stage) []*apitypes.Stage
}

type manager struct {
	stages Store
	assets asset.Store
}

func NewManager(stages Store, assets asset.Store) Manager {
	return &manager{stages: stages, assets: assets}
}

func (m *manager) Create(cfg *apitypes.Stage) error {
	if err := m.validateCreateStageConfig(cfg); err != nil {
		return err
	}
	address := m.inferStageAddress(cfg)
	rpc, err := m.inferRpc(address, cfg)
	if err != nil {
		return err
	}
	st := New(cfg.Name, address, cfg.Asset, rpc)
	return m.stages.Create(st)
}

func (m *manager) Get(query *apitypes.Stage) []*apitypes.Stage {
	stages := m.stages.GetMatching(query)
	apiStages := make([]*apitypes.Stage, 0, len(stages))
	for _, st := range stages {
		apiStages = append(apiStages, st.ToApi())
	}
	return apiStages
}

// validateCreateStageConfig verifies if all conditions to create a stage are met.
// It returns an error if a condition is not met and nil otherwise.
func (m *manager) validateCreateStageConfig(cfg *apitypes.Stage) error {
	if ok, err := validate.ArgNotNil(cfg, "cfg"); !ok {
		return err
	}
	if !naming.IsValidStageName(cfg.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			cfg.Name)
	}
	if cfg.Phase != "" {
		return errdefs.InvalidArgumentWithMsg("phase should not be specified")
	}
	// Asset is not required but if specified should exist.
	if cfg.Asset != "" && !m.assets.Contains(cfg.Asset) {
		return errdefs.NotFoundWithMsg(
			"asset '%v' not found",
			cfg.Asset)
	}
	if cfg.Address != "" && cfg.Host != "" {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and host for stage")
	}
	if cfg.Address != "" && cfg.Port != 0 {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and port for stage")
	}
	return nil
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
