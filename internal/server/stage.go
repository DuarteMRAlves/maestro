package server

import (
	"context"
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/discovery"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/reflection"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

// CreateStage creates a new stage with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateStage(config *apitypes.Stage) error {
	var err error
	s.logger.Info("Create Stage.", logStage(config, "config")...)
	if err := s.validateCreateStageConfig(config); err != nil {
		return err
	}
	address := s.inferStageAddress(config)
	rpc, err := s.inferRpc(address, config)
	if err != nil {
		return err
	}
	st := stage.New(config.Name, address, config.Asset, rpc)
	return s.stageStore.Create(st)
}

func (s *Server) GetStage(query *apitypes.Stage) []*apitypes.Stage {
	s.logger.Info("Get Stage.", logStage(query, "query")...)
	stages := s.stageStore.GetMatching(query)
	apiStages := make([]*apitypes.Stage, 0, len(stages))
	for _, st := range stages {
		apiStages = append(apiStages, st.ToApi())
	}
	return apiStages
}

func logStage(s *apitypes.Stage, field string) []zap.Field {
	if s == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", string(s.Name)),
		zap.String("asset", string(s.Asset)),
		zap.String("service", s.Service),
		zap.String("rpc", s.Rpc),
		zap.String("address", s.Address),
		zap.String("host", s.Host),
		zap.Int32("port", s.Port),
	}
}

// validateCreateStageConfig verifies if all conditions to create a stage are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateStageConfig(config *apitypes.Stage) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}
	if !naming.IsValidStageName(config.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			config.Name)
	}
	if config.Phase != "" {
		return errdefs.InvalidArgumentWithMsg("phase should not be specified")
	}
	// Asset is not required but if specified should exist.
	if config.Asset != "" && !s.assetStore.Contains(config.Asset) {
		return errdefs.NotFoundWithMsg(
			"asset '%v' not found",
			config.Asset)
	}
	if config.Address != "" && config.Host != "" {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and host for stage")
	}
	if config.Address != "" && config.Port != 0 {
		return errdefs.InvalidArgumentWithMsg(
			"Cannot simultaneously specify address and port for stage")
	}
	return nil
}

func (s *Server) inferStageAddress(config *apitypes.Stage) string {
	address := config.Address
	// If address is empty, fill it from config host and port.
	if address == "" {
		host, port := config.Host, config.Port
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

func (s *Server) inferRpc(
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
