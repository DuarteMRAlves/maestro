package server

import (
	"context"
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
func (s *Server) CreateStage(config *stage.Stage) error {
	s.logger.Info("Create Stage.", logStage(config, "config")...)
	if err := s.validateCreateStageConfig(config); err != nil {
		return err
	}
	if err := s.validateRpcExists(config); err != nil {
		return err
	}
	return s.stageStore.Create(config)
}

func (s *Server) GetStage(query *stage.Stage) []*stage.Stage {
	s.logger.Info("Get Stage.", logStage(query, "query")...)
	return s.stageStore.Get(query)
}

func logStage(s *stage.Stage, field string) []zap.Field {
	if s == nil {
		return []zap.Field{zap.String(field, "null")}
	}
	return []zap.Field{
		zap.String("name", s.Name),
		zap.String("asset", s.Asset),
		zap.String("service", s.Service),
		zap.String("method", s.Method),
		zap.String("address", s.Address),
	}
}

// validateCreateStageConfig verifies if all conditions to create a stage are met.
// It returns an error if a condition is not met and nil otherwise.
func (s *Server) validateCreateStageConfig(config *stage.Stage) error {
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
	return nil
}

func (s *Server) validateRpcExists(config *stage.Stage) error {
	address := inferAddress(config)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return errdefs.InternalWithMsg("connect to %s: %s", address, err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reflectionClient := reflection.NewClient(ctx, conn)
	availableServices, err := reflectionClient.ListServices()
	if err != nil {
		return err
	}
	serviceName, err := findService(config, availableServices)
	if err != nil {
		return err
	}
	service, err := reflectionClient.ResolveService(serviceName)
	if err != nil {
		return err
	}
	return validateRPC(config, service.RPCs())
}

func inferAddress(config *stage.Stage) string {
	if config.Address != "" {
		return config.Address
	} else {
		// FIXME: Add methods to infer address from partial address or asset.
		return "localhost:8061"
	}
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

// validateRPC verifies that the rpc to be called for the stage exists. If a
// rpc was specified in the config, then it verifies it exists in the available
// methods. Otherwise, it verifies only a single method is available.
func validateRPC(config *stage.Stage, available []reflection.RPC) error {
	search := config.Method
	if search == "" {
		if len(available) == 1 {
			return nil
		}
		return errdefs.InvalidArgumentWithMsg(
			"find rpc without name for stage %v: expected 1 available "+
				"rpc but %v found",
			config.Name,
			len(available))
	} else {
		for _, rpc := range available {
			if search == rpc.Name() {
				return nil
			}
		}
		return errdefs.NotFoundWithMsg(
			"rpc with name %v not found for stage %v",
			search,
			config.Name)
	}
}
