package server

import (
	"context"
	"fmt"
	apitypes "github.com/DuarteMRAlves/maestro/internal/api/types"
	"github.com/DuarteMRAlves/maestro/internal/discovery"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/naming"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"github.com/DuarteMRAlves/maestro/internal/worker"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

// CreateStage creates a new stage with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateStage(config *apitypes.Stage) error {
	var (
		st  *stage.Stage
		err error
	)
	s.logger.Info("Create Stage.", logStage(config, "config")...)
	if st, err = s.createStageFromConfig(config); err != nil {
		return err
	}
	if err = s.inferRpc(st, config); err != nil {
		return err
	}
	if st.Worker, err = worker.NewWorker(st.Address, st.Rpc); err != nil {
		return err
	}
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
		zap.String("name", s.Name),
		zap.String("asset", string(s.Asset)),
		zap.String("service", s.Service),
		zap.String("rpc", s.Rpc),
		zap.String("address", s.Address),
		zap.String("host", s.Host),
		zap.Int32("port", s.Port),
	}
}

func (s *Server) createStageFromConfig(
	config *apitypes.Stage,
) (*stage.Stage, error) {
	if err := s.validateCreateStageConfig(config); err != nil {
		return nil, err
	}
	st := stage.NewDefault()
	st.Name = config.Name
	st.Asset = config.Asset
	st.Address = config.Address
	// If address is empty, fill it from config host and port.
	if st.Address == "" {
		host, port := config.Host, config.Port
		if host == "" {
			host = "localhost"
		}
		if port == 0 {
			port = 8061
		}
		st.Address = fmt.Sprintf("%s:%d", host, port)
	}
	return st, nil
}

// validateCreateStageConfig verifies if all conditions to create a stage are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateStageConfig(config *apitypes.Stage) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}
	if !naming.IsValidName(config.Name) {
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

func (s *Server) inferRpc(st *stage.Stage, cfg *apitypes.Stage) error {
	address := st.Address
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		return errdefs.InternalWithMsg(
			"connect to %s for stage %s: %s",
			address,
			st.Name,
			err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	rpcDiscoveryCfg := &discovery.Config{
		Service: cfg.Service,
		Rpc:     cfg.Rpc,
	}
	st.Rpc, err = discovery.FindRpc(ctx, conn, rpcDiscoveryCfg)
	if err != nil {
		return errdefs.PrependMsg(err, "stage %v", st.Name)
	}
	return nil
}
