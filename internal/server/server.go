package server

import (
	"github.com/DuarteMRAlves/maestro/internal/execution"
	"github.com/DuarteMRAlves/maestro/internal/rpc"
	"github.com/DuarteMRAlves/maestro/internal/storage"
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
	storageManager    storage.Manager
	flowManager       execution.Manager
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
