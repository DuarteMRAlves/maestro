package server

import (
	"context"
	"fmt"
	"github.com/DuarteMRAlves/maestro/internal/api/types"
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
func (s *Server) CreateStage(config *types.Stage) error {
	var (
		st  *stage.Stage
		err error
	)
	s.logger.Info("Create Stage.", logStage(config, "config")...)
	if st, err = s.createStageFromConfig(config); err != nil {
		return err
	}
	if err = s.inferRpc(st); err != nil {
		return err
	}
	return s.stageStore.Create(st)
}

func (s *Server) GetStage(query *types.Stage) []*types.Stage {
	s.logger.Info("Get Stage.", logStage(query, "query")...)
	st := &stage.Stage{}
	st.Name = query.Name
	st.Service = query.Service
	st.Method = query.Method
	st.Asset = query.Asset
	st.Address = query.Address
	stages := s.stageStore.Get(st)
	apiStages := make([]*types.Stage, 0, len(stages))
	for _, st := range stages {
		apiStages = append(apiStages, st.ToApi())
	}
	return apiStages
}

func logStage(s *types.Stage, field string) []zap.Field {
	if s == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", s.Name),
		zap.String("asset", s.Asset),
		zap.String("service", s.Service),
		zap.String("method", s.Method),
		zap.String("address", s.Address),
		zap.String("host", s.Host),
		zap.Int32("port", s.Port),
	}
}

func (s *Server) createStageFromConfig(
	config *types.Stage,
) (*stage.Stage, error) {
	if err := s.validateCreateStageConfig(config); err != nil {
		return nil, err
	}
	st := &stage.Stage{
		Name:    config.Name,
		Service: config.Service,
		Method:  config.Method,
		Asset:   config.Asset,
		Address: config.Address,
	}
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
func (s *Server) validateCreateStageConfig(config *types.Stage) error {
	if ok, err := validate.ArgNotNil(config, "config"); !ok {
		return err
	}
	if !naming.IsValidName(config.Name) {
		return errdefs.InvalidArgumentWithMsg(
			"invalid name '%v'",
			config.Name)
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

func (s *Server) inferRpc(st *stage.Stage) error {
	address := st.Address
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return errdefs.InternalWithMsg(
			"connect to %s for stage %s: %s",
			address,
			st.Name,
			err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reflectionClient := reflection.NewClient(ctx, conn)
	availableServices, err := reflectionClient.ListServices()
	if err != nil {
		return err
	}
	serviceName, err := findService(st, availableServices)
	if err != nil {
		return err
	}
	service, err := reflectionClient.ResolveService(serviceName)
	if err != nil {
		return err
	}
	return inferRpcFromServices(st, service.RPCs())
}

// findService finds the service that should be used to call the stage method.
// It tries to find the specified service among the available services. If the
// service is not specified, then only one available service must exist that
// will be used. An error is returned if none of the above conditions is
// verified.
func findService(config *stage.Stage, available []string) (string, error) {
	search := config.Service
	if search == "" {
		if len(available) == 1 {
			return available[0], nil
		}
		return "", errdefs.InvalidArgumentWithMsg(
			"find service without name for stage %v: expected 1 "+
				"available service but %v found",
			config.Name,
			len(available))
	} else {
		for _, s := range available {
			if search == s {
				return search, nil
			}
		}
		return "", errdefs.NotFoundWithMsg(
			"service with name %v not found for stage %v",
			search,
			config.Name)
	}
}

// inferRpcFromServices verifies that the rpc to be called for the stage exists. If a
// rpc was specified in the config, then it verifies it exists in the available
// methods. Otherwise, it verifies only a single method is available.
func inferRpcFromServices(st *stage.Stage, available []reflection.RPC) error {
	search := st.Method
	if search == "" {
		if len(available) == 1 {
			return nil
		}
		return errdefs.InvalidArgumentWithMsg(
			"find rpc without name for stage %v: expected 1 available "+
				"rpc but %v found",
			st.Name,
			len(available))
	} else {
		for _, rpc := range available {
			if search == rpc.Name() {
				st.Rpc = rpc
				return nil
			}
		}
		return errdefs.NotFoundWithMsg(
			"rpc with name %v not found for stage %v",
			search,
			st.Name)
	}
}
