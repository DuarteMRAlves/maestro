package server

import (
	"github.com/DuarteMRAlves/maestro/internal/api"
	"github.com/DuarteMRAlves/maestro/internal/arch"
	"github.com/DuarteMRAlves/maestro/internal/exec"
	"github.com/DuarteMRAlves/maestro/internal/logs"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/util"
	"github.com/dgraph-io/badger/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

const grpcNotConfigured = "grpc server not configured"

// Server is the main class that handles the requests
// It implements the InternalAPI interface and manages all requests
type Server struct {
	storageManager    arch.Manager
	flowManager       exec.Manager
	reflectionManager rpc.Manager

	grpcServer *grpc.Server

	// db is a key-value store database to persist state across multiple
	// executions of the server and to ensure consistency.
	db *badger.DB

	logger *zap.Logger
}

func (s *Server) ServeGrpc(lis net.Listener) error {
	if ok, err := util.Status(s.grpcServer != nil, grpcNotConfigured); !ok {
		return err
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
	logs.LogCreateAssetRequest(s.logger, req)
	return s.db.Update(
		func(txn *badger.Txn) error {
			return s.storageManager.CreateAsset(txn, req)
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
	logs.LogGetAssetRequest(s.logger, req)
	err = s.db.View(
		func(txn *badger.Txn) error {
			assets, err = s.storageManager.GetMatchingAssets(txn, req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return assets, nil
}

func (s *Server) CreateOrchestration(req *api.CreateOrchestrationRequest) error {
	logs.LogCreateOrchestrationRequest(s.logger, req)
	return s.db.Update(
		func(txn *badger.Txn) error {
			return s.storageManager.CreateOrchestration(txn, req)
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
	logs.LogGetOrchestrationRequest(s.logger, req)
	err = s.db.View(
		func(txn *badger.Txn) error {
			orchestrations, err = s.storageManager.GetMatchingOrchestration(
				txn,
				req,
			)
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
	logs.LogCreateStageRequest(s.logger, req)
	return s.db.Update(
		func(txn *badger.Txn) error {
			return s.storageManager.CreateStage(txn, req)
		},
	)
}

func (s *Server) GetStage(req *api.GetStageRequest) ([]*api.Stage, error) {
	var (
		stages []*api.Stage
		err    error
	)
	logs.LogGetStageRequest(s.logger, req)
	err = s.db.View(
		func(txn *badger.Txn) error {
			stages, err = s.storageManager.GetMatchingStage(txn, req)
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
	logs.LogCreateLinkRequest(s.logger, req)
	return s.db.Update(
		func(txn *badger.Txn) error {
			return s.storageManager.CreateLink(txn, req)
		},
	)
}

func (s *Server) GetLink(req *api.GetLinkRequest) ([]*api.Link, error) {
	var (
		links []*api.Link
		err   error
	)
	logs.LogGetLinkRequest(s.logger, req)
	err = s.db.View(
		func(txn *badger.Txn) error {
			links, err = s.storageManager.GetMatchingLinks(txn, req)
			return err
		},
	)
	if err != nil {
		return nil, err
	}
	return links, nil
}
