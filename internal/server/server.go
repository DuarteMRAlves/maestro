package server

import (
	"github.com/DuarteMRAlves/maestro/internal/asset"
	"github.com/DuarteMRAlves/maestro/internal/flow"
	"github.com/DuarteMRAlves/maestro/internal/link"
	"github.com/DuarteMRAlves/maestro/internal/orchestration"
	"github.com/DuarteMRAlves/maestro/internal/stage"
	"github.com/DuarteMRAlves/maestro/internal/validate"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

const grpcNotConfigured = "grpc server not configured"

// Server is the main class that handles the requests
// It implements the InternalAPI interface and manages all requests
type Server struct {
	assetStore         asset.Store
	stageStore         stage.Store
	linkStore          link.Store
	orchestrationStore orchestration.Store

	flowManager *flow.Manager

	grpcServer *grpc.Server

	logger *zap.Logger
}

func (s *Server) ServeGrpc(lis net.Listener) error {
	if ok, err := validate.Status(s.grpcServer != nil, grpcNotConfigured); !ok {
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
