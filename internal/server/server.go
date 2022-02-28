package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/arch"
	"github.com/DuarteMRAlves/maestro/internal/errdefs"
	"github.com/DuarteMRAlves/maestro/internal/exec"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

const grpcNotConfigured = "grpc server not configured"

// Server is the main class that handles the requests
// It implements the InternalAPI interface and manages all requests
type Server struct {
	execManager exec.Manager

	grpcServer *grpc.Server

	// db is a key-value store database to persist state across multiple
	// executions of the server and to ensure consistency.
	db *badger.DB

	logger *zap.Logger
}

func (s *Server) ServeGrpc(lis net.Listener) error {
	if s.grpcServer == nil {
		return errdefs.FailedPreconditionWithMsg(grpcNotConfigured)
	}
	return s.grpcServer.Serve(lis)
}

func (s *Server) GracefulStopGrpc() {
	s.grpcServer.GracefulStop()
}

func (s *Server) StopGrpc() {
	s.grpcServer.Stop()
}

// CreateAsset creates a new asset with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateAsset(req *api.CreateAssetRequest) error {
	fields := logs.FieldsForCreateAssetRequest(req)
	s.logger.Info("Create Asset.", fields...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			createAsset := arch.CreateAssetWithTxn(txn)
			return createAsset(req)
		},
	)
}

func (s *Server) GetAsset(req *api.GetAssetRequest) (
	[]*api.Asset,
	error,
) {
	var (
		assets []*api.Asset
		err    error
	)
	fields := logs.FieldsForGetAssetRequest(req)
	s.logger.Info("Get Asset.", fields...)
	err = s.db.View(
		func(txn *badger.Txn) error {
			getAssets := arch.GetAssetsWithTxn(txn)
			assets, err = getAssets(req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (s *Server) CreateOrchestration(req *api.CreateOrchestrationRequest) error {
	fields := logs.FieldsForCreateOrchestrationRequest(req)
	s.logger.Info("Create Orchestration.", fields...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			createOrchestration := arch.CreateOrchestrationWithTxn(txn)
			return createOrchestration(req)
		},
	)
}

func (s *Server) GetOrchestration(
	req *api.GetOrchestrationRequest,
) ([]*api.Orchestration, error) {
	var (
		orchestrations []*api.Orchestration
		err            error
	)
	fields := logs.FieldsForGetOrchestrationRequest(req)
	s.logger.Info("Get Orchestration.", fields...)
	err = s.db.View(
		func(txn *badger.Txn) error {
			getOrchestrations := arch.GetOrchestrationsWithTxn(txn)
			orchestrations, err = getOrchestrations(req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return orchestrations, nil
}

// CreateStage creates a new stage with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateStage(req *api.CreateStageRequest) error {
	fields := logs.FieldsForCreateStageRequest(req)
	s.logger.Info("Create Stage.", fields...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			createStage := arch.CreateStageWithTxn(txn)
			return createStage(req)
		},
	)
}

func (s *Server) GetStage(req *api.GetStageRequest) ([]*api.Stage, error) {
	var (
		stages []*api.Stage
		err    error
	)
	fields := logs.FieldsForGetStageRequest(req)
	s.logger.Info("Get Stage.", fields...)
	err = s.db.View(
		func(txn *badger.Txn) error {
			getStages := arch.GetStagesWithTxn(txn)
			stages, err = getStages(req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return stages, nil
}

// CreateLink creates a new link with the specified config.
// It returns an error if the asset can not be created and nil otherwise.
func (s *Server) CreateLink(req *api.CreateLinkRequest) error {
	fields := logs.FieldsForCreateLinkRequest(req)
	s.logger.Info("Create Link.", fields...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			createLink := arch.CreateLinkWithTxn(txn)
			return createLink(req)
		},
	)
}

func (s *Server) GetLink(req *api.GetLinkRequest) ([]*api.Link, error) {
	var (
		links []*api.Link
		err   error
	)
	fields := logs.FieldsForGetLinkRequest(req)
	s.logger.Info("Get Link.", fields...)
	err = s.db.View(
		func(txn *badger.Txn) error {
			getLinks := arch.GetLinksWithTxn(txn)
			links, err = getLinks(req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (s *Server) StartExecution(req *api.StartExecutionRequest) error {
	fields := logs.FieldsForStartExecutionRequest(req)
	s.logger.Info("Start Execution.", fields...)
	return s.db.Update(
		func(txn *badger.Txn) error {
			return s.execManager.StartExecution(txn, req)
		},
	)
}

func (s *Server) AttachExecution(req *api.AttachExecutionRequest) (
	*api.Subscription,
	error,
) {
	fields := logs.FieldsForAttachExecutionRequest(req)
	s.logger.Info("Attach Execution.", fields...)
	return s.execManager.AttachExecution(req)
}
